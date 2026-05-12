package hnsw

// selectNeighborsHeuristic implements HNSW Algorithm 4
// (SELECT-NEIGHBORS-HEURISTIC) from Malkov & Yashunin (2018).
//
// Acceptance rule for a candidate c:
//
//	for every already-accepted neighbor r:
//	    distance(c.Vector, r.Vector) >= c.Dist
//
// Geometrically: c is accepted only if no already-accepted point is closer
// to c than the query is. An accepted r that is closer to c than the query
// lies in c's occluded region — you can already reach c via r, so the
// direct edge q->c is redundant.
//
// Reads other nodes' Vector fields, which are immutable after Insert
// returns, so this routine does not need to acquire any per-Node locks.
//
// Preconditions:
//   - candidates is sorted ascending by Dist (closest to query first).
//   - Every candidate.ID exists in idx.nodes.
//   - M > 0.
//
// Returns accepted neighbor IDs in acceptance order, len <= M.
func (idx *Index) selectNeighborsHeuristic(query []float32, candidates []Candidate, M int) []uint32 {
	if M <= 0 || len(candidates) == 0 {
		return nil
	}

	accepted := make([]uint32, 0, M)
	acceptedVecs := make([][]float32, 0, M)

	for _, c := range candidates {
		if len(accepted) >= M {
			break
		}

		cNode := idx.getNode(c.ID)
		if cNode == nil {
			continue
		}
		cVec := cNode.Vector

		shouldAccept := true
		for _, rVec := range acceptedVecs {
			if idx.Distance(cVec, rVec) < c.Dist {
				shouldAccept = false
				break
			}
		}

		if shouldAccept {
			accepted = append(accepted, c.ID)
			acceptedVecs = append(acceptedVecs, cVec)
		}
	}

	return accepted
}
