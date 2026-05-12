package hnsw

import (
	"container/heap"
	"sort"
)

// searchLayer runs an ef-bounded greedy search at a single layer.
// Safe to call concurrently from multiple goroutines; per-Node RLocks
// guard reads of Neighbors against concurrent Insert mutations.
//
// query: the query vector
// entryPoints: starting node IDs at this layer
// ef: candidate set capacity (1 for greedy descent, larger for layer 0)
// layer: which layer to search
//
// Returns results sorted ascending by Dist.
func (idx *Index) searchLayer(query []float32, entryPoints []uint32, ef, layer int) []Candidate {
	visited := getBitset()
	defer putBitset(visited)
	candidates := &MinHeap{}
	results := &MaxHeap{}
	dist := idx.Distance

	for _, ep := range entryPoints {
		epNode := idx.getNode(ep)
		if epNode == nil {
			continue
		}
		// epNode.Vector is immutable after Insert returns; no lock needed.
		d := dist(query, epNode.Vector)
		visited.Set(ep)
		heap.Push(candidates, Candidate{ep, d})
		heap.Push(results, Candidate{ep, d})
	}

	for candidates.Len() > 0 {
		closest := heap.Pop(candidates).(Candidate)
		if results.Len() == ef && closest.Dist > results.Peek().Dist {
			break
		}

		cn := idx.getNode(closest.ID)
		if cn == nil {
			continue
		}

		// Hold the RLock for the inner loop. Distance is fast (SIMD), so the
		// critical section is short; in exchange we don't pay the cost of
		// copying the neighbors slice on every expansion.
		cn.mu.RLock()
		if layer >= len(cn.Neighbors) {
			cn.mu.RUnlock()
			continue
		}
		for _, nbrID := range cn.Neighbors[layer] {
			if visited.Test(nbrID) {
				continue
			}
			visited.Set(nbrID)

			nbr := idx.getNode(nbrID)
			if nbr == nil {
				continue
			}
			d := dist(query, nbr.Vector)

			if results.Len() < ef || d < results.Peek().Dist {
				heap.Push(candidates, Candidate{nbrID, d})
				heap.Push(results, Candidate{nbrID, d})
				if results.Len() > ef {
					heap.Pop(results)
				}
			}
		}
		cn.mu.RUnlock()
	}

	out := make([]Candidate, results.Len())
	for i := len(out) - 1; i >= 0; i-- {
		out[i] = heap.Pop(results).(Candidate)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Dist < out[j].Dist })
	return out
}

// Search is the top-level KNN query. Safe to call concurrently.
func (idx *Index) Search(query []float32, k, efSearch int) []Candidate {
	idx.globalMu.Lock()
	hasEntry := idx.hasEntry
	entryPoint := idx.entryPoint
	maxLevel := idx.maxLevel
	idx.globalMu.Unlock()

	if !hasEntry {
		return nil
	}

	entryPoints := []uint32{entryPoint}
	for l := maxLevel; l > 0; l-- {
		res := idx.searchLayer(query, entryPoints, 1, l)
		if len(res) > 0 {
			entryPoints = []uint32{res[0].ID}
		}
	}

	candidates := idx.searchLayer(query, entryPoints, efSearch, 0)

	out := make([]Candidate, 0, k)
	for _, c := range candidates {
		n := idx.getNode(c.ID)
		if n == nil {
			continue
		}
		n.mu.RLock()
		deleted := n.Deleted
		n.mu.RUnlock()
		if deleted {
			continue
		}
		out = append(out, c)
		if len(out) == k {
			break
		}
	}
	return out
}
