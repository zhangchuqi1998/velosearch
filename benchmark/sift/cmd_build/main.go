// cmd_build builds an HNSW index over SIFT-1M and reports build time,
// insert throughput, and memory footprint. It does not persist the index;
// Day 6 rebuilds in-process and runs recall + latency benchmarks on top.
//
// Usage (from the velosearch repo root):
//
//	go run ./benchmark/sift/cmd_build
//
// Flags:
//
//	-m int            HNSW M (default 16)
//	-ef-construction  HNSW efConstruction (default 200)
//	-base string      path to sift_base.fvecs (default benchmark/sift/data/sift_base.fvecs)
//	-warmup int       number of inserts to skip when computing warmed-up rate (default 100000)
//	-report-every int log progress every N inserts (default 50000)
package main

import (
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/zhangchuqi1998/velosearch/benchmark/sift"
	"github.com/zhangchuqi1998/velosearch/internal/distance"
	"github.com/zhangchuqi1998/velosearch/internal/hnsw"
)

func main() {
	m := flag.Int("m", 16, "HNSW M")
	efC := flag.Int("ef-construction", 200, "HNSW efConstruction")
	basePath := flag.String("base", "benchmark/sift/data/sift_base.fvecs",
		"path to sift_base.fvecs")
	warmup := flag.Int("warmup", 100_000,
		"inserts to skip when computing warmed-up insert rate")
	reportEvery := flag.Int("report-every", 50_000, "log progress every N inserts")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Printf("Loading %s ...", *basePath)
	tLoad := time.Now()
	base, err := sift.LoadFvecs(*basePath)
	if err != nil {
		log.Fatalf("load failed: %v", err)
	}
	if len(base) == 0 {
		log.Fatal("loaded zero vectors")
	}
	log.Printf("Loaded %d vectors of dim %d in %v",
		len(base), len(base[0]), time.Since(tLoad))

	idx := hnsw.NewIndex(len(base[0]), *m, *efC, distance.L2Squared)
	log.Printf("Building HNSW index: M=%d, efConstruction=%d, vectors=%d",
		*m, *efC, len(base))

	tBuild := time.Now()
	var warmupAt time.Time

	for i, v := range base {
		idx.Insert(uint32(i), v)

		if i+1 == *warmup {
			warmupAt = time.Now()
		}

		if (i+1)%*reportEvery == 0 {
			elapsed := time.Since(tBuild)
			rate := float64(i+1) / elapsed.Seconds()
			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)
			log.Printf("  %d / %d  (%.0f ins/sec, heap=%.2f GB, elapsed %v)",
				i+1, len(base), rate, float64(ms.HeapAlloc)/1e9, elapsed.Round(time.Second))
		}
	}

	buildTime := time.Since(tBuild)

	// Final memory snapshot
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	overallRate := float64(len(base)) / buildTime.Seconds()

	var warmedRate float64
	if !warmupAt.IsZero() {
		warmedTime := time.Since(warmupAt)
		warmedCount := len(base) - *warmup
		warmedRate = float64(warmedCount) / warmedTime.Seconds()
	}

	log.Println()
	log.Println("================================================")
	log.Println(" Build complete")
	log.Println("================================================")
	log.Printf("  Vectors:          %d", len(base))
	log.Printf("  Build time:       %v", buildTime.Round(time.Millisecond))
	log.Printf("  Overall rate:     %.0f inserts/sec", overallRate)
	if warmedRate > 0 {
		log.Printf("  Warmed-up rate:   %.0f inserts/sec (skipped first %d cold inserts)",
			warmedRate, *warmup)
	}
	log.Printf("  HeapAlloc:        %.2f GB (current live heap)", float64(ms.HeapAlloc)/1e9)
	log.Printf("  TotalAlloc:       %.2f GB (cumulative, includes freed)", float64(ms.TotalAlloc)/1e9)
	log.Printf("  Sys:              %.2f GB (total OS memory in use)", float64(ms.Sys)/1e9)
	log.Printf("  NumGC:            %d", ms.NumGC)
}
