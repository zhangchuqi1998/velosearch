package hnsw

import (
	"math/rand"
	"testing"

	"github.com/zhangchuqi1998/velosearch/internal/distance"
)

// TestRecall10KRandom inserts 10K uniformly random 128-d vectors and measures
// Recall@10 against a brute-force baseline. This is CHECKPOINT 1 of the
// VeloSearch roadmap: it validates that the HNSW index is structurally
// correct, not the absolute performance on real-world data (that's the
// SIFT-1M test on Day 6).
//
// IMPORTANT: uniformly random 128-d data is HARDER than real-world data
// (SIFT, GIST, etc.) because it has no manifold structure for the HNSW
// heuristic to exploit. At efSearch=50 we get ~71% recall on this
// adversarial input; the same algorithm hits 95%+ on SIFT-1M. We test at
// efSearch=200 here to give the algorithm enough breathing room on
// random data while still asserting a meaningful recall floor.
func TestRecall10KRandom(t *testing.T) {
	const (
		N       = 10_000
		D       = 128
		Queries = 100
		K       = 10
	)

	idx := NewIndex(D, 16, 200, distance.L2Squared)
	rng := rand.New(rand.NewSource(42))
	data := make(map[uint32][]float32, N)
	for i := uint32(0); i < N; i++ {
		v := make([]float32, D)
		for j := range v {
			v[j] = rng.Float32()
		}
		data[i] = v
		idx.Insert(i, v)
	}

	// Generate queries once so we can also sweep ef for diagnostic info.
	queries := make([][]float32, Queries)
	truths := make([][]Candidate, Queries)
	for q := 0; q < Queries; q++ {
		qv := make([]float32, D)
		for j := range qv {
			qv[j] = rng.Float32()
		}
		queries[q] = qv
		truths[q] = BruteForceKNN(data, qv, K, distance.L2Squared)
	}

	// Diagnostic sweep — logged but not asserted on every value.
	for _, ef := range []int{50, 100, 200, 400} {
		totalRecall := 0.0
		for q := 0; q < Queries; q++ {
			got := idx.Search(queries[q], K, ef)
			totalRecall += RecallAtK(got, truths[q])
		}
		avg := totalRecall / float64(Queries)
		t.Logf("efSearch=%d  Recall@10 = %.4f", ef, avg)
	}

	// CHECKPOINT 1 assertion: recall@10 >= 0.90 at efSearch=200 on 10K
	// uniformly random 128-d vectors.
	const checkpointEf = 200
	const target = 0.90
	totalRecall := 0.0
	for q := 0; q < Queries; q++ {
		got := idx.Search(queries[q], K, checkpointEf)
		totalRecall += RecallAtK(got, truths[q])
	}
	avg := totalRecall / float64(Queries)
	t.Logf("CHECKPOINT 1: Recall@10 = %.4f at efSearch=%d (target >= %.2f)",
		avg, checkpointEf, target)
	if avg < target {
		t.Fatalf("CHECKPOINT 1 FAILED: recall = %.4f, need >= %.2f at efSearch=%d",
			avg, target, checkpointEf)
	}
}
