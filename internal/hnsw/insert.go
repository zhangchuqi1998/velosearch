package hnsw

import "sort"

// Insert adds a vector with the given ID into the index.
//
// Algorithm 1 from Malkov & Yashunin (2018):
//  1. Pick a random level for the new node (geometric distribution).
//  2. Greedy descent from maxLevel down to (level+1), ef=1.
//  3. From min(level, maxLevel) down to layer 0, ef=efConstruction:
//     search the layer, pick M neighbors via the heuristic, add
//     bidirectional edges, and prune over-full neighbors.
//  4. If the new node's level exceeds maxLevel, promote it to entry point.
func (idx *Index) Insert(id uint32, vector []float32) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	level := idx.randomLevel()
	node := &Node{
		ID:        id,
		Vector:    vector,
		Level:     level,
		Neighbors: make([][]uint32, level+1),
	}
	idx.nodes[id] = node

	// First node ever: it is the entry point.
	if !idx.hasEntry {
		idx.entryPoint = id
		idx.maxLevel = level
		idx.hasEntry = true
		return
	}

	entryPoints := []uint32{idx.entryPoint}

	// Greedy descent: maxLevel down to level+1, ef=1 per layer
	for l := idx.maxLevel; l > level; l-- {
		res := idx.searchLayer(vector, entryPoints, 1, l)
		if len(res) > 0 {
			entryPoints = []uint32{res[0].ID}
		}
	}

	// From min(level, maxLevel) down to 0, ef=efConstruction
	start := level
	if idx.maxLevel < start {
		start = idx.maxLevel
	}
	for l := start; l >= 0; l-- {
		candidates := idx.searchLayer(vector, entryPoints, idx.EfConstruction, l)

		Mlayer := idx.MaxM
		if l == 0 {
			Mlayer = idx.MaxM0
		}
		chosen := idx.selectNeighborsHeuristic(vector, candidates, Mlayer)

		// Add bidirectional edges.
		node.Neighbors[l] = chosen
		for _, nid := range chosen {
			other := idx.nodes[nid]
			other.Neighbors[l] = append(other.Neighbors[l], id)
			if len(other.Neighbors[l]) > Mlayer {
				idx.pruneNeighbors(other, l, Mlayer)
			}
		}

		// Next layer's entry points: all current-layer candidates.
		entryPoints = make([]uint32, len(candidates))
		for i, c := range candidates {
			entryPoints[i] = c.ID
		}
	}

	// Promote to entry point if this node went higher than any before.
	if level > idx.maxLevel {
		idx.entryPoint = id
		idx.maxLevel = level
	}
}

// pruneNeighbors trims n.Neighbors[layer] down to M using the same
// heuristic as initial neighbor selection. We re-run the heuristic
// over n's current neighbors with n.Vector as the "query".
//
// Note: this only updates n's outgoing edges. The dropped neighbors
// keep their back-edge to n until they themselves get pruned, which
// is acceptable for v0.1 (paper section 4.1, lazy cleanup).
func (idx *Index) pruneNeighbors(n *Node, layer, M int) {
	current := n.Neighbors[layer]
	if len(current) <= M {
		return
	}

	cands := make([]Candidate, 0, len(current))
	for _, nid := range current {
		other, ok := idx.nodes[nid]
		if !ok {
			continue
		}
		d := idx.Distance(n.Vector, other.Vector)
		cands = append(cands, Candidate{ID: nid, Dist: d})
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i].Dist < cands[j].Dist })

	n.Neighbors[layer] = idx.selectNeighborsHeuristic(n.Vector, cands, M)
}
