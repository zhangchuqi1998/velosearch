package grpcserver

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// SearchLatency observes the wall-clock duration of every Search RPC,
// partitioned by collection. Exponential buckets from 100µs to ~3.3s
// give 16 buckets that cover the realistic latency range (sub-ms
// in-memory hits up to multi-second worst cases on large GIST queries).
var SearchLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "velosearch_search_latency_seconds",
	Help:    "End-to-end Search RPC latency in seconds, including handler validation and serialization.",
	Buckets: prometheus.ExponentialBuckets(0.0001, 2, 16),
}, []string{"collection"})

// InsertCounter is the cumulative number of vectors accepted by Insert
// (counted after WAL append succeeds, before/during the in-memory write).
var InsertCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "velosearch_inserts_total",
	Help: "Total vectors successfully inserted, per collection.",
}, []string{"collection"})

// DeleteCounter is the cumulative number of IDs that were tombstoned
// (only counts IDs that existed; missing IDs are skipped silently).
var DeleteCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "velosearch_deletes_total",
	Help: "Total IDs successfully tombstoned, per collection.",
}, []string{"collection"})

// CollectionSize is the live count of vectors per collection (includes
// tombstones — tombstones aren't reclaimed until a future rebuild).
// Updated after every Insert and Delete that changes a collection.
var CollectionSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "velosearch_collection_vectors",
	Help: "Current vector count per collection (includes tombstones).",
}, []string{"collection"})
