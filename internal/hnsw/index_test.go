package hnsw

import (
	"testing"

	"github.com/zhangchuqi1998/velosearch/internal/distance"
)

func TestRandomLevelDistribution(t *testing.T) {
	idx := NewIndex(128, 16, 200, distance.L2Squared)
	counts := make(map[int]int)
	const n = 100_000
	for i := 0; i < n; i++ {
		counts[idx.randomLevel()]++
	}
	p0 := float64(counts[0]) / float64(n)
	if p0 < 0.92 || p0 > 0.96 {
		t.Errorf("layer 0 ratio = %.3f, expected ~0.9375", p0)
	}
	t.Logf("Level distribution: %+v", counts)
}
