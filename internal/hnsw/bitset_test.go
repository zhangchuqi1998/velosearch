package hnsw

import "testing"

func TestBitset_BasicSetTest(t *testing.T) {
	var b bitset
	if b.Test(42) {
		t.Error("empty bitset should not contain 42")
	}
	b.Set(42)
	if !b.Test(42) {
		t.Error("after Set(42), Test(42) should be true")
	}
	if b.Test(43) {
		t.Error("Set(42) should not affect bit 43")
	}
}

func TestBitset_Reset(t *testing.T) {
	var b bitset
	for _, i := range []uint32{0, 1, 63, 64, 65, 1000, 999_999} {
		b.Set(i)
	}
	for _, i := range []uint32{0, 1, 63, 64, 65, 1000, 999_999} {
		if !b.Test(i) {
			t.Errorf("Test(%d) should be true after Set", i)
		}
	}
	b.Reset()
	for _, i := range []uint32{0, 1, 63, 64, 65, 1000, 999_999} {
		if b.Test(i) {
			t.Errorf("Test(%d) should be false after Reset", i)
		}
	}
	// Spot-check a few unset positions.
	for _, i := range []uint32{2, 99, 12345} {
		if b.Test(i) {
			t.Errorf("Test(%d) should be false (never set)", i)
		}
	}
}

func TestBitset_GrowsOnDemand(t *testing.T) {
	var b bitset
	b.Set(1_000_000)
	if !b.Test(1_000_000) {
		t.Error("Test should be true for large ID after Set grows storage")
	}
	if b.Test(999_999) {
		t.Error("Test should be false for adjacent unset ID")
	}
}

func TestBitset_PoolRoundTrip(t *testing.T) {
	b1 := getBitset()
	b1.Set(7)
	b1.Set(99)
	putBitset(b1)
	b2 := getBitset()
	// Pool may give us back b1 or a fresh one; either way, getBitset called
	// Reset, so it must be empty.
	if b2.Test(7) || b2.Test(99) {
		t.Error("getBitset must return a freshly-reset bitset")
	}
}

func BenchmarkBitset_SetTest(b *testing.B) {
	var bs bitset
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.Set(uint32(i % 10000))
		_ = bs.Test(uint32((i + 1) % 10000))
		if i%500 == 499 {
			bs.Reset()
		}
	}
}
