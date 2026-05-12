package hnsw

// Min/max heaps over Candidate values for use in HNSW search.
// The heaps store concrete []Candidate slices rather than interface{}
// elements, so Push/Pop avoid boxing allocations on the hot search path.

// Candidate is an (ID, Dist) pair used throughout HNSW search.
// Heaps order candidates by Dist; ID is the application-level node id.
type Candidate struct {
	ID   uint32
	Dist float32
}

// MinHeap pops the candidate with the smallest Dist first.
// Used as the "frontier" / candidate set in SEARCH-LAYER: we always
// want to expand the closest-to-query unvisited node next.
//
// Implements heap.Interface. Use with container/heap:
//
//	h := &MinHeap{}
//	heap.Push(h, Candidate{ID: 1, Dist: 0.5})
//	c := heap.Pop(h).(Candidate)
type MinHeap []Candidate

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].Dist < h[j].Dist }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push appends x to the heap. The heap package guarantees x is a Candidate
// here; the type assertion is checked at compile time by container/heap's
// generic-ish interface, so a wrong type would panic at runtime — fine,
// it's a programmer bug, not a runtime error.
func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(Candidate))
}

// Pop removes and returns the last element (which container/heap has
// already swapped into place as the min). Standard idiom from the
// container/heap docs.
func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	c := old[n-1]
	// Zero the removed slot so the GC can reclaim any pointer fields.
	// Candidate has no pointers today, but this keeps the code safe if
	// the struct grows later.
	old[n-1] = Candidate{}
	*h = old[:n-1]
	return c
}

// Peek returns the minimum without removing it. Panics if empty.
// Useful in searchLayer's termination check.
func (h MinHeap) Peek() Candidate {
	return h[0]
}

// MaxHeap pops the candidate with the largest Dist first.
// Used as the "results" set W (size capped at ef): when W is full
// we compare the new candidate against the current furthest, and
// only insert if it's closer.
type MaxHeap []Candidate

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].Dist > h[j].Dist } // reversed
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x any) {
	*h = append(*h, x.(Candidate))
}

func (h *MaxHeap) Pop() any {
	old := *h
	n := len(old)
	c := old[n-1]
	old[n-1] = Candidate{}
	*h = old[:n-1]
	return c
}

// Peek returns the maximum without removing it. Panics if empty.
func (h MaxHeap) Peek() Candidate {
	return h[0]
}
