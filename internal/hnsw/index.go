package hnsw

import (
	"math"
	"math/rand"
	"sync"

	"github.com/zhangchuqi1998/velosearch/internal/distance"
)

type Node struct {
	ID        uint32
	Vector    []float32
	Level     int        // 出现在 0..Level
	Neighbors [][]uint32 // Neighbors[layer] = neighbor node IDs at that layer
	Deleted   bool
}

type Index struct {
	Dim            int
	M              int     // typical 16
	MaxM           int     // = M (for layers > 0)
	MaxM0          int     // = 2*M (for layer 0)
	EfConstruction int     // typical 200
	ML             float64 // 1.0 / ln(M)
	Distance       distance.DistanceFunc

	mu         sync.RWMutex
	nodes      map[uint32]*Node
	entryPoint uint32
	maxLevel   int
	hasEntry   bool

	rng *rand.Rand
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

// randomLevel geometric distribution for layer selection. Layer 0 takes ~93.75% (M=16),
// layer k probability ≈ 1/M^k.
func (idx *Index) randomLevel() int {
	r := idx.rng.Float64()
	return int(math.Floor(-math.Log(r) * idx.ML))
}
