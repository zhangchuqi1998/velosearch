package hnsw

// selectNeighborsHeuristic implements HNSW Algorithm 4
// (SELECT-NEIGHBORS-HEURISTIC) from Malkov & Yashunin (2018).
//
// The heuristic picks neighbors that "cover" different directions
// around the query rather than just the M nearest ones, which keeps
// the graph globally connected even on highly clustered data.
//
// Acceptance rule for a candidate c:
//
//	for every already-accepted neighbor r:
//	    distance(c.Vector, r.Vector) >= c.Dist
//
// Read it as: "c is accepted only if no already-accepted point is
// closer to c than the query is." Geometrically, an accepted r that
// is closer to c than the query lies in c's "occluded" region —
// you can already reach c by routing through r, so adding the direct
// edge q→c is redundant.
//
// Preconditions:
//   - candidates is sorted ascending by Dist (closest to query first).
//   - Every candidate.ID exists in idx.nodes.
//   - M > 0.
//
// Returns the IDs of accepted neighbors, in acceptance order, with
// len <= M.
func (idx *Index) selectNeighborsHeuristic(query []float32, candidates []Candidate, M int) []uint32 {
	if M <= 0 || len(candidates) == 0 {
		return nil
	}

	accepted := make([]uint32, 0, M)
	// Cache accepted vectors so we don't re-look-up idx.nodes[id] inside
	// the inner loop. For ef=200 / M=16 this is a small win, but it also
	// reads better.
	acceptedVecs := make([][]float32, 0, M)

	for _, c := range candidates {
		if len(accepted) >= M {
			break
		}

		cNode, ok := idx.nodes[c.ID]
		if !ok {
			// Defensive: a candidate referring to a missing node would
			// only happen via a programming bug elsewhere. Skip it
			// rather than panicking so a stale search doesn't bring
			// down the process.
			continue
		}
		cVec := cNode.Vector

		// Check c against every already-accepted neighbor.
		shouldAccept := true
		for _, rVec := range acceptedVecs {
			if idx.Distance(cVec, rVec) < c.Dist {
				// An accepted point r is closer to c than the query is —
				// reject c as occluded.
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
