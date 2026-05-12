package hnsw

import "sync"

// bitset is a sparse bit-set keyed by uint32 IDs.
//
// It is designed for HNSW's "visited" tracking inside searchLayer: each
// call typically touches a few hundred to a few thousand nodes, even when
// the graph holds millions. We therefore:
//
//   - Use one bit per ID (so a 1M-node index fits in 125 KB).
//   - Track set positions in a tiny "set" slice so Reset() runs in
//     O(numSet), not O(numNodes). This keeps the visited set's cost
//     proportional to actual work, not graph size.
//
// Compared to the previous map[uint32]bool, this eliminates per-call
// map allocation (the dominant contributor to TotalAlloc in Day 5/6
// builds: 118 GB cumulative, of which the visited maps were the largest
// share). Pooled via bitsetPool below.
type bitset struct {
	bits []uint64
	set  []uint32
}

// Test reports whether bit i is set.
func (b *bitset) Test(i uint32) bool {
	w := int(i >> 6)
	if w >= len(b.bits) {
		return false
	}
	return b.bits[w]&(1<<(i&63)) != 0
}

// Set marks bit i. Growth is amortized via a doubling strategy so repeated
// Sets at growing IDs do not cause quadratic copy cost.
func (b *bitset) Set(i uint32) {
	w := int(i >> 6)
	if w >= cap(b.bits) {
		newCap := cap(b.bits)*2 + 1
		if w+1 > newCap {
			newCap = w + 1
		}
		newBits := make([]uint64, w+1, newCap)
		copy(newBits, b.bits)
		b.bits = newBits
	} else if w >= len(b.bits) {
		b.bits = b.bits[:w+1]
	}
	b.bits[w] |= 1 << (i & 63)
	b.set = append(b.set, i)
}

// Reset clears all currently-set bits in O(numSet) time. The backing
// storage is preserved so subsequent Sets do not reallocate.
func (b *bitset) Reset() {
	for _, i := range b.set {
		b.bits[int(i>>6)] &^= 1 << (i & 63)
	}
	b.set = b.set[:0]
}

// bitsetPool reuses bitsets across searchLayer calls. Concurrent searches
// each get their own; Insert serializes through idx.mu.Lock so its
// searchLayer chain effectively reuses one bitset.
var bitsetPool = sync.Pool{
	New: func() any { return &bitset{} },
}

func getBitset() *bitset {
	b, _ := bitsetPool.Get().(*bitset)
	b.Reset()
	return b
}

func putBitset(b *bitset) {
	bitsetPool.Put(b)
}
