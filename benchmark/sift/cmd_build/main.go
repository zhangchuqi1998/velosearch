// cmd_build builds an HNSW index over SIFT-1M and reports build time,
// insert throughput, and memory footprint.
//
// With -workers > 1, inserts run concurrently via a worker pool. Per-Node
// RWMutexes and ID-sorted lock acquisition keep this safe (see insert.go).
//
// Usage (from the velosearch repo root):
//
//	go run ./benchmark/sift/cmd_build                 # single-threaded
//	go run ./benchmark/sift/cmd_build -workers 8      # 8 concurrent inserters
package main

import (
	"flag"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
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
	workers := flag.Int("workers", 1, "number of concurrent insert goroutines")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Printf("Loading %s ...", *basePath)
	tLoad := time.Now()
	base, err := sift.LoadFvecs(*basePath)
	if err != nil {
		log.Fatalf("load failed: %v", err)
	}
	log.Printf("Loaded %d vectors of dim %d in %v",
		len(base), len(base[0]), time.Since(tLoad))

	idx := hnsw.NewIndex(len(base[0]), *m, *efC, distance.L2Squared)
	log.Printf("Building HNSW: M=%d, efConstruction=%d, vectors=%d, workers=%d",
		*m, *efC, len(base), *workers)

	tBuild := time.Now()
	var warmupAt time.Time
	var warmupMu sync.Mutex
	var done atomic.Int64

	if *workers <= 1 {
		runSequential(idx, base, *warmup, *reportEvery, tBuild, &warmupAt)
	} else {
		runConcurrent(idx, base, *warmup, *reportEvery, *workers,
			tBuild, &warmupAt, &warmupMu, &done)
	}

	buildTime := time.Since(tBuild)

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	overallRate := float64(len(base)) / buildTime.Seconds()
	var warmedRate float64
	if !warmupAt.IsZero() {
		warmedTime := time.Since(warmupAt)
		warmedRate = float64(len(base)-*warmup) / warmedTime.Seconds()
	}

	log.Println()
	log.Println("================================================")
	log.Println(" Build complete")
	log.Println("================================================")
	log.Printf("  Vectors:          %d", len(base))
	log.Printf("  Workers:          %d", *workers)
	log.Printf("  Build time:       %v", buildTime.Round(time.Millisecond))
	log.Printf("  Overall rate:     %.0f inserts/sec", overallRate)
	if warmedRate > 0 {
		log.Printf("  Warmed-up rate:   %.0f inserts/sec (skipped first %d)",
			warmedRate, *warmup)
	}
	log.Printf("  HeapAlloc:        %.2f GB (live)", float64(ms.HeapAlloc)/1e9)
	log.Printf("  TotalAlloc:       %.2f GB (cumulative)", float64(ms.TotalAlloc)/1e9)
	log.Printf("  Sys:              %.2f GB (OS RSS)", float64(ms.Sys)/1e9)
	log.Printf("  NumGC:            %d", ms.NumGC)
}

func runSequential(idx *hnsw.Index, base [][]float32, warmup, reportEvery int,
	tBuild time.Time, warmupAt *time.Time) {
	for i, v := range base {
		idx.Insert(uint32(i), v)
		if i+1 == warmup {
			*warmupAt = time.Now()
		}
		if (i+1)%reportEvery == 0 {
			elapsed := time.Since(tBuild)
			rate := float64(i+1) / elapsed.Seconds()
			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)
			log.Printf("  %d / %d  (%.0f ins/sec, heap=%.2f GB, elapsed %v)",
				i+1, len(base), rate, float64(ms.HeapAlloc)/1e9, elapsed.Round(time.Second))
		}
	}
}

func runConcurrent(idx *hnsw.Index, base [][]float32, warmup, reportEvery, workers int,
	tBuild time.Time, warmupAt *time.Time, warmupMu *sync.Mutex,
	done *atomic.Int64) {
	type job struct {
		id uint32
		v  []float32
	}
	ch := make(chan job, workers*8)

	// Reporter: prints progress every reportEvery completions.
	reporterDone := make(chan struct{})
	go func() {
		defer close(reporterDone)
		lastReport := int64(0)
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			n := done.Load()
			if n == int64(len(base)) {
				return
			}
			milestone := (n / int64(reportEvery)) * int64(reportEvery)
			if milestone > lastReport && milestone > 0 {
				lastReport = milestone
				elapsed := time.Since(tBuild)
				rate := float64(milestone) / elapsed.Seconds()
				var ms runtime.MemStats
				runtime.ReadMemStats(&ms)
				log.Printf("  %d / %d  (%.0f ins/sec, heap=%.2f GB, elapsed %v)",
					milestone, len(base), rate, float64(ms.HeapAlloc)/1e9,
					elapsed.Round(time.Second))
			}
		}
	}()

	// Worker pool.
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range ch {
				idx.Insert(j.id, j.v)
				n := done.Add(1)
				if int(n) == warmup {
					warmupMu.Lock()
					if warmupAt.IsZero() {
						*warmupAt = time.Now()
					}
					warmupMu.Unlock()
				}
			}
		}()
	}

	for i, v := range base {
		ch <- job{id: uint32(i), v: v}
	}
	close(ch)
	wg.Wait()
	<-reporterDone
}
