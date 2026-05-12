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
