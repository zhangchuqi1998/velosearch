package hnsw

import (
	"container/heap"
	"testing"
)

func TestMinHeapPopOrder(t *testing.T) {
	h := &MinHeap{}
	heap.Push(h, Candidate{ID: 1, Dist: 3.0}) // A
	heap.Push(h, Candidate{ID: 2, Dist: 1.0}) // B
	heap.Push(h, Candidate{ID: 3, Dist: 2.0}) // C

	want := []uint32{2, 3, 1} // B, C, A — smallest Dist first
	for i, wantID := range want {
		got := heap.Pop(h).(Candidate)
		if got.ID != wantID {
			t.Errorf("MinHeap pop[%d]: got ID=%d (Dist=%v), want ID=%d",
				i, got.ID, got.Dist, wantID)
		}
	}
	if h.Len() != 0 {
		t.Errorf("MinHeap not empty after draining: Len=%d", h.Len())
	}
}

func TestMaxHeapPopOrder(t *testing.T) {
	h := &MaxHeap{}
	heap.Push(h, Candidate{ID: 1, Dist: 3.0}) // A
	heap.Push(h, Candidate{ID: 2, Dist: 1.0}) // B
	heap.Push(h, Candidate{ID: 3, Dist: 2.0}) // C

	want := []uint32{1, 3, 2} // A, C, B — largest Dist first
	for i, wantID := range want {
		got := heap.Pop(h).(Candidate)
		if got.ID != wantID {
			t.Errorf("MaxHeap pop[%d]: got ID=%d (Dist=%v), want ID=%d",
				i, got.ID, got.Dist, wantID)
		}
	}
	if h.Len() != 0 {
		t.Errorf("MaxHeap not empty after draining: Len=%d", h.Len())
	}
}

func TestMinHeapPeek(t *testing.T) {
	h := &MinHeap{}
	heap.Push(h, Candidate{ID: 1, Dist: 3.0})
	heap.Push(h, Candidate{ID: 2, Dist: 1.0})
	heap.Push(h, Candidate{ID: 3, Dist: 2.0})

	if got := h.Peek(); got.ID != 2 {
		t.Errorf("MinHeap.Peek: got ID=%d, want 2", got.ID)
	}
	if h.Len() != 3 {
		t.Errorf("Peek should not remove: Len=%d, want 3", h.Len())
	}
}

func TestMaxHeapPeek(t *testing.T) {
	h := &MaxHeap{}
	heap.Push(h, Candidate{ID: 1, Dist: 3.0})
	heap.Push(h, Candidate{ID: 2, Dist: 1.0})
	heap.Push(h, Candidate{ID: 3, Dist: 2.0})

	if got := h.Peek(); got.ID != 1 {
		t.Errorf("MaxHeap.Peek: got ID=%d, want 1", got.ID)
	}
	if h.Len() != 3 {
		t.Errorf("Peek should not remove: Len=%d, want 3", h.Len())
	}
}

// TestSliceLiteralInit verifies you can also seed a heap from an existing slice
// via heap.Init — useful when batch-loading candidates.
func TestSliceLiteralInit(t *testing.T) {
	h := &MinHeap{
		{ID: 1, Dist: 3.0},
		{ID: 2, Dist: 1.0},
		{ID: 3, Dist: 2.0},
	}
	heap.Init(h)

	if got := heap.Pop(h).(Candidate); got.ID != 2 {
		t.Errorf("after Init, first Pop: got ID=%d, want 2", got.ID)
	}
}
