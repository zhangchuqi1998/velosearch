package hnsw

import (
	"testing"

	"github.com/zhangchuqi1998/velosearch/internal/distance"
)

// TestSearchLayer_HandBuiltGraph exercises searchLayer on a small 2D graph
// constructed by hand (bypassing Insert, which doesn't exist yet on Day 3).
//
// Layout (IDs in parentheses):
//
//	A(1)=(0,0)  B(2)=(1,0)  C(3)=(0,1)  D(4)=(2,0)  E(5)=(0,2)  F(6)=(3,3)
//
// Adjacency at layer 0 only:
//
//	A: [B, C]
//	B: [A, C, D]
//	C: [A, B, E]
//	D: [B, F]
//	E: [C, F]
//	F: [D, E]
//
// The roadmap example uses query=(0.5, 0.5), which puts A, B, C all at the
// same L2 distance (0.5). That makes the heap pop order between equal-Dist
// candidates non-deterministic. We use (0.4, 0.3) instead so the top-3 has
// strictly distinct distances and the expected order is unambiguous:
//
//	A: 0.16 + 0.09 = 0.25
//	B: 0.36 + 0.09 = 0.45
//	C: 0.16 + 0.49 = 0.65
//	D: 2.56 + 0.09 = 2.65
//	E: 0.16 + 2.89 = 3.05
//	F: 6.76 + 7.29 = 14.05
//
// Entry point is F (far from query), so the search must traverse F -> D -> B
// before reaching A and C. This exercises the "expand even when closest is
// far" path: F itself goes into results initially, then gets evicted as
// closer candidates arrive.
func TestSearchLayer_HandBuiltGraph(t *testing.T) {
	idx := NewIndex(2, 4, 10, distance.L2Squared)

	add := func(id uint32, x, y float32, neighbors []uint32) {
		idx.nodes[id] = &Node{
			ID:        id,
			Vector:    []float32{x, y},
			Level:     0,
			Neighbors: [][]uint32{neighbors},
		}
	}
	add(1, 0, 0, []uint32{2, 3})    // A: [B, C]
	add(2, 1, 0, []uint32{1, 3, 4}) // B: [A, C, D]
	add(3, 0, 1, []uint32{1, 2, 5}) // C: [A, B, E]
	add(4, 2, 0, []uint32{2, 6})    // D: [B, F]
	add(5, 0, 2, []uint32{3, 6})    // E: [C, F]
	add(6, 3, 3, []uint32{4, 5})    // F: [D, E]

	idx.entryPoint = 6
	idx.maxLevel = 0
	idx.hasEntry = true

	query := []float32{0.4, 0.3}
	got := idx.searchLayer(query, []uint32{6}, 3, 0)

	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d: %+v", len(got), got)
	}

	wantIDs := []uint32{1, 2, 3} // A, B, C in ascending Dist
	for i, want := range wantIDs {
		if got[i].ID != want {
			t.Errorf("result[%d]: got ID %d (Dist %.4f), want ID %d",
				i, got[i].ID, got[i].Dist, want)
		}
	}

	// Sanity: ensure the output is actually sorted ascending by Dist.
	for i := 1; i < len(got); i++ {
		if got[i].Dist < got[i-1].Dist {
			t.Errorf("results not sorted ascending: got[%d].Dist=%.4f < got[%d].Dist=%.4f",
				i, got[i].Dist, i-1, got[i-1].Dist)
		}
	}
}
