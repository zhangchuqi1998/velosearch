package hnsw

import (
	"math"
	"math/rand"
	"sync"

	"github.com/zhangchuqi1998/velosearch/internal/distance"
)

type Node struct {
	ID    uint32
	Vector []float32 // immutable after Insert returns; safe to read without lock
	Level int        // immutable after Insert returns

	mu        sync.RWMutex // protects Neighbors and Deleted
	Neighbors [][]uint32   // Neighbors[layer] = neighbor IDs at that layer
	Deleted   bool         // tombstone marker (Day 11)
}

type Index struct {
	Dim            int
	M              int
	MaxM           int
	MaxM0          int
	EfConstruction int
	ML             float64
	Distance       distance.DistanceFunc

	nodesMu sync.RWMutex
	nodes   map[uint32]*Node

	globalMu   sync.Mutex
	entryPoint uint32
	maxLevel   int
	hasEntry   bool

	rngMu sync.Mutex
	rng   *rand.Rand
}

func NewIndex(dim, M, efConstruction int, dist distance.DistanceFunc) *Index {
	return &Index{
		Dim:            dim,
		M:              M,
		MaxM:           M,
		MaxM0:          2 * M,
		EfConstruction: efConstruction,
		ML:             1.0 / math.Log(float64(M)),
		Distance:       dist,
		nodes:          make(map[uint32]*Node),
		rng:            rand.New(rand.NewSource(42)),
	}
}

// getNode looks up a node by ID under a read lock on the nodes map.
// Returns nil if the node is missing. Holding the returned pointer is safe;
// the Index only adds to nodes (never deletes), so the pointer remains
// valid for the lifetime of the Index.
func (idx *Index) getNode(id uint32) *Node {
	idx.nodesMu.RLock()
	n := idx.nodes[id]
	idx.nodesMu.RUnlock()
	return n
}

// randomLevel geometric distribution for layer selection. Goroutine-safe.
func (idx *Index) randomLevel() int {
	idx.rngMu.Lock()
	r := idx.rng.Float64()
	idx.rngMu.Unlock()
	return int(math.Floor(-math.Log(r) * idx.ML))
}
