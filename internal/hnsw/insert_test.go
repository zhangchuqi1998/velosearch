package hnsw

import (
	"math/rand"
	"testing"

	"github.com/zhangchuqi1998/velosearch/internal/distance"
)

// TestInsert1000RandomVectors inserts 1000 random 128-d vectors and asserts
// structural invariants of the resulting HNSW graph.
//
// Note: HNSW edges are NOT guaranteed to be bidirectional after Algorithm 1
// line 13 pruning shrinks an over-full neighbor's connection list. This is
// intentional (paper section 4.1, and matches hnswlib's behavior). We check
// weaker but correct invariants instead.
func TestInsert1000RandomVectors(t *testing.T) {
	const (
		N = 1000
		D = 128
	)

	idx := NewIndex(D, 16, 200, distance.L2Squared)
	rng := rand.New(rand.NewSource(42))

	for i := uint32(0); i < N; i++ {
		v := make([]float32, D)
		for j := range v {
			v[j] = rng.Float32()
		}
		idx.Insert(i, v)
	}

	if len(idx.nodes) != N {
		t.Fatalf("nodes count = %d, want %d", len(idx.nodes), N)
	}

	// Invariant 1: every neighbor ID points to a real node, no self-loops,
	// neighbor count respects per-layer caps.
	for id, n := range idx.nodes {
		for layer, neighbors := range n.Neighbors {
			lim := idx.MaxM
			if layer == 0 {
				lim = idx.MaxM0
			}
			if len(neighbors) > lim {
				t.Errorf("node %d layer %d has %d neighbors, limit %d",
					id, layer, len(neighbors), lim)
			}
			seen := make(map[uint32]bool, len(neighbors))
			for _, nid := range neighbors {
				if nid == id {
					t.Errorf("node %d has self-loop at layer %d", id, layer)
				}
				if _, ok := idx.nodes[nid]; !ok {
					t.Errorf("node %d references missing neighbor %d at layer %d", id, nid, layer)
				}
				if seen[nid] {
					t.Errorf("node %d has duplicate neighbor %d at layer %d", id, nid, layer)
				}
				seen[nid] = true
			}
		}
	}

	// Invariant 2: graph at layer 0 is well-connected. BFS from the entry
	// point should reach almost every node (HNSW guarantees a connected
	// layer-0 graph in practice; we allow a small slack for the unlikely
	// case that pruning isolated a node).
	visited := bfsLayer0(idx, idx.entryPoint)
	reachRatio := float64(len(visited)) / float64(N)
	if reachRatio < 0.99 {
		t.Errorf("BFS at layer 0 reached only %d / %d nodes (%.2f%%), expected >= 99%%",
			len(visited), N, reachRatio*100)
	}

	// Invariant 3: layers thin out roughly geometrically. Not a strict
	// assertion — just log so we can spot weirdness.
	perLayer := make(map[int]int)
	for _, n := range idx.nodes {
		for l := 0; l <= n.Level; l++ {
			perLayer[l]++
		}
	}
	t.Logf("index: %d nodes, maxLevel=%d, entryPoint=%d, layer 0 BFS reach=%d, per-layer counts=%v",
		len(idx.nodes), idx.maxLevel, idx.entryPoint, len(visited), perLayer)
}

// bfsLayer0 returns the set of node IDs reachable from start via layer-0
// edges only.
func bfsLayer0(idx *Index, start uint32) map[uint32]bool {
	visited := make(map[uint32]bool)
	queue := []uint32{start}
	visited[start] = true
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		n, ok := idx.nodes[id]
		if !ok || len(n.Neighbors) == 0 {
			continue
		}
		for _, nid := range n.Neighbors[0] {
			if !visited[nid] {
				visited[nid] = true
				queue = append(queue, nid)
			}
		}
	}
	return visited
}

// TestInsertEmptyIndex_FirstNodeBecomesEntry verifies the "first ever insert"
// path sets the entry point and hasEntry flag.
func TestInsertEmptyIndex_FirstNodeBecomesEntry(t *testing.T) {
	idx := NewIndex(4, 16, 200, distance.L2Squared)
	if idx.hasEntry {
		t.Fatal("new index should not have an entry point")
	}
	idx.Insert(42, []float32{1, 2, 3, 4})
	if !idx.hasEntry {
		t.Fatal("after first insert, hasEntry should be true")
	}
	if idx.entryPoint != 42 {
		t.Errorf("entryPoint = %d, want 42", idx.entryPoint)
	}
	if _, ok := idx.nodes[42]; !ok {
		t.Errorf("node 42 not in idx.nodes")
	}
}
