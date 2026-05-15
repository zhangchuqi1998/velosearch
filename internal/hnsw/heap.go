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

// --- Typed direct push/pop --------------------------------------------------
//
// container/heap.Push/Pop go through the heap.Interface — each call boxes a
// Candidate into an `any`, calls our Push/Pop methods, plus walks the heap
// via Less/Swap/Len through indirect method dispatch. That's a few interface
// calls and one allocation per heap operation, which adds up to thousands of
// allocations per Search at ef=200.
//
// The direct functions below operate on concrete *MinHeap / *MaxHeap, inline
// the sift, and inline the Dist comparison. No boxing, no interface dispatch.
// Behaviour is identical to container/heap's Push/Pop on the corresponding
// heap.Interface, just specialised for Candidate.

func minPush(h *MinHeap, c Candidate) {
	*h = append(*h, c)
	a := *h
	j := len(a) - 1
	for {
		i := (j - 1) / 2
		if i == j || a[j].Dist >= a[i].Dist {
			break
		}
		a[i], a[j] = a[j], a[i]
		j = i
	}
}

func minPop(h *MinHeap) Candidate {
	a := *h
	n := len(a) - 1
	a[0], a[n] = a[n], a[0]
	// sift down on a[:n]
	i := 0
	for {
		l := 2*i + 1
		if l >= n {
			break
		}
		j := l
		if r := l + 1; r < n && a[r].Dist < a[l].Dist {
			j = r
		}
		if a[i].Dist <= a[j].Dist {
			break
		}
		a[i], a[j] = a[j], a[i]
		i = j
	}
	c := a[n]
	a[n] = Candidate{} // GC-safe even if Candidate grows pointer fields later
	*h = a[:n]
	return c
}

func maxPush(h *MaxHeap, c Candidate) {
	*h = append(*h, c)
	a := *h
	j := len(a) - 1
	for {
		i := (j - 1) / 2
		if i == j || a[j].Dist <= a[i].Dist {
			break
		}
		a[i], a[j] = a[j], a[i]
		j = i
	}
}

// maxReplaceTop swaps in c as the new root and sifts down. Equivalent to
// (maxPop + maxPush) when the caller already knows c.Dist < currentRoot.Dist
// — saves one full sift pass. Requires len(h) >= 1.
func maxReplaceTop(h *MaxHeap, c Candidate) {
	a := *h
	a[0] = c
	n := len(a)
	i := 0
	for {
		l := 2*i + 1
		if l >= n {
			break
		}
		j := l
		if r := l + 1; r < n && a[r].Dist > a[l].Dist {
			j = r
		}
		if a[i].Dist >= a[j].Dist {
			break
		}
		a[i], a[j] = a[j], a[i]
		i = j
	}
}

func maxPop(h *MaxHeap) Candidate {
	a := *h
	n := len(a) - 1
	a[0], a[n] = a[n], a[0]
	i := 0
	for {
		l := 2*i + 1
		if l >= n {
			break
		}
		j := l
		if r := l + 1; r < n && a[r].Dist > a[l].Dist {
			j = r
		}
		if a[i].Dist >= a[j].Dist {
			break
		}
		a[i], a[j] = a[j], a[i]
		i = j
	}
	c := a[n]
	a[n] = Candidate{}
	*h = a[:n]
	return c
}
