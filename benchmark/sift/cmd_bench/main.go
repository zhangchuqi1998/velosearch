// cmd_bench builds the SIFT-1M index, then runs the 10K query set and
// measures Recall@k plus latency percentiles for a list of efSearch values.
//
// Building the index takes ~20-25 min on the unoptimized v0.1 code, so the
// runner accepts a comma-separated list of efSearch values via -ef and runs
// them all against the same in-memory index. This avoids rebuilding for
// each parameter point.
//
// Usage from the repo root:
//
//	go run ./benchmark/sift/cmd_bench
//	go run ./benchmark/sift/cmd_bench -ef 100
//	go run ./benchmark/sift/cmd_bench -ef 30,50,75,100,150,200,300,400
//
// Flags:
//
//	-ef string            comma-separated efSearch values (default "50,100,200,400")
//	-k int                top-k neighbors (default 10, must be <= 100)
//	-m int                HNSW M for build (default 16)
//	-ef-construction int  HNSW efConstruction (default 200)
//	-queries int          number of queries to run (max 10000)
//	-base / -query / -gt  data file paths
//
// Output:
//   per-ef summary lines during the run, plus a final markdown table
//   suitable for pasting into NOTES.md.
//
// Future work: persist the built index to skip rebuilds. The current Index
// has unencodable fields (sync.RWMutex, *rand.Rand, distance.DistanceFunc),
// so gob.Encode does not work without custom MarshalBinary methods. WAL on
// Day 10 will replace this gap.
package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/zhangchuqi1998/velosearch/benchmark/sift"
	"github.com/zhangchuqi1998/velosearch/internal/distance"
	"github.com/zhangchuqi1998/velosearch/internal/hnsw"
)

type result struct {
	ef                  int
	recall              float64
	p50, p95, p99, mean time.Duration
	wall                time.Duration
	qps                 float64
}

func main() {
	efFlag := flag.String("ef", "50,100,200,400",
		"comma-separated efSearch values to benchmark")
	k := flag.Int("k", 10, "top-k neighbors")
	m := flag.Int("m", 16, "HNSW M (build-time)")
	efC := flag.Int("ef-construction", 200, "HNSW efConstruction")
	nQueries := flag.Int("queries", 10_000, "number of queries to run (max 10000)")
	basePath := flag.String("base", "benchmark/sift/data/sift_base.fvecs",
		"path to sift_base.fvecs")
	queryPath := flag.String("query", "benchmark/sift/data/sift_query.fvecs",
		"path to sift_query.fvecs")
	gtPath := flag.String("gt", "benchmark/sift/data/sift_groundtruth.ivecs",
		"path to sift_groundtruth.ivecs")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	efList, err := parseEfList(*efFlag, *k)
	if err != nil {
		log.Fatal(err)
	}
	if *k > 100 {
		log.Fatalf("k (%d) must be <= 100 (ground truth only has 100 neighbors per query)", *k)
	}

	// Phase 1: build index
	log.Printf("Loading %s ...", *basePath)
	base, err := sift.LoadFvecs(*basePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d base vectors of dim %d", len(base), len(base[0]))

	idx := hnsw.NewIndex(len(base[0]), *m, *efC, distance.L2Squared)
	log.Printf("Building HNSW (M=%d, efC=%d, vectors=%d) ...", *m, *efC, len(base))
	tBuild := time.Now()
	for i, v := range base {
		idx.Insert(uint32(i), v)
		if (i+1)%100_000 == 0 {
			log.Printf("  built %d/%d (%v elapsed)",
				i+1, len(base), time.Since(tBuild).Round(time.Second))
		}
	}
	buildTime := time.Since(tBuild)
	log.Printf("Build done in %v", buildTime.Round(time.Second))

	// Phase 2: load queries + ground truth
	queries, err := sift.LoadFvecs(*queryPath)
	if err != nil {
		log.Fatal(err)
	}
	gt, err := sift.LoadIvecs(*gtPath)
	if err != nil {
		log.Fatal(err)
	}
	if len(queries) != len(gt) {
		log.Fatalf("mismatch: %d queries but %d ground-truth rows", len(queries), len(gt))
	}
	if *nQueries > len(queries) {
		*nQueries = len(queries)
	}
	log.Printf("Loaded %d queries, using %d of them", len(queries), *nQueries)

	// Phase 3: sweep efSearch
	results := make([]result, 0, len(efList))
	for _, ef := range efList {
		log.Printf("Running ef=%d ...", ef)
		r := benchOne(idx, queries[:*nQueries], gt[:*nQueries], ef, *k)
		results = append(results, r)
		log.Printf("  ef=%d  Recall@%d=%.4f  P50=%v  P95=%v  P99=%v  mean=%v  QPS=%.0f",
			ef, *k, r.recall,
			r.p50.Round(time.Microsecond),
			r.p95.Round(time.Microsecond),
			r.p99.Round(time.Microsecond),
			r.mean.Round(time.Microsecond),
			r.qps)
	}

	// Phase 4: final summary
	printSummary(results, *k, buildTime)
}

// parseEfList parses "50,100,200" into a sorted unique slice; verifies each
// value is >= k.
func parseEfList(s string, k int) ([]int, error) {
	parts := strings.Split(s, ",")
	out := make([]int, 0, len(parts))
	seen := make(map[int]bool, len(parts))
	for _, p := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, fmt.Errorf("invalid ef value %q: %w", p, err)
		}
		if v < k {
			return nil, fmt.Errorf("efSearch (%d) must be >= k (%d)", v, k)
		}
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	slices.Sort(out)
	return out, nil
}

func benchOne(idx *hnsw.Index, queries [][]float32, gt [][]int32, ef, k int) result {
	latencies := make([]time.Duration, len(queries))
	totalRecall := 0.0

	tStart := time.Now()
	for q := range queries {
		t0 := time.Now()
		got := idx.Search(queries[q], k, ef)
		latencies[q] = time.Since(t0)

		// Recall@k against first k ground-truth IDs.
		truthSet := make(map[uint32]bool, k)
		for i := 0; i < k && i < len(gt[q]); i++ {
			truthSet[uint32(gt[q][i])] = true
		}
		hit := 0
		for _, c := range got {
			if truthSet[c.ID] {
				hit++
			}
		}
		totalRecall += float64(hit) / float64(k)
	}
	wall := time.Since(tStart)

	slices.Sort(latencies)
	avg := totalRecall / float64(len(queries))

	var sum time.Duration
	for _, d := range latencies {
		sum += d
	}
	mean := sum / time.Duration(len(latencies))

	pct := func(p float64) time.Duration {
		i := int(float64(len(latencies)-1) * p)
		return latencies[i]
	}

	return result{
		ef:     ef,
		recall: avg,
		p50:    pct(0.50),
		p95:    pct(0.95),
		p99:    pct(0.99),
		mean:   mean,
		wall:   wall,
		qps:    float64(len(queries)) / wall.Seconds(),
	}
}

func printSummary(results []result, k int, buildTime time.Duration) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println(" SIFT-1M benchmark summary")
	fmt.Println("================================================================")
	fmt.Printf("Build time:  %v\n", buildTime.Round(time.Second))
	fmt.Printf("HeapAlloc:   %.2f GB\n", float64(ms.HeapAlloc)/1e9)
	fmt.Printf("Sys (RSS):   %.2f GB\n", float64(ms.Sys)/1e9)
	fmt.Println()
	fmt.Println("Markdown table (paste into NOTES.md):")
	fmt.Println()
	fmt.Printf("| efSearch | Recall@%d | Mean (us) | P50 (us) | P95 (us) | P99 (us) | QPS |\n", k)
	fmt.Println("|----------|-----------|-----------|----------|----------|----------|-----|")
	for _, r := range results {
		fmt.Printf("| %d | %.4f | %.0f | %.0f | %.0f | %.0f | %.0f |\n",
			r.ef, r.recall,
			float64(r.mean.Microseconds()),
			float64(r.p50.Microseconds()),
			float64(r.p95.Microseconds()),
			float64(r.p99.Microseconds()),
			r.qps)
	}
}
