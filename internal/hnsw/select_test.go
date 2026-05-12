package hnsw

import (
	"testing"
)

func TestSelectNeighborsHeuristic_OccludingPoints(t *testing.T) {
	// Candidates laid out on the x-axis with one outlier off-axis:
	//   A=(1.0, 0)   close, points east
	//   B=(1.1, 0)   "behind" A — occluded
	//   C=(1.2, 0)   also "behind" A — occluded
	//   D=(-5, 5)    completely different direction (NW)
	//
	// With M=2 the heuristic should pick A (cheapest) and D
	// (different direction), rejecting B and C as occluded by A.
	//
	// Using L2-squared as the metric throughout (consistent with
	// idx.Distance below):
	//   A.Dist = 1.0    pairwise A-B = 0.01,  A-C = 0.04
	//   B.Dist = 1.21   → 0.01  < 1.21  → reject
	//   C.Dist = 1.44   → 0.04  < 1.44  → reject
	//   D.Dist = 50     → A-D  = 61    → 61 >= 50 → accept
	idx := &Index{
		nodes: map[uint32]*Node{
			1: {ID: 1, Vector: []float32{1.0, 0}},
			2: {ID: 2, Vector: []float32{1.1, 0}},
			3: {ID: 3, Vector: []float32{1.2, 0}},
			4: {ID: 4, Vector: []float32{-5, 5}},
		},
		Distance: l2SquaredForSelectTest,
	}

	query := []float32{0, 0}
	candidates := []Candidate{
		{ID: 1, Dist: 1.0},  // A
		{ID: 2, Dist: 1.21}, // B
		{ID: 3, Dist: 1.44}, // C
		{ID: 4, Dist: 50.0}, // D
	}

	got := idx.selectNeighborsHeuristic(query, candidates, 2)
	want := []uint32{1, 4} // A, D

	if len(got) != len(want) {
		t.Fatalf("selectNeighborsHeuristic returned %d ids, want %d (got=%v)",
			len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("position %d: got ID=%d, want ID=%d (full result=%v)",
				i, got[i], want[i], got)
		}
	}
}

func TestSelectNeighborsHeuristic_EdgeCases(t *testing.T) {
	idx := &Index{
		nodes: map[uint32]*Node{
			1: {ID: 1, Vector: []float32{1, 0}},
			2: {ID: 2, Vector: []float32{2, 0}},
		},
		Distance: l2SquaredForSelectTest,
	}
	query := []float32{0, 0}
	cands := []Candidate{
		{ID: 1, Dist: 1.0},
		{ID: 2, Dist: 4.0},
	}

	t.Run("M=0 returns nil", func(t *testing.T) {
		got := idx.selectNeighborsHeuristic(query, cands, 0)
		if got != nil {
			t.Errorf("M=0: want nil, got %v", got)
		}
	})

	t.Run("M negative returns nil", func(t *testing.T) {
		got := idx.selectNeighborsHeuristic(query, cands, -1)
		if got != nil {
			t.Errorf("M=-1: want nil, got %v", got)
		}
	})

	t.Run("empty candidates returns nil", func(t *testing.T) {
		got := idx.selectNeighborsHeuristic(query, nil, 5)
		if got != nil {
			t.Errorf("empty candidates: want nil, got %v", got)
		}
	})

	t.Run("M larger than candidates", func(t *testing.T) {
		// A=(1,0) dist=1.0; B=(2,0) dist=4.0; A-B=1.0; 1.0 < 4.0 → reject B.
		got := idx.selectNeighborsHeuristic(query, cands, 10)
		if len(got) != 1 || got[0] != 1 {
			t.Errorf("want [1], got %v", got)
		}
	})

	t.Run("M=1 takes only the closest", func(t *testing.T) {
		got := idx.selectNeighborsHeuristic(query, cands, 1)
		if len(got) != 1 || got[0] != 1 {
			t.Errorf("want [1], got %v", got)
		}
	})
}

func TestSelectNeighborsHeuristic_AllAcceptedWhenSpread(t *testing.T) {
	// Three candidates 120° apart on the unit circle — none occludes
	// any other, so all three should be accepted.
	const r = 1.0
	idx := &Index{
		nodes: map[uint32]*Node{
			1: {ID: 1, Vector: []float32{r, 0}},
			2: {ID: 2, Vector: []float32{-0.5 * r, 0.866 * r}},
			3: {ID: 3, Vector: []float32{-0.5 * r, -0.866 * r}},
		},
		Distance: l2SquaredForSelectTest,
	}
	query := []float32{0, 0}
	cands := []Candidate{
		{ID: 1, Dist: 1.0},
		{ID: 2, Dist: 1.0},
		{ID: 3, Dist: 1.0},
	}

	got := idx.selectNeighborsHeuristic(query, cands, 3)
	if len(got) != 3 {
		t.Errorf("expected 3 accepted (all spread out), got %d: %v", len(got), got)
	}
}

// l2SquaredForSelectTest avoids importing the distance package so this
// file compiles standalone in the hnsw package.
func l2SquaredForSelectTest(a, b []float32) float32 {
	var sum float32
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return sum
}
