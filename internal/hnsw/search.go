package hnsw

import (
	"container/heap"
	"sort"
)

// searchLayer runs an ef-bounded greedy search at a single layer.
//
// query: the query vector
// entryPoints: starting node IDs at this layer
// ef: candidate set capacity (== 1 for greedy descent, efSearch/efConstruction otherwise)
// layer: which layer to search
//
// Returns results sorted ascending by Dist.
func (idx *Index) searchLayer(query []float32, entryPoints []uint32, ef, layer int) []Candidate {
	visited := make(map[uint32]bool, ef*2)
	candidates := &MinHeap{}
	results := &MaxHeap{}
	dist := idx.Distance

	// Init: push every entry point into both candidates and results
	for _, ep := range entryPoints {
		d := dist(query, idx.nodes[ep].Vector)
		visited[ep] = true
		heap.Push(candidates, Candidate{ep, d})
		heap.Push(results, Candidate{ep, d})
	}

	for candidates.Len() > 0 {
		// 1. Pop closest unexpanded node c from candidates
		closest := heap.Pop(candidates).(Candidate)

		// 2. If c.Dist > furthest in results AND len(results) == ef, break
		if results.Len() == ef && closest.Dist > results.Peek().Dist {
			break
		}

		// 3. For each unvisited neighbor of c at this layer:
		node := idx.nodes[closest.ID]
		if layer >= len(node.Neighbors) {
			// node does not exist on this layer; skip
			continue
		}
		for _, nbrID := range node.Neighbors[layer] {
			if visited[nbrID] {
				continue
			}
			visited[nbrID] = true

			d := dist(query, idx.nodes[nbrID].Vector)

			// 3c. push if closer than furthest result OR results not full
			if results.Len() < ef || d < results.Peek().Dist {
				heap.Push(candidates, Candidate{nbrID, d})
				heap.Push(results, Candidate{nbrID, d})
				// 3d. if len(results) > ef -> evict furthest
				if results.Len() > ef {
					heap.Pop(results)
				}
			}
		}
	}

	// Return results sorted ascending by Dist
	out := make([]Candidate, results.Len())
	for i := len(out) - 1; i >= 0; i-- {
		out[i] = heap.Pop(results).(Candidate)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Dist < out[j].Dist })
	return out
}

// Search is the top-level KNN query.
// k: k: number of nearest neighbors to return
// efSearch: efSearch: candidate set size at query time (runtime tunable)
func (idx *Index) Search(query []float32, k, efSearch int) []Candidate {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if !idx.hasEntry {
		return nil
	}

	entryPoints := []uint32{idx.entryPoint}
	// Greedy descent from maxLevel down to 1
	for l := idx.maxLevel; l > 0; l-- {
		res := idx.searchLayer(query, entryPoints, 1, l)
		if len(res) > 0 {
			entryPoints = []uint32{res[0].ID}
		}
	}

	// Layer 0 with efSearch
	candidates := idx.searchLayer(query, entryPoints, efSearch, 0)

	// Filter tombstones (Day 11 starts using Deleted; this is a no-op until then)
	out := make([]Candidate, 0, k)
	for _, c := range candidates {
		if idx.nodes[c.ID].Deleted {
			continue
		}
		out = append(out, c)
		if len(out) == k {
			break
		}
	}
	return out
}
