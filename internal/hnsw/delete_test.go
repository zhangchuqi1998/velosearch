package hnsw

import (
	"math/rand"
	"testing"

	"github.com/zhangchuqi1998/velosearch/internal/distance"
)

func TestDelete_FiltersFromResults(t *testing.T) {
	const N = 1000
	const dim = 128
	idx := NewIndex(dim, 16, 200, distance.L2Squared)

	rng := rand.New(rand.NewSource(42))
	for i := 0; i < N; i++ {
		v := make([]float32, dim)
		for j := range v {
			v[j] = rng.Float32()
		}
		idx.Insert(uint32(i), v)
	}

	deleted := make(map[uint32]bool)
	for id := uint32(100); id <= 990; id += 10 {
		if err := idx.Delete(id); err != nil {
			t.Fatalf("Delete(%d): %v", id, err)
		}
		deleted[id] = true
	}

	queryRng := rand.New(rand.NewSource(99))
	leaks := 0
	for q := 0; q < 100; q++ {
		query := make([]float32, dim)
		for j := range query {
			query[j] = queryRng.Float32()
		}
		hits := idx.Search(query, 10, 50)
		for _, c := range hits {
			if deleted[c.ID] {
				leaks++
				t.Errorf("query %d: deleted ID %d leaked into top-10 results", q, c.ID)
			}
		}
	}
	t.Logf("Delete-filter check: %d deleted IDs, 100 queries × k=10, leaks=%d", len(deleted), leaks)
}

func TestDelete_StatsCount(t *testing.T) {
	const N = 50
	idx := NewIndex(8, 8, 50, distance.L2Squared)
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < N; i++ {
		v := make([]float32, 8)
		for j := range v {
			v[j] = rng.Float32()
		}
		idx.Insert(uint32(i), v)
	}
	const wantDeleted = 10
	for id := 0; id < wantDeleted*5; id += 5 {
		if err := idx.Delete(uint32(id)); err != nil {
			t.Fatalf("Delete(%d): %v", id, err)
		}
	}
	st := idx.SnapshotStats()
	if st.NumVectors != N {
		t.Errorf("NumVectors: got %d, want %d (tombstones must not reduce count)", st.NumVectors, N)
	}
	if st.NumDeleted != wantDeleted {
		t.Errorf("NumDeleted: got %d, want %d", st.NumDeleted, wantDeleted)
	}
}
