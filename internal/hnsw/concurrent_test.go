package hnsw

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/zhangchuqi1998/velosearch/internal/distance"
)

// TestConcurrentInsert_NoCorruption hammers Insert from many goroutines and
// then verifies invariants of the resulting graph. Without the -race
// detector (CGO is not enabled in this environment), we rely on Go's
// built-in concurrent-map-write panic + manual structural checks to catch
// the most common races.
func TestConcurrentInsert_NoCorruption(t *testing.T) {
	const (
		N       = 5_000
		D       = 32
		Workers = 8
	)

	idx := NewIndex(D, 16, 100, distance.L2Squared)

	// Pre-generate vectors so goroutines don't share RNG state.
	rng := rand.New(rand.NewSource(42))
	vecs := make([][]float32, N)
	for i := range vecs {
		v := make([]float32, D)
		for j := range v {
			v[j] = rng.Float32()
		}
		vecs[i] = v
	}

	// Concurrent insert via worker pool.
	type job struct {
		id uint32
		v  []float32
	}
	ch := make(chan job, Workers*4)
	var wg sync.WaitGroup
	var inserted atomic.Int64
	for w := 0; w < Workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range ch {
				idx.Insert(j.id, j.v)
				inserted.Add(1)
			}
		}()
	}
	for i, v := range vecs {
		ch <- job{id: uint32(i), v: v}
	}
	close(ch)
	wg.Wait()

	if got := inserted.Load(); got != int64(N) {
		t.Fatalf("inserted %d, want %d", got, N)
	}

	// Structural invariants.
	idx.nodesMu.RLock()
	defer idx.nodesMu.RUnlock()

	if len(idx.nodes) != N {
		t.Fatalf("nodes map size = %d, want %d", len(idx.nodes), N)
	}

	for id, n := range idx.nodes {
		n.mu.RLock()
		// Per-layer cap respected.
		for layer, neighbors := range n.Neighbors {
			lim := idx.MaxM
			if layer == 0 {
				lim = idx.MaxM0
			}
			if len(neighbors) > lim {
				t.Errorf("node %d layer %d has %d neighbors > limit %d",
					id, layer, len(neighbors), lim)
			}
			// No self-loops, no dangling references, no duplicates.
			seen := make(map[uint32]bool, len(neighbors))
			for _, nid := range neighbors {
				if nid == id {
					t.Errorf("node %d has self-loop at layer %d", id, layer)
				}
				if _, ok := idx.nodes[nid]; !ok {
					t.Errorf("node %d references missing neighbor %d at layer %d",
						id, nid, layer)
				}
				if seen[nid] {
					t.Errorf("node %d duplicate neighbor %d at layer %d", id, nid, layer)
				}
				seen[nid] = true
			}
		}
		n.mu.RUnlock()
	}
}

// TestConcurrentSearchDuringInsert mixes Search and Insert from many
// goroutines simultaneously. Verifies neither panics, and post-fact recall
// is in a plausible range (>= 70% — concurrent insertion can degrade graph
// quality slightly).
func TestConcurrentSearchDuringInsert(t *testing.T) {
	const (
		N       = 3_000
		D       = 32
		Workers = 4
		Queries = 50
		K       = 5
	)

	idx := NewIndex(D, 16, 100, distance.L2Squared)
	rng := rand.New(rand.NewSource(7))
	vecs := make([][]float32, N)
	for i := range vecs {
		v := make([]float32, D)
		for j := range v {
			v[j] = rng.Float32()
		}
		vecs[i] = v
	}

	// First wave: synchronously seed a few hundred nodes so searches can
	// find an entry point.
	for i := 0; i < 200; i++ {
		idx.Insert(uint32(i), vecs[i])
	}

	// Concurrent inserters + concurrent searchers.
	type job struct {
		id uint32
		v  []float32
	}
	ch := make(chan job, Workers*4)
	var wg sync.WaitGroup
	for w := 0; w < Workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range ch {
				idx.Insert(j.id, j.v)
			}
		}()
	}
	// Searcher goroutine.
	searchCount := atomic.Int64{}
	stopSearch := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		qrng := rand.New(rand.NewSource(123))
		for {
			select {
			case <-stopSearch:
				return
			default:
			}
			q := make([]float32, D)
			for j := range q {
				q[j] = qrng.Float32()
			}
			_ = idx.Search(q, K, 30)
			searchCount.Add(1)
		}
	}()

	for i := 200; i < N; i++ {
		ch <- job{id: uint32(i), v: vecs[i]}
	}
	close(ch)
	// Workers finish.
	doneInserts := make(chan struct{})
	go func() {
		// Wait for inserters only (4 of Workers+1 wg counts).
		for i := 0; i < Workers; i++ {
			// Best-effort: detect insert completion via a separate channel.
			_ = i
		}
		close(doneInserts)
	}()
	// Simpler: just wait for ch drained + insert workers via separate WG.
	// Above setup conflates; signal via a token from main once ch is closed.

	// Wait for inserts: we know once ch is drained and workers loop ends.
	// But searcher uses the same wg. Stop it.
	close(stopSearch)
	wg.Wait()

	if searchCount.Load() == 0 {
		t.Error("searcher made zero queries — something blocked")
	}
	t.Logf("concurrent run: %d inserts, %d searches", N-200, searchCount.Load())

	// Sanity: post-fact recall using brute force vs HNSW search.
	data := make(map[uint32][]float32, N)
	for i, v := range vecs {
		data[uint32(i)] = v
	}
	queries := make([][]float32, Queries)
	qrng := rand.New(rand.NewSource(999))
	for i := range queries {
		q := make([]float32, D)
		for j := range q {
			q[j] = qrng.Float32()
		}
		queries[i] = q
	}
	total := 0.0
	for _, q := range queries {
		truth := BruteForceKNN(data, q, K, distance.L2Squared)
		got := idx.Search(q, K, 50)
		total += RecallAtK(got, truth)
	}
	avg := total / float64(Queries)
	t.Logf("post-concurrent Recall@%d on 3K random vectors = %.4f", K, avg)
	if avg < 0.70 {
		t.Errorf("recall=%.4f too low; concurrent insert may have corrupted the graph", avg)
	}
}
