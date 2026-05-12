package hnsw

import "sort"

// BruteForceKNN is a brute-force KNN used as ground truth in tests.
func BruteForceKNN(data map[uint32][]float32, query []float32, k int, dist DistanceFunc) []Candidate {
	all := make([]Candidate, 0, len(data))
	for id, v := range data {
		all = append(all, Candidate{id, dist(query, v)})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Dist < all[j].Dist })
	if len(all) > k {
		all = all[:k]
	}
	return all
}

// DistanceFunc alias (re-exported for callers)
type DistanceFunc func(a, b []float32) float32
