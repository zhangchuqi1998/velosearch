package hnsw

import (
	"math/rand"
	"testing"

	"github.com/zhangchuqi1998/velosearch/internal/distance"
)

// BenchmarkSearch_10K builds a 10K random 128-d index once and benches
// Search at the same recall regime ann-benchmarks probes (k=10, ef=200).
// Use to compare against faiss-hnsw and measure Go-side optimizations.
func BenchmarkSearch_10K_ef200(b *testing.B) {
	const N, dim = 10_000, 128
	idx := NewIndex(dim, 16, 200, distance.L2Squared)
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < N; i++ {
		v := make([]float32, dim)
		for j := range v {
			v[j] = rng.Float32()
		}
		idx.Insert(uint32(i), v)
	}

	queries := make([][]float32, 1000)
	qrng := rand.New(rand.NewSource(99))
	for i := range queries {
		q := make([]float32, dim)
		for j := range q {
			q[j] = qrng.Float32()
		}
		queries[i] = q
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = idx.Search(queries[i%len(queries)], 10, 200)
	}
}
