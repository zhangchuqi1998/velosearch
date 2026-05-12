package hnsw

// RecallAtK computes the overlap ratio between HNSW results and ground truth.
// Both inputs are length-k slices sorted ascending by Dist.
func RecallAtK(hnsw, truth []Candidate) float64 {
	set := make(map[uint32]bool, len(truth))
	for _, c := range truth {
		set[c.ID] = true
	}
	hit := 0
	for _, c := range hnsw {
		if set[c.ID] {
			hit++
		}
	}
	return float64(hit) / float64(len(truth))
}
