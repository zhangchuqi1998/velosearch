package hnsw

import (
	"errors"
	"sort"
)

// Insert adds (id, vector) to the index. Safe to call from multiple
// goroutines concurrently; per-Node RWMutexes protect Neighbors mutation,
// global-state locks protect entryPoint/maxLevel, and lock acquisition
// proceeds in ascending Node-ID order to prevent deadlock.
func (idx *Index) Insert(id uint32, vector []float32) {
	level := idx.randomLevel()
	node := &Node{
		ID:        id,
		Vector:    vector,
		Level:     level,
		Neighbors: make([][]uint32, level+1),
	}

	// First-insertion fast path. globalMu serializes all "is this the first
	// node?" checks so only one goroutine wins.
	idx.globalMu.Lock()
	if !idx.hasEntry {
		idx.nodesMu.Lock()
		idx.nodes[id] = node
		idx.nodesMu.Unlock()
		idx.entryPoint = id
		idx.maxLevel = level
		idx.hasEntry = true
		idx.globalMu.Unlock()
		return
	}
	entryPoint := idx.entryPoint
	maxLevel := idx.maxLevel
	idx.globalMu.Unlock()

	// Publish the new node to the map BEFORE any other goroutine could see
	// it referenced as a neighbor. Its Neighbors slices are empty per layer
	// so concurrent searchLayer calls will just skip past it.
	idx.nodesMu.Lock()
	idx.nodes[id] = node
	idx.nodesMu.Unlock()

	// Greedy descent: maxLevel down to level+1, ef=1.
	entryPoints := []uint32{entryPoint}
	for l := maxLevel; l > level; l-- {
		res := idx.searchLayer(vector, entryPoints, 1, l)
		if len(res) > 0 {
			entryPoints = []uint32{res[0].ID}
		}
	}

	// From min(level, maxLevel) down to 0, with ef = efConstruction.
	start := level
	if maxLevel < start {
		start = maxLevel
	}
	for l := start; l >= 0; l-- {
		candidates := idx.searchLayer(vector, entryPoints, idx.EfConstruction, l)

		Mlayer := idx.MaxM
		if l == 0 {
			Mlayer = idx.MaxM0
		}
		chosen := idx.selectNeighborsHeuristic(vector, candidates, Mlayer)

		// Build lock set: self + chosen neighbors, deduped, sorted by ID.
		// Sorted acquisition order is what prevents deadlock between
		// concurrent Inserts whose lock sets overlap.
		lockSet := idx.buildLockSet(node, chosen)
		for _, n := range lockSet {
			n.mu.Lock()
		}

		// Outgoing edges from self at this layer. Use append (not assign)
		// because another goroutine's back-edge may have already landed
		// while we were searching, and we must not clobber it.
		node.Neighbors[l] = append(node.Neighbors[l], chosen...)
		if len(node.Neighbors[l]) > Mlayer {
			idx.pruneNeighborsLocked(node, l, Mlayer)
		}

		// Back-edges from chosen neighbors to self.
		for _, cid := range chosen {
			if cid == id {
				continue
			}
			cn := idx.getNode(cid)
			if cn == nil {
				continue
			}
			cn.Neighbors[l] = append(cn.Neighbors[l], id)
			if len(cn.Neighbors[l]) > Mlayer {
				idx.pruneNeighborsLocked(cn, l, Mlayer)
			}
		}

		// Release in reverse order (safe regardless, but conventional).
		for i := len(lockSet) - 1; i >= 0; i-- {
			lockSet[i].mu.Unlock()
		}

		// Next layer's entry points come from current layer's full candidate set.
		entryPoints = make([]uint32, len(candidates))
		for i, c := range candidates {
			entryPoints[i] = c.ID
		}
	}

	// Promote to entry point if our level exceeds the current max.
	if level > maxLevel {
		idx.globalMu.Lock()
		if level > idx.maxLevel {
			idx.entryPoint = id
			idx.maxLevel = level
		}
		idx.globalMu.Unlock()
	}
}

// buildLockSet collects the set of Nodes whose mu must be held while
// adding bidirectional edges between `node` and `chosen`. The result is
// sorted by Node.ID ascending so all callers acquire in the same order.
func (idx *Index) buildLockSet(node *Node, chosen []uint32) []*Node {
	out := make([]*Node, 0, len(chosen)+1)
	out = append(out, node)
	seen := make(map[uint32]bool, len(chosen)+1)
	seen[node.ID] = true
	for _, cid := range chosen {
		if seen[cid] {
			continue
		}
		seen[cid] = true
		if cn := idx.getNode(cid); cn != nil {
			out = append(out, cn)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// pruneNeighborsLocked trims n.Neighbors[layer] down to M. Caller MUST hold
// n.mu.Lock(). Reads other nodes' Vector fields, which are immutable
// after Insert so do not need locks.
func (idx *Index) pruneNeighborsLocked(n *Node, layer, M int) {
	current := n.Neighbors[layer]
	if len(current) <= M {
		return
	}

	cands := make([]Candidate, 0, len(current))
	for _, nid := range current {
		other := idx.getNode(nid)
		if other == nil {
			continue
		}
		d := idx.Distance(n.Vector, other.Vector)
		cands = append(cands, Candidate{ID: nid, Dist: d})
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i].Dist < cands[j].Dist })

	n.Neighbors[layer] = idx.selectNeighborsHeuristic(n.Vector, cands, M)
}

// Delete marks the node as a tombstone. Searches filter tombstones out
// (see Search in search.go) but the graph structure is preserved so
// deleted nodes can still serve as transit points.
func (idx *Index) Delete(id uint32) error {
	n := idx.getNode(id)
	if n == nil {
		return errors.New("id not found")
	}
	n.mu.Lock()
	n.Deleted = true
	n.mu.Unlock()
	return nil
}

// Stats is the on-the-fly aggregate used by gRPC handlers.
type Stats struct {
	NumVectors int
	NumDeleted int
	NumLayers  int
	MemBytes   int64
}

// SnapshotStats reports approximate live counts. Acquires global locks
// briefly; usable from any goroutine.
func (idx *Index) SnapshotStats() Stats {
	idx.nodesMu.RLock()
	idx.globalMu.Lock()
	s := Stats{NumVectors: len(idx.nodes), NumLayers: idx.maxLevel + 1}
	idx.globalMu.Unlock()

	deleted := 0
	for _, n := range idx.nodes {
		n.mu.RLock()
		if n.Deleted {
			deleted++
		}
		n.mu.RUnlock()
	}
	s.NumDeleted = deleted
	s.MemBytes = int64(s.NumVectors) * int64(idx.Dim) * 4
	idx.nodesMu.RUnlock()
	return s
}
