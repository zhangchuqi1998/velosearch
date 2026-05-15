# VeloSearch — 14 天每日任务清单 / 14-Day Daily Roadmap

> **目标 / Goal**: 从零用 Go 实现 HNSW 向量搜索引擎,跑出可以放上简历的指标。
> Build an HNSW-based vector search engine from scratch in Go, with measurable resume-grade metrics.
>
> **节奏 / Pace**: 4 小时/天 × 14 天 + 1 天收尾 = ~60 小时
>
> **栈 / Stack**: Go 1.22+, gRPC, Protocol Buffers, Docker

---

## 📛 代码块标签说明 / Code Block Tag Legend

文档里每个代码块前面都有一个 tag,告诉你这块代码该怎么处理:

| Tag | 含义 / Meaning | 你该做啥 / Action |
|------|---------|---------|
| 📋 **COPY** | 完整脚手架,直接抄到指定文件 | Ctrl+C → Ctrl+V 到 `<path>`,**不要改** |
| 🧩 **SKELETON** | 框架已写好,留了 `TODO:` / `...` | 填空。逻辑给你了,只补 1-2 段实现 |
| ✍️ **WRITE** | 这部分你写 | 上面通常有 🤖 LLM prompt,复制给 LLM 出 Go 代码 |
| 🤖 **LLM** | 复制给 LLM 的提示词 | 粘贴到 ChatGPT / Claude / Cursor,**不是你的项目代码** |
| ▶️ **RUN** | 终端命令 | 在 PowerShell 跑(或 bash,如果标了) |
| 📚 **PSEUDO** | 算法伪代码或参考表 | 看懂思路,**不能直接编译**,要么自己翻成 Go,要么交给上面的 LLM prompt |
| 🎯 **VERIFY** | 验证命令/期望输出 | 跑这些命令确认上一步没翻车 |

**简单决策树 / Quick decision tree:**
- 看到 📋 → 抄
- 看到 🤖 → 给 LLM
- 看到 ✍️ + 上面有 🤖 → 让 LLM 出代码
- 看到 🧩 → 抄,然后填 TODO
- 看到 📚 → 不能编译,只读思路

---

## 怎么用这份文档 / How to use this doc

1. **每天打开本文件,找到对应 Day** — 顺序执行 Block 1 → Break → Block 2 → Wrap-up
2. **每个 task 有四件套**:
   - 做什么 / What — 一句话目标
   - 怎么做 / How — 编号步骤
   - 验证 / Verify — 怎么知道做对了
   - LLM 提示词 / LLM prompt — 卡住直接复制粘贴
3. **Checkpoint 没过不要往下走** — 跳过 = 简历数字编的
4. **结尾必做**: `git commit` + 在 `WEEKLY_LOG.md` 写三行 (今天完成 / 卡点 / 明天先做)

**LLM 提示原则 / How to prompt LLM well**: 给 (1) 你想做啥 (2) 你试了啥 (3) 报错原文 (4) 最小复现 (30-50 行,别贴整个文件)。

---

## 目标指标 / Target Metrics

| Metric | Target | 在哪测 / Where measured |
|--------|--------|---------|
| Recall@10 on SIFT-1M | ≥ 95% | Day 6, Day 7 |
| P99 search latency | < 15ms | Day 6, Day 7 |
| Insert throughput | ≥ 3K/sec | Day 6 |
| Crash recovery | `kill -9` 后 0 数据丢失 / 0 data loss after kill -9 | Day 11 |

⚠️ **没测出来的数字不要写简历。** 实测低于目标就改简历到实测值 — 22ms 可辩护,15ms 编的会被面试官当场拆。

---

## 关键里程碑 / Key Checkpoints

| ID | Day | 检查项 / Check |
|----|-----|----------------|
| **C1** | Day 4 | Recall@10 ≥ 90% on 10K random 128-d vectors |
| **C2** | Day 7 | Recall ≥ 95% AND P99 < 15ms on SIFT-1M |
| **C3** | Day 11 | `kill -9` 中途服务 → 重启 → 已确认 insert 全部可搜到 |

每个 checkpoint 没过当天就要发现,不要拖。

---

## 通用每日节奏 / Daily rhythm

```
[ 90 min Block 1 ] → [ 20 min break ] → [ 90 min Block 2 ] → [ 20 min wrap-up ]
```

**Block 中不要看微信/X/邮箱**。Break 才看。
**Wrap-up 必做 3 件事**:
1. `git add . && git commit -m "feat(day-N): <一句话>"`
2. 在 `WEEKLY_LOG.md` 写当天 3 行
3. 在脑子里过一遍明天 Block 1 第一个 task

---
---

# Day 1 — 2026-05-11 (周一 / Mon) — 项目搭建 + 读 HNSW 论文 ✅ 已完成

> **你的当前状态:Day 1 Block 1 已经由 Claude 帮你做完。** 工具链装好了,目录结构和脚手架文件都写好了,build/test/lint 都过。
>
> **你还要做的:Day 1 Block 2(读 HNSW 论文)+ Wrap-up(commit + WEEKLY_LOG)。** 详见下面 Block 2。

## Block 1 (90 min) — 装工具 + 搭骨架 ✅ DONE

(略 — 已完成,详情见 git log)

**📋 COPY 已完成的内容(参考)/ Already done (for reference):**
- 装好的工具: Go 1.26.3, protoc 34.1, grpcurl 1.9.3, golangci-lint 2.12.2, protoc-gen-go, protoc-gen-go-grpc, make
- 已建目录: `F:\app\job\velosearch\` + 15 个子目录
- 已写文件: `Makefile`, `.golangci.yml`, `.github/workflows/ci.yml`, `.gitignore`, `cmd/server/main.go`, `README.md`, `NOTES.md`, `WEEKLY_LOG.md`
- `go mod init` + `git init` 已完成

## Block 2 (90 min) — 读 HNSW 论文 / Read HNSW Paper ← 你做这块

### 2.A Malkov & Yashunin 2018 Section 4 (60 min)

**Paper:** https://arxiv.org/abs/1603.09320

**读哪部分:**
- Section 3 (Background, 略读) — 5 min
- **Section 4 (Algorithm Description, 精读)** — 40 min
  - Algorithm 1: INSERT
  - Algorithm 2: SEARCH-LAYER ← 最关键
  - Algorithm 3: SELECT-NEIGHBORS-SIMPLE
  - Algorithm 4: SELECT-NEIGHBORS-HEURISTIC ← 用这个
  - Algorithm 5: K-NN-SEARCH
- Section 5 (Experiments, 略读) — 15 min

**✍️ WRITE: 在 `NOTES.md` 填这张表**(模板已经在文件里)

| Param | 含义 | Typical |
|------|------|---------|
| `M` | ... | 16 |
| `mL` | ... | ~0.36 |
| `efConstruction` | ... | 200 |
| `efSearch` | ... | 50-400 |

### 2.B Pinecone 博客 (20 min)

https://www.pinecone.io/learn/series/faiss/hnsw/

### 2.C LLM 加深理解 (10 min)

**🤖 LLM PROMPT** (复制到 ChatGPT/Claude):
```
Explain HNSW (Hierarchical Navigable Small World) assuming I'm about to implement it
from scratch in Go. Walk through Algorithms 1-5 from Malkov & Yashunin 2018.

For each algorithm:
1. What it does in one sentence
2. What data structures it touches (heaps, visited sets, neighbor lists)
3. Where it's called from
4. THE most common implementation mistake people make

Then explain mL, M, M_max, M_max0, efConstruction, efSearch — what happens when each
parameter is set too high vs. too low.

End with: a 6-node 2D example I can trace by hand to verify my searchLayer implementation.
```

## 🏁 Wrap-up (20 min)

**✍️ WRITE: 在 `NOTES.md` 末尾回答这 3 个问题(口头先说,然后写下来):**
1. 为啥 HNSW 用分层图(不是单层)?
2. 搜索时,什么时候从 greedy descent 切到 ef-bounded search?
3. Algorithm 4 的邻居启发式防止了什么?

**▶️ RUN:**
```powershell
cd F:\app\job\velosearch
git add .
git status
git commit -m "chore(day-1): project skeleton, Makefile, lint config, CI"
# After finishing the paper:
git add NOTES.md
git commit -m "docs(day-1): HNSW paper notes filled"
```

**✍️ WRITE: 在 `WEEKLY_LOG.md` 填 Day 1 那段(模板已经在文件里)**

---
---

# Day 2 — 2026-05-12 (周二 / Tue) — 距离 + 优先队列 + 索引数据结构

## 今日目标
- L2/Cosine 距离 + 单测 + bench
- 类型安全的 MinHeap / MaxHeap
- `Node` / `Index` 数据结构定型

## Block 1 (90 min) — Distance + Priority Queue

### 2.1 距离函数 (~45 min)

**目标文件:** `internal/distance/distance.go`

**🤖 LLM PROMPT** (复制给 LLM):
```
Write Go for internal/distance/distance.go:
- L2Squared(a, b []float32) float32   // no sqrt
- Cosine(a, b []float32) float32      // 1 - dot/(|a|*|b|)
- DotProduct(a, b []float32) float32  (helper)
- Norm(a []float32) float32           (helper)

Both panic if len(a) != len(b).
Add internal/distance/distance_test.go with table tests covering:
  - L2Squared([0,0], [3,4]) == 25
  - Cosine of orthogonal == 1.0
  - Cosine of identical == 0.0 (up to 1e-6)
  - panic on length mismatch
Add benchmarks BenchmarkL2Squared and BenchmarkCosine on 128-d vectors.

Define type DistanceFunc func(a, b []float32) float32 at the top.
```

**✍️ WRITE:** 把 LLM 输出贴到 `internal/distance/distance.go` 和 `distance_test.go`,简单看一遍逻辑正确。

**🎯 VERIFY:**
```powershell
cd F:\app\job\velosearch
go test ./internal/distance/ -v
go test -bench=. -benchmem ./internal/distance/
```
记下 ns/op,Day 7 你会优化这个。

### 2.2 优先队列 (~45 min)

**目标文件:** `internal/hnsw/heap.go`, `internal/hnsw/heap_test.go`

**🤖 LLM PROMPT:**
```
I'm implementing HNSW in Go. I need two heaps over candidate values:

type Candidate struct {
    ID   uint32
    Dist float32
}

Implement using container/heap:
1. MinHeap: pops the candidate with smallest Dist first (used as the "frontier" / candidates set in searchLayer)
2. MaxHeap: pops largest Dist first (used as "results" set, size capped at ef)

Provide:
- type MinHeap []Candidate (with Len, Less, Swap, Push, Pop methods to implement heap.Interface)
- type MaxHeap []Candidate
- A tiny heap_test.go that:
  pushes (3.0, A), (1.0, B), (2.0, C) into a MinHeap and asserts pop order B,C,A
  pushes the same into a MaxHeap and asserts pop order A,C,B

CRITICAL: use the concrete []Candidate type, NOT interface{}. This is a hot path,
GC pressure from boxing will kill latency.
```

**🎯 VERIFY:**
```powershell
go test ./internal/hnsw/ -run TestHeap -v
```

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — Index / Node Data Structures

### 2.3 类型定义 (~30 min)

**目标文件:** `internal/hnsw/index.go`

**📋 COPY → `internal/hnsw/index.go`** (这个完整可用,直接抄):
```go
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
	Level     int        // present on layers 0..Level
	Neighbors [][]uint32 // Neighbors[layer] = neighbor node IDs at that layer
	Deleted   bool       // tombstone marker, used on Day 11
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
```

### 2.4 随机层级生成 (~30 min)

**📋 COPY → `internal/hnsw/index.go`** (追加到上面的文件):
```go
// randomLevel generates a geometric-distribution level. Layer 0 holds ~93.75% (M=16),
// P(layer k) ~= 1/M^k.
func (idx *Index) randomLevel() int {
	r := idx.rng.Float64()
	return int(math.Floor(-math.Log(r) * idx.ML))
}
```

**📋 COPY → `internal/hnsw/index_test.go`** (测试这个分布):
```go
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
```

**🎯 VERIFY:**
```powershell
go test ./internal/hnsw/ -v
```

### 2.5 距离调用快路径 (~30 min)

**📚 PSEUDO** — 现在只 stub,Day 7 才优化。本质上是把 `idx.Distance` 在热循环外捕获到局部变量,帮助 Go 编译器内联:
```go
// in searchLayer (you write this on Day 3):
//   dist := idx.Distance  ← local copy
//   for _, neighborID := range node.Neighbors[layer] {
//       d := dist(query, otherVector)
//       ...
//   }
```

**今天什么都不用做 — Day 3 写 searchLayer 时记住把 `dist := idx.Distance` 提到循环外即可。**

---

## 🏁 Wrap-up (20 min)

**▶️ RUN:**
```powershell
git add .
git commit -m "feat(day-2): distance funcs, typed heaps, Index/Node types, randomLevel"
```

**✍️ WRITE WEEKLY_LOG.md:**
```
## Day 2 — 2026-05-12
- 完成: L2/Cosine + 测试 + bench, MinHeap/MaxHeap, Index/Node 数据结构, randomLevel 分布测试通过
- L2Squared ns/op (baseline): <填>
- 卡点: <如果有>
- 明天: searchLayer (Algorithm 2) — 最难一天,留好精力
```

**Today's checklist:**
- [ ] `go test ./...` 全绿
- [ ] L2Squared bench ns/op 记下来
- [ ] randomLevel 分布测试通过 (~93.75% on layer 0)

---
---

# Day 3 — 2026-05-13 (周三 / Wed) — searchLayer + Insert(核心)

## 今日目标
- 实现 `searchLayer` (Algorithm 2),手工 6 节点图验证
- 实现 `Insert` + 启发式邻居选择
- 1000 个随机向量后图全部双向

⚠️ **今天技术含量最高,留好精力。**

## Block 1 (90 min) — searchLayer

### 3.1 函数签名 + 框架 (30 min)

**🧩 SKELETON → `internal/hnsw/search.go`** (这是骨架,你用 3.2 的 prompt 让 LLM 补完中间循环):
```go
package hnsw

import (
	"container/heap"
	"sort"
)

// searchLayer runs an ef-bounded greedy search at a single layer.
//
// query: the query vector
// entryPoints: starting node IDs at this layer
// ef: candidate set capacity (== 1 for greedy descent, efSearch/efConstruction otherwise)
// layer: which layer to search
//
// Returns results sorted ascending by Dist.
func (idx *Index) searchLayer(query []float32, entryPoints []uint32, ef, layer int) []Candidate {
	visited := make(map[uint32]bool, ef*2)
	candidates := &MinHeap{}
	results := &MaxHeap{}
	dist := idx.Distance

	// Init: push every entry point into both candidates and results
	for _, ep := range entryPoints {
		d := dist(query, idx.nodes[ep].Vector)
		visited[ep] = true
		heap.Push(candidates, Candidate{ep, d})
		heap.Push(results, Candidate{ep, d})
	}

	// TODO: main loop -- ask LLM with the prompt in section 3.2 to fill this in
	// 1. Pop closest unexpanded node c from candidates
	// 2. If c.Dist > furthest in results AND len(results) == ef, break
	// 3. For each unvisited neighbor of c at this layer:
	//    a. compute distance
	//    b. mark visited
	//    c. if closer than furthest result OR results not full -> push to candidates and results
	//    d. if len(results) > ef -> evict furthest

	// Return results sorted ascending by Dist
	out := make([]Candidate, results.Len())
	for i := len(out) - 1; i >= 0; i-- {
		out[i] = heap.Pop(results).(Candidate)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Dist < out[j].Dist })
	return out
}
```

### 3.2 主循环 + 早停 (45 min)

**📚 PSEUDO** — 论文 Algorithm 2 思路:
```
while candidates not empty:
    c = pop min from candidates
    f = top of results (furthest)
    if c.dist > f.dist AND len(results) == ef:
        break
    for each neighbor n of c at layer:
        if n in visited: continue
        visited[n] = true
        d = dist(query, n.vector)
        if d < f.dist OR len(results) < ef:
            push (d, n) into candidates
            push (d, n) into results
            if len(results) > ef: pop max from results
```

**🤖 LLM PROMPT:**
```
Implement HNSW Algorithm 2 (SEARCH-LAYER) in Go. I have this skeleton with the
init done and the TODO marking where the main loop should go:

<paste your 3.1 skeleton from internal/hnsw/search.go>

Available types/methods:
- type Candidate struct { ID uint32; Dist float32 }
- MinHeap and MaxHeap (container/heap implementations of []Candidate)
- idx.nodes is map[uint32]*Node
- Node has fields: ID, Vector []float32, Neighbors [][]uint32 (Neighbors[layer] = neighbor IDs)
- idx.Distance is a DistanceFunc

Fill in the TODO with the main loop:
1. Pop closest from candidates (heap.Pop)
2. Peek results' top (results[0] since it's a max-heap, top is largest)
3. Early termination check
4. Iterate node's neighbors at layer
5. Update visited, push to both heaps, evict from results if over ef

Then add this test in search_test.go:
A hand-built 6-node graph in 2D:
  A=(0,0), B=(1,0), C=(0,1), D=(2,0), E=(0,2), F=(3,3)
  Adjacency at layer 0 only:
    A: [B, C]
    B: [A, C, D]
    C: [A, B, E]
    D: [B, F]
    E: [C, F]
    F: [D, E]
Manually wire this graph (create Nodes with these neighbors, populate idx.nodes).
Then call searchLayer(query=[0.5, 0.5], entryPoints=[F], ef=3, layer=0).
Expected results in order: A, B, C (sorted by L2Squared from [0.5,0.5]).
```

### 3.3 跑测试 + 手工验证 (15 min)

**🎯 VERIFY:**
```powershell
go test ./internal/hnsw/ -run TestSearchLayer -v
```

如果挂了,在纸上画图标出每一步 candidates / results 的状态,找哪步算错了。

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — Insert + Algorithm 4

### 3.4 Algorithm 4 启发式 (30 min)

**目标文件:** `internal/hnsw/select.go`

**🤖 LLM PROMPT:**
```
Implement HNSW Algorithm 4 (SELECT-NEIGHBORS-HEURISTIC) in Go.

Place in internal/hnsw/select.go:

func (idx *Index) selectNeighborsHeuristic(query []float32, candidates []Candidate, M int) []uint32

`candidates` is already sorted ascending by Dist (distance from query).
For each candidate c in order:
  accept c if for ALL already-accepted neighbors r:
    distance(c.Vector, r.Vector) >= c.Dist
  (i.e., no accepted neighbor is closer to c than c is to the query)
Stop when len(accepted) == M.

Return the IDs of accepted neighbors.

Available:
- idx.Distance is a DistanceFunc
- idx.nodes is map[uint32]*Node, where Node has Vector []float32

Add a test in select_test.go:
Hand-construct 4 candidates with this distance pattern:
  A at (1.0, 0)   — dist 1.0
  B at (1.1, 0)   — dist 1.1
  C at (1.2, 0)   — dist 1.2
  D at (5.0, 0)   — dist 5.0
Query at origin.

With M=2:
  Accept A (first one, trivially).
  Reject B (B is at distance 0.1 from A < 1.1).
  Reject C (C is at distance 0.2 from A < 1.2).
  Accept D (D is at distance 4.0 from A >= 5.0 fails... actually 4.0 < 5.0, so accept). Wait:
  rule: accept c if for all accepted r, distance(c, r) >= c.Dist
  D vs A: distance(D, A) = 4.0; c.Dist = 5.0; 4.0 >= 5.0 is FALSE → reject D too.
  Result: [A].

Adjust the test data if needed so we get expected behavior [A, D]:
Try A=(1,0), B=(1.1,0), C=(1.2,0), D=(5,5). Then:
  D vs A: distance = sqrt(16+25) ≈ 6.4; c.Dist for D = sqrt(50) ≈ 7.07; 6.4 >= 7.07 false → still reject.
Use D=(-5, 5): distance(D,A)=sqrt(36+25)≈7.81; D.Dist=sqrt(50)≈7.07; 7.81 >= 7.07 true → accept.
Final test data:
  A=(1,0), B=(1.1,0), C=(1.2,0), D=(-5,5)
  Expected: [A, D] with M=2.
```

### 3.5 Insert 主体 (50 min)

**目标文件:** `internal/hnsw/insert.go`

**📋 COPY → `internal/hnsw/insert.go`** (完整可用,只有 `pruneNeighbors` 是 TODO):
```go
package hnsw

func (idx *Index) Insert(id uint32, vector []float32) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	level := idx.randomLevel()
	node := &Node{
		ID:        id,
		Vector:    vector,
		Level:     level,
		Neighbors: make([][]uint32, level+1),
	}
	idx.nodes[id] = node

	// Empty index: first node becomes the entry point
	if !idx.hasEntry {
		idx.entryPoint = id
		idx.maxLevel = level
		idx.hasEntry = true
		return
	}

	entryPoints := []uint32{idx.entryPoint}

	// Greedy descent from maxLevel down to level+1, ef=1 to find best entry per layer
	for l := idx.maxLevel; l > level; l-- {
		res := idx.searchLayer(vector, entryPoints, 1, l)
		if len(res) > 0 {
			entryPoints = []uint32{res[0].ID}
		}
	}

	// From min(level, maxLevel) down to 0, ef = efConstruction
	start := level
	if idx.maxLevel < start {
		start = idx.maxLevel
	}
	for l := start; l >= 0; l-- {
		candidates := idx.searchLayer(vector, entryPoints, idx.EfConstruction, l)
		Mlayer := idx.MaxM
		if l == 0 {
			Mlayer = idx.MaxM0
		}
		chosen := idx.selectNeighborsHeuristic(vector, candidates, Mlayer)

		// Add bidirectional edges
		node.Neighbors[l] = chosen
		for _, nid := range chosen {
			other := idx.nodes[nid]
			other.Neighbors[l] = append(other.Neighbors[l], id)
			if len(other.Neighbors[l]) > Mlayer {
				idx.pruneNeighbors(other, l, Mlayer)
			}
		}

		// Next layer's entry: all current candidates
		entryPoints = make([]uint32, len(candidates))
		for i, c := range candidates {
			entryPoints[i] = c.ID
		}
	}

	// Update entry point
	if level > idx.maxLevel {
		idx.entryPoint = id
		idx.maxLevel = level
	}
}

// pruneNeighbors trims node n's neighbors at the given layer down to M.
// Reuses selectNeighborsHeuristic with current neighbors as candidates and n.Vector as the query.
// TODO: implement (~10 lines).
func (idx *Index) pruneNeighbors(n *Node, layer, M int) {
	// 1. Build Candidate{ID, dist from n.Vector} for each id in n.Neighbors[layer]
	// 2. Sort candidates by Dist ascending
	// 3. Call idx.selectNeighborsHeuristic(n.Vector, sortedCandidates, M)
	// 4. Assign result to n.Neighbors[layer]
	panic("TODO(day-3): implement pruneNeighbors")
}
```

**✍️ WRITE: 你来实现 `pruneNeighbors` 函数体,~10 行。** 用上面注释里的 4 步,或者让 LLM 出:

**🤖 LLM PROMPT (如果卡住):**
```
Implement pruneNeighbors in Go. Signature:
func (idx *Index) pruneNeighbors(n *Node, layer, M int)

Steps:
1. Build []Candidate from n.Neighbors[layer]:
   for each neighbor ID nid in n.Neighbors[layer]:
     d := idx.Distance(n.Vector, idx.nodes[nid].Vector)
     candidates = append(candidates, Candidate{nid, d})
2. Sort candidates by Dist ascending (sort.Slice)
3. Call kept := idx.selectNeighborsHeuristic(n.Vector, candidates, M)
4. Assign n.Neighbors[layer] = kept

That's it, no edge removal cleanup needed (the dropped neighbors' back-edges to n will be
cleaned up the next time they get pruned, which is acceptable for v0.1).
```

### 3.6 烟雾测试 (10 min)

**📋 COPY → `internal/hnsw/insert_test.go`**:
```go
package hnsw

import (
	"math/rand"
	"testing"

	"github.com/zhangchuqi1998/velosearch/internal/distance"
)

func TestInsert1000RandomVectors(t *testing.T) {
	idx := NewIndex(128, 16, 200, distance.L2Squared)
	rng := rand.New(rand.NewSource(42))
	for i := uint32(0); i < 1000; i++ {
		v := make([]float32, 128)
		for j := range v {
			v[j] = rng.Float32()
		}
		idx.Insert(i, v)
	}

	// Assert all edges are bidirectional
	for id, n := range idx.nodes {
		for layer, neighbors := range n.Neighbors {
			for _, nid := range neighbors {
				other := idx.nodes[nid]
				found := false
				for _, back := range other.Neighbors[layer] {
					if back == id {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("edge %d -> %d at layer %d is not bidirectional", id, nid, layer)
				}
			}
		}
	}
	t.Logf("Index has %d nodes, maxLevel = %d", len(idx.nodes), idx.maxLevel)
}
```

**🎯 VERIFY:**
```powershell
go test ./internal/hnsw/ -v
```

---

## 🏁 Wrap-up (20 min)

**▶️ RUN:**
```powershell
git add .
git commit -m "feat(day-3): searchLayer, Algorithm 4 heuristic, Insert; bidirectional edges verified"
```

**Today's checklist:**
- [ ] 6-node `TestSearchLayer` 通过
- [ ] `TestSelectNeighborsHeuristic` 通过 (expected [A, D])
- [ ] `TestInsert1000RandomVectors` 通过

---
---

# Day 4 — 2026-05-14 (周四 / Thu) — 🎯 CHECKPOINT 1

## 今日目标
- 实现 `Search` (Algorithm 5)
- 实现 `BruteForceKNN` 作 ground truth
- 在 10K 随机向量上跑出 **Recall@10 ≥ 90%**

## Block 1 (90 min) — Search + Brute Force

### 4.1 Search (30 min)

**📋 COPY → `internal/hnsw/search.go`** (追加到已有文件):
```go
// Search is the top-level KNN query.
// k: k: number of nearest neighbors to return
// efSearch: efSearch: candidate set size at query time (runtime tunable)
func (idx *Index) Search(query []float32, k, efSearch int) []Candidate {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if !idx.hasEntry {
		return nil
	}

	entryPoints := []uint32{idx.entryPoint}
	// Greedy descent from maxLevel down to 1
	for l := idx.maxLevel; l > 0; l-- {
		res := idx.searchLayer(query, entryPoints, 1, l)
		if len(res) > 0 {
			entryPoints = []uint32{res[0].ID}
		}
	}

	// Layer 0 with efSearch
	candidates := idx.searchLayer(query, entryPoints, efSearch, 0)

	// Filter tombstones (Day 11 starts using Deleted; this is a no-op until then)
	out := make([]Candidate, 0, k)
	for _, c := range candidates {
		if idx.nodes[c.ID].Deleted {
			continue
		}
		out = append(out, c)
		if len(out) == k {
			break
		}
	}
	return out
}
```

### 4.2 Brute Force baseline (15 min)

**📋 COPY → `internal/hnsw/brute.go`**:
```go
package hnsw

import "sort"

// BruteForceKNN is a brute-force KNN used as ground truth in tests.
func BruteForceKNN(data map[uint32][]float32, query []float32, k int, dist DistanceFunc) []Candidate {
	all := make([]Candidate, 0, len(data))
	for id, v := range data {
		all = append(all, Candidate{id, dist(query, v)})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Dist < all[j].Dist })
	if len(all) > k {
		all = all[:k]
	}
	return all
}

// DistanceFunc alias (re-exported for callers)
type DistanceFunc func(a, b []float32) float32
```

**注:** 如果 `DistanceFunc` 已在 `internal/distance` 包定义,这里别重复,直接 import 用。

### 4.3 Recall helper (15 min)

**📋 COPY → `internal/hnsw/recall.go`**:
```go
package hnsw

// RecallAtK computes the overlap ratio between HNSW results and ground truth.
// Both inputs are length-k slices sorted ascending by Dist.
func RecallAtK(hnsw, truth []Candidate) float64 {
	set := make(map[uint32]bool, len(truth))
	for _, c := range truth {
		set[c.ID] = true
	}
	hit := 0
	for _, c := range hnsw {
		if set[c.ID] {
			hit++
		}
	}
	return float64(hit) / float64(len(truth))
}
```

### 4.4 10K 测试 (30 min)

**📋 COPY → `internal/hnsw/recall_test.go`**:
```go
package hnsw

import (
	"math/rand"
	"testing"

	"github.com/zhangchuqi1998/velosearch/internal/distance"
)

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

	// IMPORTANT: efSearch=50 on uniformly random 128-d data gives only ~71%
	// recall (no manifold structure → curse of dimensionality). Use efSearch=200
	// here. SIFT-1M (Day 6) hits 95%+ at efSearch=50 because real data has
	// structure HNSW can exploit.
	totalRecall := 0.0
	for q := 0; q < Queries; q++ {
		qv := make([]float32, D)
		for j := range qv {
			qv[j] = rng.Float32()
		}
		truth := BruteForceKNN(data, qv, K, distance.L2Squared)
		got := idx.Search(qv, K, 200)
		totalRecall += RecallAtK(got, truth)
	}
	avg := totalRecall / float64(Queries)
	t.Logf("Recall@10 = %.4f (target >= 0.90)", avg)
	if avg < 0.90 {
		t.Fatalf("CHECKPOINT 1 FAILED: recall = %.4f, need >= 0.90", avg)
	}
}
```

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — 跑 + 调

### 4.5 跑测试 (10 min)

**▶️ RUN:**
```powershell
go test ./internal/hnsw/ -run TestRecall10KRandom -v -timeout 5m
```

### 4.6 Recall < 90% 排查 (剩余时间)

**📚 REF — 症状 → 原因表:**

| 症状 | 可能原因 |
|---|---|
| Recall ~ 0 | searchLayer 早停条件反了 |
| Recall ~ 50% | greedy descent 没跑到 layer 0 |
| Recall ~ 80% | 用了 Algorithm 3 而不是 4 |
| Recall ~ 85% | pruneNeighbors 没用启发式 |
| Recall 90-94% | 几乎正常,试着把 efSearch 升到 100 |

**📚 REF — 最常见 5 个 bug:**
1. Min-heap / max-heap 弄反
2. early termination 用了 `>=` 而不是 `>`
3. greedy descent 上一层 entry 用了多个而不是最近 1 个
4. selectNeighborsHeuristic 比较反向 (`<` vs `>=`)
5. 双向边只加了 outgoing 没加 incoming

**🤖 LLM PROMPT (debug 用):**
```
My HNSW implementation hits only X% recall@10 on 10K random 128-d vectors (target: 90%+).

Here's my searchLayer: <paste>
Here's my Insert: <paste>
Here's my selectNeighborsHeuristic: <paste>
Here's my test setup: <paste>

Walk through the 6-node example by hand to verify searchLayer correctness, then identify
the most likely cause based on my measured recall: <X%>.
Suggest ONE targeted change at a time.
```

---

## 🏁 Wrap-up (20 min)

**▶️ RUN:**
```powershell
git add .
git commit -m "feat(day-4): Search, BruteForceKNN, RecallAtK; recall@10 = X% on 10K (C1)"
```
(把 X% 改成实测)

**Today's checklist:**
- [ ] **CHECKPOINT 1: Recall@10 ≥ 90%** — ✅ / ❌

❌ 的话 Day 5 全天 debug,不要进 SIFT-1M。

---
---

# Day 5 — 2026-05-15 (周五 / Fri) — SIFT-1M 加载 + 建索引

> 如果 Day 4 C1 没过 → 今天全天 debug。

## Block 1 (90 min) — 下载 + Loader

### 5.1 下载 SIFT-1M (~15 min)

**▶️ RUN (PowerShell):**
```powershell
cd F:\app\job\velosearch\benchmark\sift
mkdir data -Force
cd data
# Option A: download sift.tar.gz (~170 MB) in a browser from http://corpus-texmex.irisa.fr/
# Option B: command line
Invoke-WebRequest -Uri "ftp://ftp.irisa.fr/local/texmex/corpus/sift.tar.gz" -OutFile sift.tar.gz
tar -xzf sift.tar.gz
# After extraction:
#   sift_base.fvecs        ~516 MB (1M × 128 float)
#   sift_query.fvecs       ~5.1 MB (10K × 128 float)
#   sift_groundtruth.ivecs ~4.1 MB (10K × 100 int)
#   sift_learn.fvecs       not used
cd F:\app\job\velosearch
```

(`.gitignore` 已包含 `benchmark/sift/data/`,这些文件不会进 git。)

### 5.2 fvecs/ivecs reader (~60 min)

**目标文件:** `benchmark/sift/loader.go`, `loader_test.go`

**🤖 LLM PROMPT:**
```
Write a Go package "sift" at benchmark/sift/loader.go to read the SIFT-1M binary format
from http://corpus-texmex.irisa.fr/.

Format spec:
  Each vector preceded by 4-byte little-endian int32 = dimension.
  Then dimension × 4 bytes = float32 (.fvecs) or int32 (.ivecs).
  Repeated until EOF.

Provide:
  func LoadFvecs(path string) ([][]float32, error)
  func LoadIvecs(path string) ([][]int32, error)

Performance: use bufio.Reader with 1 MB buffer, batch-read each vector.

Also write benchmark/sift/loader_test.go that:
  - asserts LoadFvecs("data/sift_base.fvecs") returns len=1_000_000, len(out[0])=128
  - asserts LoadIvecs("data/sift_groundtruth.ivecs") returns len=10_000, len(out[0])=100
  - skips with t.Skip() if file not found (so CI passes without dataset)
```

**🎯 VERIFY:**
```powershell
go test ./benchmark/sift/ -v
```

### 5.3 内存预估 (~15 min)

**📚 REF:**
```
1M vectors × 128 floats × 4 bytes = 512 MB (raw)
+ HNSW graph: avg ~M edges/node × 2 directions × 4 bytes × 1M ≈ 128 MB
+ Go map overhead: ~50-100 MB
≈ 700 MB - 1 GB total in RAM
```

打开任务管理器确认空闲 RAM > 2 GB。

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — 建 1M 索引

### 5.4 Runner (30 min)

**📋 COPY → `benchmark/sift/cmd_build/main.go`** (新建子目录,这是 main 包):
```go
// This file lives at benchmark/sift/cmd_build/main.go
// because loader.go declares package sift; main can't share that package.
// We use a sub-directory cmd_build/ to host the main package.
package main

import (
	"encoding/gob"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/zhangchuqi1998/velosearch/benchmark/sift"
	"github.com/zhangchuqi1998/velosearch/internal/distance"
	"github.com/zhangchuqi1998/velosearch/internal/hnsw"
)

func main() {
	log.Println("Loading sift_base.fvecs...")
	t0 := time.Now()
	base, err := sift.LoadFvecs("benchmark/sift/data/sift_base.fvecs")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d vectors in %v", len(base), time.Since(t0))

	idx := hnsw.NewIndex(128, 16, 200, distance.L2Squared)
	log.Println("Building index...")
	t1 := time.Now()
	for i, v := range base {
		idx.Insert(uint32(i), v)
		if (i+1)%100_000 == 0 {
			elapsed := time.Since(t1)
			rate := float64(i+1) / elapsed.Seconds()
			log.Printf("  inserted %d/%d  (%.0f/sec, elapsed %v)", i+1, len(base), rate, elapsed)
		}
	}
	buildTime := time.Since(t1)
	log.Printf("Build complete in %v", buildTime)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("HeapAlloc: %.2f GB", float64(m.HeapAlloc)/1e9)
	log.Printf("Insert throughput (cold): %.0f/sec", float64(len(base))/buildTime.Seconds())

	// Simple gob save -- replaced by WAL on Day 10
	f, err := os.Create("benchmark/sift/data/sift_index.gob")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := gob.NewEncoder(f).Encode(idx); err != nil {
		log.Printf("gob encode failed (will rebuild Day 6): %v", err)
	}
}
```

**注:** `idx` 包含 `sync.RWMutex` 和 `*rand.Rand`,gob 编码可能失败 — 没关系,Day 6 直接重建一次就行,或者你也可以加 `MarshalBinary` 方法。这是非阻塞,可以略。

### 5.5 跑 + 看输出 (60 min)

**▶️ RUN:**
```powershell
cd F:\app\job\velosearch
go run ./benchmark/sift/cmd_build
# Build takes 10-30 minutes; memory stays < 1.5 GB
```

**期间不要干瞪眼** — 把 Day 6 的 bench 脚手架先建起来,或者读读 Day 7 的 SIMD 部分。

---

## 🏁 Wrap-up (20 min)

**▶️ RUN:**
```powershell
git add benchmark/sift/loader.go benchmark/sift/loader_test.go benchmark/sift/cmd_build/main.go
git commit -m "feat(day-5): SIFT-1M loader, 1M build runner; build time = X min"
```

**WEEKLY_LOG:** 把 build time / HeapAlloc / cold insert rate 都写下来。

---
---

# Day 6 — 2026-05-16 (周六 / Sat) — Benchmark + 参数 Sweep

## 今日目标
- 10K 查询测 recall@10 + P50/P95/P99
- (M, efConstruction, efSearch) 参数 sweep
- 找到最佳 efSearch 使 recall ≥ 95%

## Block 1 (90 min) — Baseline

### 6.1 Bench Runner (40 min)

**目标文件:** `benchmark/sift/cmd_bench/main.go`

**🤖 LLM PROMPT:**
```
Write a benchmark runner in Go at benchmark/sift/cmd_bench/main.go.

Loads from benchmark/sift/data/:
  sift_query.fvecs        (10000 × 128 float32)
  sift_groundtruth.ivecs  (10000 × 100 int32)

Takes flags:
  -ef int (default 50)         efSearch parameter
  -k int (default 10)          number of neighbors to retrieve

Either:
  a) loads pre-built index from sift_index.gob (if it loads), OR
  b) rebuilds in-memory by calling sift.LoadFvecs("sift_base.fvecs") and idx.Insert in a loop.

For each of 10000 queries:
  - times idx.Search(query, k, ef) with time.Now()
  - computes recall@k against the first k ground truth IDs

Output:
  Mean recall: X.XXXX
  Latency P50: X us, P95: X us, P99: X us
  Total wall time: Xs

Use slices.Sort for percentiles. import: "github.com/zhangchuqi1998/velosearch/internal/hnsw"
and "github.com/zhangchuqi1998/velosearch/internal/distance".
```

### 6.2 Baseline 跑分 (20 min)

**▶️ RUN:**
```powershell
go run ./benchmark/sift/cmd_bench -ef 50
# Expected (unoptimized):
#   Recall ≈ 0.90-0.95
#   P99 ≈ 50-200 ms
```

**✍️ WRITE: 记到 NOTES.md:**
```markdown
## SIFT-1M Baseline (Day 6)
Config: M=16, efC=200, efSearch=50
- Recall@10: 0.XX
- P50: X µs, P95: X µs, P99: X µs
```

### 6.3 Insert throughput micro-bench (30 min)

**📋 COPY → `benchmark/sift/insert_bench_test.go`**:
```go
package sift

import (
	"testing"

	"github.com/zhangchuqi1998/velosearch/internal/distance"
	"github.com/zhangchuqi1998/velosearch/internal/hnsw"
)

func BenchmarkInsertSIFT(b *testing.B) {
	base, err := LoadFvecs("data/sift_base.fvecs")
	if err != nil {
		b.Skip("sift_base.fvecs not found")
	}
	if len(base) < 200_000 {
		b.Fatal("need at least 200k vectors for benchmark")
	}
	idx := hnsw.NewIndex(128, 16, 200, distance.L2Squared)
	// Warm up with 100K
	for i := 0; i < 100_000; i++ {
		idx.Insert(uint32(i), base[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.Insert(uint32(100_000+i), base[100_000+(i%100_000)])
	}
}
```

**▶️ RUN:**
```powershell
go test -bench=BenchmarkInsertSIFT -benchtime=10000x -run=^$ ./benchmark/sift
# Outputs X ns/op; throughput = 1e9 / X ops/sec
```

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — Parameter Sweep

### 6.4 efSearch sweep (30 min)

**▶️ RUN (PowerShell):**
```powershell
foreach ($ef in 30, 50, 75, 100, 150, 200, 300, 400) {
    Write-Host "=== efSearch=$ef ==="
    go run ./benchmark/sift/cmd_bench -ef $ef
}
```

**✍️ WRITE: 把数据填到 NOTES.md 这个表:**

| efSearch | Recall@10 | P50 (µs) | P99 (µs) |
|----------|-----------|----------|----------|
| 30 | | | |
| 50 | | | |
| 75 | | | |
| 100 | | | |
| 150 | | | |
| 200 | | | |
| 300 | | | |
| 400 | | | |

挑出"最小 efSearch 使 recall ≥ 95%"。

### 6.5 M / efC sweep (40 min)

只有当 efSearch sweep 后 recall 卡在 < 95% 才需要做这步。

**▶️ RUN:** 改 `cmd_build/main.go` 里的 `hnsw.NewIndex(128, M, efC, ...)`,重建索引,重测。 每次 ~30 min。

### 6.6 Pareto 表 (20 min)

**✍️ WRITE: NOTES.md Pareto 表:**

| Config (M, efC, efS) | Recall | P99 (µs) | Build (min) |
|----------------------|--------|----------|-------------|
| 16, 200, 50 | | | |
| 16, 200, 100 | | | |
| 24, 400, 100 | | | |

**关键决策:** 选**最低 (M, efC) 使 recall ≥ 95%** — 加 M / efC 只多耗内存和建图时间,简历看不出区别。

---

## 🏁 Wrap-up (20 min)

```powershell
git add .
git commit -m "feat(day-6): SIFT bench runner, param sweep; chosen M=X efC=X ef=X recall=X.XX P99=X ms"
```

---
---

# Day 7 — 2026-05-17 (周日 / Sun) — 🎯 CHECKPOINT 2

## 今日目标
- pprof 找热点
- SIMD 距离 + 连续内存布局
- **P99 < 15ms @ recall ≥ 95% on SIFT-1M**

## Block 1 (90 min) — pprof + SIMD

### 7.1 CPU profile (15 min)

**▶️ RUN:**
```powershell
go test -bench=BenchmarkInsertSIFT -cpuprofile=cpu.prof -benchtime=10000x -run=^$ ./benchmark/sift
go tool pprof -http=:8080 cpu.prof
```

打开 `http://localhost:8080`,Flame Graph。**预期 70-90% 时间在 L2Squared。**

### 7.2 SIMD distance via vek (45 min)

**▶️ RUN:**
```powershell
go get github.com/viterin/vek
```

**📋 COPY → `internal/distance/distance.go`** (替换原 L2Squared 实现):
```go
package distance

import "github.com/viterin/vek/vek32"

type DistanceFunc func(a, b []float32) float32

// L2Squared is the squared L2 distance, AVX2-accelerated via vek32.
func L2Squared(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("length mismatch")
	}
	return vek32.DistanceSquared(a, b)
}

// Cosine distance: 1 - dot/(|a|*|b|)
func Cosine(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("length mismatch")
	}
	dot := vek32.Dot(a, b)
	na := vek32.Norm(a)
	nb := vek32.Norm(b)
	if na == 0 || nb == 0 {
		return 1.0
	}
	return 1.0 - dot/(na*nb)
}
```

**🎯 VERIFY:**
```powershell
go test ./internal/distance/ -v
go test -bench=BenchmarkL2Squared -benchmem ./internal/distance/
# Expect ns/op to drop from ~150 ns to ~30-40 ns
```

### 7.3 BCE hint (30 min) — Optional

**📚 REF — 不用 vek 时的手写版本:**
```go
func L2SquaredManual(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("length mismatch")
	}
	_ = b[len(a)-1] // BCE hint: tell compiler b is at least len(a) long
	var sum float32
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return sum
}
```

`go build -gcflags="-d=ssa/check_bce/debug=1" ./internal/distance/` 可以看哪些 bounds-check 没被消除。

**用了 vek 就不用这一步。**

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — Memory Layout + Re-measure

### 7.4 连续内存布局 (60 min)

**📚 REF — 当前问题:**
每个 Node 自带 `Vector []float32` — 1M 个独立 slice 头,GC 跟踪 1M 对象,缓存不友好。

**🤖 LLM PROMPT:**
```
Refactor my HNSW Index in Go to store all vectors in one contiguous []float32 instead of per-Node slices.

Current code:
<paste your internal/hnsw/index.go, internal/hnsw/insert.go, internal/hnsw/search.go>

Refactor plan:
1. Add to Index struct:
     vectorStore []float32         // length = N * Dim
     nextOffset  uint32            // running counter for next free slot in N space
2. Remove Vector []float32 from Node struct.
3. Add method:
     func (idx *Index) vectorAt(id uint32) []float32 {
         start := int(id) * idx.Dim
         return idx.vectorStore[start : start+idx.Dim]
     }
4. In Insert: append vector to vectorStore.
5. Replace all `idx.nodes[id].Vector` and `node.Vector` reads with idx.vectorAt(id).

Assumption (for v0.1): IDs are contiguous starting from 0. Generalization later.

Show me ONLY the diff for:
- Index struct (in index.go)
- NewIndex (need to pre-allocate vectorStore based on a capacity hint — accept a new param `capacity int`)
- Node struct
- vectorAt method (new)
- Insert function (changes only)
- searchLayer function (changes only)
- Search function (changes only)
- pruneNeighbors function (changes only)
- All existing tests (update calls to NewIndex with capacity argument)
```

**注:** `NewIndex` 签名要从 `NewIndex(dim, M, efC, dist)` 改成 `NewIndex(dim, M, efC, capacity, dist)`。所有测试要跟着改。

### 7.5 Re-measure + 决策 (30 min)

**▶️ RUN:**
```powershell
go test ./internal/hnsw/ -v   # All previous tests must still pass
go run ./benchmark/sift/cmd_bench -ef <chosen efSearch>
```

**✍️ WRITE: NOTES.md 优化进展表:**

| Stage | Recall | P99 (µs) |
|-------|--------|----------|
| Baseline (Day 6) | 0.97 | 14000 |
| + SIMD L2 (vek32) | | |
| + Contiguous store | | |

**📚 REF — 目标决策:**
- ✅ 到 `P99 < 15000 µs at recall ≥ 0.95` → CHECKPOINT 2 通过
- ❌ 没到 → **简历改成实测的 ms 数,不要硬撑** — 诚实 22ms 胜过编 15ms

---

## 🏁 Wrap-up (20 min)

```powershell
git add .
git commit -m "perf(day-7): SIMD L2 + contiguous vectorStore; P99 = X ms @ recall X.XX (C2)"
```

**WEEKLY_LOG:** 写一段 Week 1 总结。

---
---

# Day 8 — 2026-05-18 (周一 / Mon) — Protobuf + Collection Manager

## 今日目标
- 设计 gRPC API
- `make proto` 生成代码
- Collection Manager (多租户路由)

## Block 1 (90 min) — Proto schema

### 8.1 Proto schema (45 min)

**📋 COPY → `proto/velosearch.proto`**:
```protobuf
syntax = "proto3";
package velosearch.v1;
option go_package = "github.com/zhangchuqi1998/velosearch/proto;velosearchv1";

service VectorSearch {
  rpc CreateCollection(CreateCollectionRequest) returns (CreateCollectionResponse);
  rpc DropCollection(DropCollectionRequest) returns (DropCollectionResponse);
  rpc ListCollections(ListCollectionsRequest) returns (ListCollectionsResponse);

  rpc Insert(InsertRequest) returns (InsertResponse);
  rpc Delete(DeleteRequest) returns (DeleteResponse);
  rpc Search(SearchRequest) returns (SearchResponse);

  rpc Stats(StatsRequest) returns (StatsResponse);
}

enum Metric {
  METRIC_UNSPECIFIED = 0;
  METRIC_L2 = 1;
  METRIC_COSINE = 2;
}

message Vector {
  repeated float values = 1;
}

message CreateCollectionRequest {
  string name = 1;
  int32 dim = 2;
  Metric metric = 3;
  int32 m = 4;
  int32 ef_construction = 5;
}
message CreateCollectionResponse { bool created = 1; }

message DropCollectionRequest  { string name = 1; }
message DropCollectionResponse { bool dropped = 1; }

message ListCollectionsRequest {}
message ListCollectionsResponse { repeated string names = 1; }

message Item {
  uint32 id = 1;
  Vector vector = 2;
}
message InsertRequest  { string collection = 1; repeated Item items = 2; }
message InsertResponse { int32 inserted = 1; }

message DeleteRequest  { string collection = 1; repeated uint32 ids = 2; }
message DeleteResponse { int32 deleted = 1; }

message SearchRequest {
  string collection = 1;
  Vector query = 2;
  int32 k = 3;
  int32 ef_search = 4;
}
message Hit { uint32 id = 1; float distance = 2; }
message SearchResponse { repeated Hit hits = 1; }

message StatsRequest  { string collection = 1; }
message StatsResponse {
  int32 num_vectors = 1;
  int32 num_deleted = 2;
  int32 num_layers = 3;
  int64 mem_bytes = 4;
}
```

### 8.2 生成 Go 代码 (10 min)

**▶️ RUN:**
```powershell
make proto
# Generates proto/velosearch.pb.go and proto/velosearch_grpc.pb.go
```

### 8.3 Commit 生成代码 (5 min)

**▶️ RUN:**
```powershell
git add proto/
git commit -m "feat(day-8): proto schema for VectorSearch gRPC API"
```

(惯例: 生成代码入库,clone 后不需要装 protoc 也能 build。)

### 8.4 Config types (30 min)

**📋 COPY → `internal/collection/config.go`**:
```go
package collection

import "github.com/zhangchuqi1998/velosearch/internal/distance"

type Metric int

const (
	MetricL2 Metric = iota
	MetricCosine
)

type Config struct {
	Name           string
	Dim            int
	Metric         Metric
	M              int
	EfConstruction int
	Capacity       int // estimated capacity, used to preallocate the contiguous vectorStore
}

func (m Metric) DistanceFunc() distance.DistanceFunc {
	switch m {
	case MetricL2:
		return distance.L2Squared
	case MetricCosine:
		return distance.Cosine
	}
	return distance.L2Squared
}
```

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — Manager

### 8.5 Manager 实现 (60 min)

**📋 COPY → `internal/collection/manager.go`**:
```go
package collection

import (
	"errors"
	"sync"

	"github.com/zhangchuqi1998/velosearch/internal/hnsw"
)

var (
	ErrAlreadyExists = errors.New("collection already exists")
	ErrNotFound      = errors.New("collection not found")
	ErrDimMismatch   = errors.New("vector dimension does not match collection")
)

type Collection struct {
	Config Config
	Index  *hnsw.Index
}

type Manager struct {
	mu   sync.RWMutex
	cols map[string]*Collection
}

func NewManager() *Manager {
	return &Manager{cols: make(map[string]*Collection)}
}

func (m *Manager) Create(cfg Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.cols[cfg.Name]; ok {
		return ErrAlreadyExists
	}
	idx := hnsw.NewIndex(cfg.Dim, cfg.M, cfg.EfConstruction, cfg.Capacity, cfg.Metric.DistanceFunc())
	m.cols[cfg.Name] = &Collection{Config: cfg, Index: idx}
	return nil
}

func (m *Manager) Get(name string) (*Collection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.cols[name]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (m *Manager) Drop(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.cols[name]; !ok {
		return ErrNotFound
	}
	delete(m.cols, name)
	return nil
}

func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.cols))
	for n := range m.cols {
		names = append(names, n)
	}
	return names
}
```

### 8.6 测试 (30 min)

**🤖 LLM PROMPT:**
```
Write internal/collection/manager_test.go with these tests:
1. TestCreate_Success: Create with valid cfg, Get returns same Collection instance
2. TestCreate_Duplicate: Create twice with same name → second returns ErrAlreadyExists
3. TestGet_NotFound: Get on non-existent name → ErrNotFound
4. TestDrop: Drop existing → success; subsequent Get → ErrNotFound
5. TestDrop_NotFound: Drop non-existent → ErrNotFound
6. TestList: empty manager → []; after creating 3 → list contains all 3 (order-agnostic)
7. TestConcurrentCreate: 100 goroutines Create different names; check `go test -race` passes

Use github.com/zhangchuqi1998/velosearch/internal/collection package.
Use t.Run for subtests where appropriate.
```

**🎯 VERIFY:**
```powershell
go test -race ./internal/collection/ -v
```

---

## 🏁 Wrap-up (20 min)

```powershell
git add .
git commit -m "feat(day-8): Collection Manager + tests (race-clean)"
```

---
---

# Day 9 — 2026-05-19 (周二 / Tue) — gRPC Handlers + Server + Integration Test

## 今日目标
- 实现所有 gRPC handler
- main.go 启动服务
- bufconn 端到端测试

## Block 1 (90 min) — Handlers

### 9.1 Server scaffold (20 min)

**📋 COPY → `internal/grpcserver/server.go`** (起始,handler 你接下来加):
```go
package grpcserver

import (
	pb "github.com/zhangchuqi1998/velosearch/proto"
	"github.com/zhangchuqi1998/velosearch/internal/collection"
)

type Server struct {
	pb.UnimplementedVectorSearchServer
	mgr *collection.Manager
}

func New(mgr *collection.Manager) *Server {
	return &Server{mgr: mgr}
}
```

### 9.2 Handlers (60 min)

**🤖 LLM PROMPT:**
```
Implement all gRPC handlers as methods on *Server (in internal/grpcserver/server.go).

I have:
- type Server struct { pb.UnimplementedVectorSearchServer; mgr *collection.Manager }
- Manager methods: Create(cfg), Get(name), Drop(name), List()
- Each Collection has Index *hnsw.Index with:
    Insert(id uint32, v []float32)
    Search(q []float32, k, efSearch int) []hnsw.Candidate
    Delete(id uint32) error
    Dim int

Proto: <paste proto/velosearch.proto>

For each RPC handler:
1. Validate inputs:
   - CreateCollection: dim > 0, m >= 4, ef_construction >= m
   - Insert: every Item.vector.values length must equal collection's Dim
   - Search: k > 0, ef_search >= k, query.values length == Dim
2. Map errors to gRPC status codes (google.golang.org/grpc/codes + status):
   - collection.ErrAlreadyExists → codes.AlreadyExists
   - collection.ErrNotFound → codes.NotFound
   - validation → codes.InvalidArgument
   - other → codes.Internal
3. Stats: count vectors / deleted / layers / approximate mem bytes by walking idx.nodes
   (you may add a Stats() method on *hnsw.Index that returns these — show me that too)
4. Log every RPC at info level via log/slog with collection name and request size.

Output complete server.go.
```

### 9.3 Index.Delete + Index.Stats (10 min)

**📋 COPY → `internal/hnsw/index.go`** (追加):
```go
import "errors"

func (idx *Index) Delete(id uint32) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	n, ok := idx.nodes[id]
	if !ok {
		return errors.New("id not found")
	}
	n.Deleted = true
	return nil
}

type Stats struct {
	NumVectors int
	NumDeleted int
	NumLayers  int
	MemBytes   int64
}

func (idx *Index) Stats() Stats {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	s := Stats{NumVectors: len(idx.nodes), NumLayers: idx.maxLevel + 1}
	for _, n := range idx.nodes {
		if n.Deleted {
			s.NumDeleted++
		}
	}
	// rough memory estimate
	s.MemBytes = int64(len(idx.nodes)) * int64(idx.Dim) * 4
	return s
}
```

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — main.go + Integration Test

### 9.4 main.go (30 min)

**📋 COPY → `cmd/server/main.go`** (替换之前 Day 1 的 stub):
```go
package main

import (
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/zhangchuqi1998/velosearch/internal/collection"
	"github.com/zhangchuqi1998/velosearch/internal/grpcserver"
	pb "github.com/zhangchuqi1998/velosearch/proto"
	"google.golang.org/grpc"
)

func main() {
	addr := flag.String("addr", ":50051", "gRPC listen address")
	dataDir := flag.String("data-dir", "./data", "WAL data directory (used Day 10)")
	flag.Parse()
	_ = dataDir

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mgr := collection.NewManager()
	srv := grpcserver.New(mgr)

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		slog.Error("listen failed", "err", err)
		os.Exit(1)
	}

	grpcSrv := grpc.NewServer()
	pb.RegisterVectorSearchServer(grpcSrv, srv)

	go func() {
		slog.Info("server listening", "addr", *addr)
		if err := grpcSrv.Serve(lis); err != nil {
			slog.Error("serve failed", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	slog.Info("shutting down")
	grpcSrv.GracefulStop()
}
```

**🎯 VERIFY:**
```powershell
go run ./cmd/server
# In another terminal:
grpcurl -plaintext localhost:50051 list
# Expected: velosearch.v1.VectorSearch
```

### 9.5 Integration test (60 min)

**目标文件:** `internal/grpcserver/integration_test.go`

**🤖 LLM PROMPT:**
```
Write an integration test for my gRPC server using google.golang.org/grpc/test/bufconn
(in-memory listener, no TCP).

Test flow:
1. Setup: bufconn.Listen(1024*1024), start *grpc.Server with VectorSearch handler
   in a goroutine, return a *grpc.ClientConn over the bufconn.
2. CreateCollection("test", dim=8, metric=METRIC_L2, m=16, ef_construction=200)
3. Insert 100 random 8-d vectors with IDs 0..99 (use math/rand seeded with 42)
4. Search using vector at index 42 as query, k=5, ef_search=50
   Assert: SearchResponse.Hits[0].Id == 42 (because identical vectors give 0 distance)
5. Delete id=42
6. Search same query again, assert top hit ID != 42
7. Stats — assert NumVectors == 100, NumDeleted == 1

Use t.Cleanup() to gracefully stop the server.
Imports needed: pb "...velosearch/proto", "...velosearch/internal/grpcserver",
"...velosearch/internal/collection", "google.golang.org/grpc",
"google.golang.org/grpc/credentials/insecure", "google.golang.org/grpc/test/bufconn".

Output the complete integration_test.go.
```

**🎯 VERIFY:**
```powershell
go test -race ./internal/grpcserver/ -v
```

---

## 🏁 Wrap-up (20 min)

```powershell
git add .
git commit -m "feat(day-9): gRPC handlers, main.go, bufconn integration test"
```

---
---

# Day 10 — 2026-05-20 (周三 / Wed) — WAL Persistence

## 今日目标
- WAL record 格式 (length-prefixed protobuf + CRC32)
- WAL writer + replay
- Handlers 全部写 WAL

## Block 1 (90 min) — WAL Format + Writer

### 10.1 WAL proto (20 min)

**📋 COPY → `proto/wal.proto`**:
```protobuf
syntax = "proto3";
package velosearch.wal.v1;
option go_package = "github.com/zhangchuqi1998/velosearch/proto;velosearchwalv1";

message WALRecord {
  oneof op {
    CreateColl create_coll = 1;
    DropColl   drop_coll   = 2;
    InsertOp   insert      = 3;
    DeleteOp   delete      = 4;
  }
}

message CreateColl {
  string name = 1;
  int32  dim = 2;
  int32  metric = 3;   // 0=L2, 1=Cosine
  int32  m = 4;
  int32  ef_construction = 5;
  int32  capacity = 6;
}
message DropColl { string name = 1; }
message InsertOp { string collection = 1; uint32 id = 2; repeated float vector = 3; }
message DeleteOp { string collection = 1; uint32 id = 2; }
```

**▶️ RUN:**
```powershell
# Update Makefile's proto target to include wal.proto, or run manually:
protoc --go_out=. --go_opt=paths=source_relative proto/wal.proto
```

### 10.2 WAL writer + replay (70 min)

**目标文件:** `internal/storage/wal.go`, `wal_test.go`

**📚 REF — On-disk format:**
```
+----------+----------+----------+
| 4 bytes  | N bytes  | 4 bytes  |
| BE len N | payload  | CRC32    |
+----------+----------+----------+
... repeated ...
```

**🤖 LLM PROMPT:**
```
Write a write-ahead log for a vector database in Go.

File format per record:
  - 4-byte big-endian uint32 = payload length N
  - N bytes = protobuf-encoded WALRecord (from proto/wal.proto)
  - 4-byte big-endian uint32 = CRC32 (Castagnoli polynomial) of payload

Place in internal/storage/wal.go:

type WAL struct {
    f    *os.File
    mu   sync.Mutex
    size int64
}

func Open(path string) (*WAL, error)   // os.O_APPEND|O_CREATE|O_WRONLY, mode 0644
func (w *WAL) Append(rec *walpb.WALRecord) error   // marshal, write len+payload+crc, call f.Sync()
func (w *WAL) Close() error
func Replay(path string, handler func(*walpb.WALRecord) error) error
  // Reads records sequentially, verifies CRC, calls handler for each.
  // On partial-write tail (truncated record from crash), LOG WARN and stop — don't fail.
  // On CRC mismatch mid-file, return error.

Imports:
  walpb "github.com/zhangchuqi1998/velosearch/proto"
  "google.golang.org/protobuf/proto"
  "hash/crc32"
  "encoding/binary"
  "io"

Use crc32.MakeTable(crc32.Castagnoli) as a package-level var.

Also write wal_test.go with:
1. TestWAL_AppendReplay: open new WAL, append 5 records, close, replay → 5 records read in order
2. TestWAL_TruncatedTail: append 5 records, manually truncate file to mid-record, replay → warn + 4 records (or however many were complete)
3. TestWAL_CorruptedBody: flip 1 byte in payload of middle record, replay → error at that record

Use t.TempDir() for the WAL path.
```

**🎯 VERIFY:**
```powershell
go test ./internal/storage/ -v
```

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — Integration

### 10.3 main.go 加 replay (30 min)

**🧩 SKELETON → `cmd/server/main.go`** (在 `mgr := collection.NewManager()` 之后,`grpcSrv.Serve` 之前插入):
```go
import (
	"path/filepath"
	"github.com/zhangchuqi1998/velosearch/internal/storage"
	walpb "github.com/zhangchuqi1998/velosearch/proto"
)

// ... inside main():

if err := os.MkdirAll(*dataDir, 0755); err != nil {
	slog.Error("mkdir failed", "err", err)
	os.Exit(1)
}
walPath := filepath.Join(*dataDir, "wal.log")

w, err := storage.Open(walPath)
if err != nil {
	slog.Error("open wal failed", "err", err)
	os.Exit(1)
}
defer w.Close()

// TODO: implement ApplyWALRecord (see 10.4) and call it here
slog.Info("replaying WAL...")
nRec := 0
if err := storage.Replay(walPath, func(rec *walpb.WALRecord) error {
	nRec++
	return ApplyWALRecord(mgr, rec) // see 10.4
}); err != nil {
	slog.Error("replay failed", "err", err)
	os.Exit(1)
}
slog.Info("replay done", "records", nRec)

srv := grpcserver.New(mgr, w) // Server now takes a WAL handle (see 10.4)
```

### 10.4 ApplyWALRecord + Server 改造 (40 min)

**📋 COPY → `cmd/server/apply.go`** (新文件):
```go
package main

import (
	"github.com/zhangchuqi1998/velosearch/internal/collection"
	walpb "github.com/zhangchuqi1998/velosearch/proto"
)

// ApplyWALRecord applies a single WAL record to the Manager.
// Used by both startup replay and runtime apply.
func ApplyWALRecord(mgr *collection.Manager, rec *walpb.WALRecord) error {
	switch op := rec.Op.(type) {
	case *walpb.WALRecord_CreateColl:
		return mgr.Create(collection.Config{
			Name:           op.CreateColl.Name,
			Dim:            int(op.CreateColl.Dim),
			Metric:         collection.Metric(op.CreateColl.Metric),
			M:              int(op.CreateColl.M),
			EfConstruction: int(op.CreateColl.EfConstruction),
			Capacity:       int(op.CreateColl.Capacity),
		})
	case *walpb.WALRecord_DropColl:
		return mgr.Drop(op.DropColl.Name)
	case *walpb.WALRecord_Insert:
		c, err := mgr.Get(op.Insert.Collection)
		if err != nil {
			return err
		}
		c.Index.Insert(op.Insert.Id, op.Insert.Vector)
		return nil
	case *walpb.WALRecord_Delete:
		c, err := mgr.Get(op.Delete.Collection)
		if err != nil {
			return err
		}
		return c.Index.Delete(op.Delete.Id)
	}
	return nil
}
```

**🧩 SKELETON: 改造 `internal/grpcserver/server.go`**

把 `Server` 加上 `wal *storage.WAL` 字段,所有 mutating handler 改成"先 WAL,后 memory"。

**📚 REF — 改造模板 (以 Insert 为例):**
```go
func (s *Server) Insert(ctx context.Context, req *pb.InsertRequest) (*pb.InsertResponse, error) {
	// 1. validate (same as before)
	// ...

	// 2. WAL first
	for _, item := range req.Items {
		rec := &walpb.WALRecord{Op: &walpb.WALRecord_Insert{
			Insert: &walpb.InsertOp{
				Collection: req.Collection,
				Id:         item.Id,
				Vector:     item.Vector.Values,
			},
		}}
		if err := s.wal.Append(rec); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	// 3. memory update
	c, err := s.mgr.Get(req.Collection)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	for _, item := range req.Items {
		c.Index.Insert(item.Id, item.Vector.Values)
	}
	return &pb.InsertResponse{Inserted: int32(len(req.Items))}, nil
}
```

**⚠️ 顺序关键:** WAL 先,memory 后。倒过来 = crash 时 memory 写了 WAL 没记录,数据丢失。

**✍️ WRITE: 把 Insert/Delete/CreateCollection/DropCollection 四个 handler 都改成这个模式。**

### 10.5 改 New 函数 (10 min)

**🧩 SKELETON → `internal/grpcserver/server.go`**:
```go
type Server struct {
	pb.UnimplementedVectorSearchServer
	mgr *collection.Manager
	wal *storage.WAL
}

func New(mgr *collection.Manager, wal *storage.WAL) *Server {
	return &Server{mgr: mgr, wal: wal}
}
```

集成测试也要同步改 — 测试里要传一个临时 WAL(用 `t.TempDir()`)。

---

## 🏁 Wrap-up (20 min)

```powershell
git add .
git commit -m "feat(day-10): WAL with CRC32, replay-on-startup, handlers write WAL first"
```

**手动冒烟测试:**
```powershell
go run ./cmd/server -data-dir=./testdata/smoke
# In another terminal insert a few records, then Ctrl+C to stop
# Restart:
go run ./cmd/server -data-dir=./testdata/smoke
# Expect 'replay done' followed by a non-zero record count
```

---
---

# Day 11 — 2026-05-21 (周四 / Thu) — 🎯 CHECKPOINT 3 / Crash Recovery

## 今日目标
- Tombstone delete 在 search 里正确过滤
- `kill -9` → 重启 → 已确认 insert 全部可搜
- 10 次 crash test 不挂

## Block 1 (90 min) — Tombstone Delete

### 11.1 Search 已经过滤 deleted(Day 4 已加) (5 min)

确认你 Day 4 写的 `Search` 函数里这段还在:
```go
for _, c := range candidates {
    if idx.nodes[c.ID].Deleted { continue }
    out = append(out, c)
    if len(out) == k { break }
}
```

**⚠️ 注意:** `searchLayer` 里**继续遍历** deleted 节点的邻居(它们可能是通往非 deleted 节点的桥)。只有最终结果过滤 deleted。

### 11.2 Delete-search 测试 (50 min)

**📋 COPY → `internal/hnsw/delete_test.go`**:
```go
package hnsw

import (
	"testing"

	"github.com/zhangchuqi1998/velosearch/internal/distance"
)

func TestDelete_FiltersFromResults(t *testing.T) {
	idx := NewIndex(128, 16, 200, 1000, distance.L2Squared)
	// ... insert 1000 vectors, delete every 10th, run 100 queries, assert no deleted IDs in results
	// (adapt from TestRecall10KRandom)
}
```

**🤖 LLM PROMPT:**
```
Write internal/hnsw/delete_test.go with TestDelete_FiltersFromResults:

1. Build idx with 1000 random 128-d vectors (rand.NewSource(42), capacity=1000)
2. Delete every 10th ID (100, 110, 120, ..., 990)
3. Run 100 random queries with k=10, efSearch=50
4. Assert: NO result has an ID that was deleted

Compare against TestRecall10KRandom (in internal/hnsw/recall_test.go) for the test pattern.
```

### 11.3 Stats 验证 (35 min)

加一个简单测试:Delete 之后 Stats.NumDeleted 应该等于删的数量。

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — Crash Recovery Test

### 11.4 Crash client (40 min)

**🧩 SKELETON → `benchmark/crash_client/main.go`**:
```go
// Usage:
//   go run ./benchmark/crash_client -addr=localhost:50052 -mode=write -n=1000
//   go run ./benchmark/crash_client -addr=localhost:50052 -mode=verify -n=1000
//
// write mode: CreateCollection(name=crash, dim=32) + Insert n deterministic vectors with ID = 0..n-1
// verify mode: for each ID i, search using vector identical to inserted i with k=1,
//             assert returned Hit.ID == i; any miss exits with code 1
package main

func main() {
	// parse -addr -mode -n with flag.Parse()
	// dial the server with grpc.Dial
	// write branch: CreateCollection + Insert in batches of 100; await each batch
	// verify branch: loop Search calls, track hits vs misses, non-zero exit on any miss
	panic("TODO: implement")
}
```

**🤖 LLM PROMPT:**
```
Implement benchmark/crash_client/main.go (Go).

Flags:
  -addr string   gRPC server address (default "localhost:50052")
  -mode string   "write" or "verify"
  -n int         number of vectors

In WRITE mode:
1. CreateCollection(name="crash", dim=32, METRIC_L2, m=16, ef_construction=200, capacity=n*2)
2. Generate n deterministic vectors: vector[i] = float32 array of length 32 where
   vector[i][j] = float32(i + j) (any deterministic function, just must be reproducible in verify mode)
3. Insert in batches of 100, await server response each batch
4. Print "DONE writing N vectors"

In VERIFY mode:
1. For i in 0..n-1:
   - Build the same deterministic vector for id i
   - Search(k=1)
   - If Hits[0].Id != i: increment miss counter
2. Print "verified n=N misses=M"
3. Exit 1 if any miss, 0 otherwise

Use proto package github.com/zhangchuqi1998/velosearch/proto.
```

### 11.5 Crash test script (40 min)

**📋 COPY → `benchmark/crash_test.ps1`**:
```powershell
# Usage: .\benchmark\crash_test.ps1
# Returns exit 0 on success, non-zero on failure
$ErrorActionPreference = "Stop"
$root = "F:\app\job\velosearch"
$dataDir = Join-Path $root "testdata\crash_$([guid]::NewGuid().Guid)"
New-Item -ItemType Directory -Force -Path $dataDir | Out-Null

try {
    Push-Location $root

    # 1. Start the server
    Write-Host "[1/4] Starting server..."
    $server = Start-Process -FilePath "go" `
        -ArgumentList "run","./cmd/server","-data-dir=$dataDir","-addr=:50052" `
        -PassThru `
        -RedirectStandardOutput "$dataDir\server.log" `
        -RedirectStandardError "$dataDir\server.err"
    Start-Sleep -Seconds 3

    # 2. Insert 1000
    Write-Host "[2/4] Writing 1000 vectors..."
    go run ./benchmark/crash_client -addr=localhost:50052 -mode=write -n=1000
    if ($LASTEXITCODE -ne 0) { throw "write failed" }

    # 3. SIGKILL
    Write-Host "[3/4] Killing server..."
    Stop-Process -Id $server.Id -Force
    Start-Sleep -Seconds 1

    # 4. Restart + verify
    Write-Host "[4/4] Restart + verify..."
    $server2 = Start-Process -FilePath "go" `
        -ArgumentList "run","./cmd/server","-data-dir=$dataDir","-addr=:50052" `
        -PassThru `
        -RedirectStandardOutput "$dataDir\server2.log" `
        -RedirectStandardError "$dataDir\server2.err"
    Start-Sleep -Seconds 4
    go run ./benchmark/crash_client -addr=localhost:50052 -mode=verify -n=1000
    $verifyExit = $LASTEXITCODE
    Stop-Process -Id $server2.Id -Force
    if ($verifyExit -ne 0) { throw "verify failed" }
    Write-Host "OK ✅"
}
finally {
    Pop-Location
    Remove-Item -Recurse -Force $dataDir -ErrorAction SilentlyContinue
}
```

### 11.6 10 次循环 (10 min)

**▶️ RUN:**
```powershell
1..10 | ForEach-Object {
    Write-Host "`n=== Run $_ ==="
    & .\benchmark\crash_test.ps1
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Run $_ FAILED"
        break
    }
}
```

**目标: 10/10 全过。**

---

## 🏁 Wrap-up (20 min)

```powershell
git add .
git commit -m "feat(day-11): tombstone delete + 10x crash recovery test (C3)"
```

**Today's checklist:**
- [ ] **C3: 10/10 crash test pass — ✅ / ❌**

---
---

# Day 12 — 2026-05-22 (周五 / Fri) — ann-benchmarks Adapter

## 今日目标
- Python adapter 接入 ann-benchmarks
- 后台启动 SIFT-1M 跑分
- Block 2 写 README 框架 / proto 注释 / 面试答案大纲

## Block 1 (90 min) — Adapter

### 12.1 Clone ann-benchmarks (10 min)

**▶️ RUN (PowerShell):**
```powershell
cd F:\app\job
git clone https://github.com/erikbern/ann-benchmarks.git
cd ann-benchmarks
pip install -r requirements.txt
```

### 12.2 Adapter module (80 min)

**目标文件:** `ann-benchmarks/ann_benchmarks/algorithms/velosearch/module.py`

**🤖 LLM PROMPT:**
```
The ann-benchmarks framework (https://github.com/erikbern/ann-benchmarks) expects each
algorithm to provide a Python class implementing BaseANN:

class BaseANN:
    def fit(self, X)                  # X is np.ndarray (n, dim); build index
    def set_query_arguments(self, ef_search)
    def query(self, v, n)             # return top-n nearest neighbor IDs for query v
    def get_memory_usage(self)        # return MB

My algorithm is a Go gRPC server (velosearch) on localhost:50051 with this proto:
<paste proto/velosearch.proto>

Write ann_benchmarks/algorithms/velosearch/module.py that:

class VeloSearch(BaseANN):
    def __init__(self, metric: str, m: int, ef_construction: int):
        self._m = m
        self._efc = ef_construction
        self._metric = "L2" if metric == "euclidean" else "Cosine"
        self._server = None
        self._channel = None
        self._stub = None
        self._collection = "ann"
        self._ef_search = 50

    def fit(self, X):
        # 1. tempdir for data
        # 2. subprocess.Popen Go server: go run github.com/.../cmd/server -addr=:50051 -data-dir=tmpdir
        #    (assume `go` is on PATH; the server binary path is the velosearch repo we cloned)
        # 3. wait 2s for startup
        # 4. open gRPC channel + stub
        # 5. CreateCollection (dim=X.shape[1], capacity=X.shape[0]+1000, ...)
        # 6. Insert in batches of 1000 (await each batch's response)

    def set_query_arguments(self, ef_search: int):
        self._ef_search = ef_search

    def query(self, v, n):
        # call Search RPC, return list of hit.id

    def get_memory_usage(self):
        # use psutil.Process(self._server.pid).memory_info().rss / 1024  (in KB, the framework expects KB)

    def __del__(self):
        # gracefully terminate Popen, clean tempdir

Also write the entry to add in ann-benchmarks/algos.yaml under section sift-128-euclidean:
  velosearch:
    docker-tag: ann-benchmarks-velosearch
    module: ann_benchmarks.algorithms.velosearch
    constructor: VeloSearch
    base-args: []
    run-groups:
      base:
        args: ["euclidean", 16, 200]   # metric, M, efConstruction
        query-args: [[50, 100, 200, 400]]

Note: you need to install grpcio and grpcio-tools in the ann-benchmarks env,
and generate Python stubs from velosearch.proto. Provide the python -m grpc_tools.protoc
command at the end.
```

**▶️ RUN — 生成 Python proto stubs:**
```powershell
cd F:\app\job\ann-benchmarks
pip install grpcio grpcio-tools psutil
python -m grpc_tools.protoc -I=../velosearch --python_out=ann_benchmarks/algorithms/velosearch --grpc_python_out=ann_benchmarks/algorithms/velosearch ../velosearch/proto/velosearch.proto
```

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — Start Run + Parallel Work

### 12.3 启动 SIFT 跑分 (5 min)

**▶️ RUN:**
```powershell
cd F:\app\job\ann-benchmarks
python run.py --algorithm velosearch --dataset sift-128-euclidean --local
# Runs 30-60 minutes in the background
```

### 12.4 等的时候做 (85 min)

**并行任务三选 (避免空等):**

1. **README 框架** (45 min) — 见 Day 14 详细要求,先把骨架写出来,benchmark 数据留空
2. **Proto 文件注释** (20 min) — 每个 message 加一行 doc comment
3. **NOTES.md 面试答案大纲** (20 min) — 把 10 个白板问题各写 60 秒大纲

### 12.5 看结果 (如果跑完)

**▶️ RUN:**
```powershell
python plot.py --dataset sift-128-euclidean
# Generates results/sift-128-euclidean.png
```

VeloSearch 的点应该在 recall-vs-QPS Pareto 前沿附近(不必超过 hnswlib,在同 ballpark 就赢了)。

---

## 🏁 Wrap-up (20 min)

```powershell
git add ann-benchmarks/ann_benchmarks/algorithms/velosearch/
git commit -m "feat(day-12): ann-benchmarks adapter for VeloSearch"
```

---
---

# Day 13 — 2026-05-23 (周六 / Sat) — GIST-1M + Report + Dockerfile

## 今日目标
- SIFT 结果完成 + GIST-1M 跑分
- 写 BENCHMARK.md
- 多阶段 Dockerfile,镜像 < 30 MB

## Block 1 (90 min) — GIST + Report

### 13.1 SIFT 结果截图 (15 min)

```powershell
cd F:\app\job\ann-benchmarks
python plot.py --dataset sift-128-euclidean
Copy-Item results\sift-128-euclidean.png F:\app\job\velosearch\docs\img\
```

### 13.2 GIST-1M (10 min 启动)

```powershell
# Dataset auto-downloads (~3.6 GB)
python run.py --algorithm velosearch --dataset gist-960-euclidean --local
# Runs 1-2 hours in the background
```

### 13.3 BENCHMARK.md (65 min)

**📋 COPY → `docs/BENCHMARK.md`** (模板,数据你填):
```markdown
# VeloSearch Benchmark Report

## Test Environment

- CPU: <你的 CPU 型号 e.g. "AMD Ryzen 7 5800X3D">
- RAM: <e.g. "32 GB DDR4-3600">
- OS: Windows 11
- Go: 1.26.3
- Storage: <e.g. "NVMe SSD">

All queries are single-threaded; no concurrent reads.

## SIFT-1M (128-d, Euclidean)

Dataset: 1M base vectors, 10K queries, 100 ground-truth neighbors per query.

| Engine | Config | Recall@10 | Mean Latency | P99 Latency | QPS |
|--------|--------|-----------|--------------|-------------|-----|
| VeloSearch | M=16, efC=200, efS=100 | <X.XX> | <X> µs | <X> µs | <X> |
| hnswlib (ref) | M=16, efC=200, efS=100 | <X.XX> | <X> µs | <X> µs | <X> |

![SIFT-1M Recall vs QPS](img/sift-128-euclidean.png)

## GIST-1M (960-d, Euclidean)

Dataset: 1M base vectors, 1K queries, 100 ground-truth neighbors per query.

[同上表]

## Observations

- Higher dimensionality (960 vs 128) drops recall ~3-5% at same efSearch parameters,
  consistent with the curse of dimensionality.
- VeloSearch memory footprint on SIFT-1M: ~700 MB (contiguous vectorStore) + ~150 MB
  (graph + Go map overhead) = ~850 MB total.
- HNSW deletion is currently tombstone-based; rebuilding the index to reclaim
  tombstone space is on the roadmap.
```

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — Dockerfile + Image

### 13.4 Multi-stage Dockerfile (45 min)

**📋 COPY → `deploy/Dockerfile`**:
```dockerfile
# Stage 1: build
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /out/velosearch \
    ./cmd/server

# Stage 2: minimal runtime (distroless, non-root)
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/velosearch /velosearch
EXPOSE 50051
USER nonroot
ENTRYPOINT ["/velosearch"]
CMD ["-addr=:50051", "-data-dir=/data"]
```

### 13.5 Build + verify (25 min)

**▶️ RUN:**
```powershell
docker build -t velosearch:dev -f deploy/Dockerfile .
docker images velosearch:dev
# Expect image size < 30 MB
docker run --rm -p 50051:50051 velosearch:dev &
Start-Sleep -Seconds 3
grpcurl -plaintext localhost:50051 list
```

### 13.6 .dockerignore (10 min)

**📋 COPY → `.dockerignore`**:
```
.git
.github
benchmark/sift/data
benchmark/gist/data
*.prof
*.pprof
data
testdata
docs
ann-benchmarks
*.md
LICENSE
```

### 13.7 GIST 数据 (10 min)

如果 GIST 跑完了,把数字填进 BENCHMARK.md。否则明天补。

---

## 🏁 Wrap-up (20 min)

```powershell
git add .
git commit -m "feat(day-13): GIST-1M benchmark, BENCHMARK.md, multi-stage Dockerfile (~28MB)"
```

---
---

# Day 14 — 2026-05-24 (周日 / Sun) — docker-compose + README + GitHub Release

## 今日目标
- docker-compose 一键起 velosearch + prometheus + grafana
- 写完 README
- 推 GitHub,v0.1.0 tag,Pin 到 profile

## Block 1 (90 min) — Metrics + Compose

### 14.1 Prometheus metrics (45 min)

**📋 COPY → `internal/grpcserver/metrics.go`**:
```go
package grpcserver

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	SearchLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "velosearch_search_latency_seconds",
		Help:    "Search RPC latency",
		Buckets: prometheus.ExponentialBuckets(0.0001, 2, 16),
	}, []string{"collection"})

	InsertCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "velosearch_inserts_total",
		Help: "Total vectors inserted",
	}, []string{"collection"})

	CollectionSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "velosearch_collection_vectors",
		Help: "Current vector count per collection",
	}, []string{"collection"})
)
```

**▶️ RUN:**
```powershell
go get github.com/prometheus/client_golang
```

**🧩 SKELETON: 在 `internal/grpcserver/server.go` 的 Search handler 加 timer:**
```go
import "github.com/prometheus/client_golang/prometheus"

func (s *Server) Search(...) (..., error) {
	timer := prometheus.NewTimer(SearchLatency.WithLabelValues(req.Collection))
	defer timer.ObserveDuration()
	// ... rest of handler
}
```
对 Insert handler 加 `InsertCounter.WithLabelValues(req.Collection).Add(float64(len(req.Items)))`。

**🧩 SKELETON: 在 `cmd/server/main.go` 加 HTTP /metrics:**
```go
import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// in main(), before the graceful-shutdown wait:
go func() {
	http.Handle("/metrics", promhttp.Handler())
	slog.Info("metrics http listening", "addr", ":9090")
	_ = http.ListenAndServe(":9090", nil)
}()
```

### 14.2 docker-compose (45 min)

**📋 COPY → `deploy/docker-compose.yml`**:
```yaml
version: "3.8"
services:
  velosearch:
    build:
      context: ..
      dockerfile: deploy/Dockerfile
    ports:
      - "50051:50051"
      - "9090:9090"
    volumes:
      - velosearch-data:/data

  prometheus:
    image: prom/prometheus:latest
    ports: ["9091:9090"]
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana:latest
    ports: ["3000:3000"]
    environment:
      - GF_AUTH_ANONYMOUS_ENABLED=true
      - GF_AUTH_ANONYMOUS_ORG_ROLE=Admin

volumes:
  velosearch-data:
```

**📋 COPY → `deploy/prometheus.yml`**:
```yaml
global:
  scrape_interval: 5s
scrape_configs:
  - job_name: velosearch
    static_configs:
      - targets: ['velosearch:9090']
```

**▶️ RUN:**
```powershell
cd deploy
docker-compose up -d
# After a few seconds, open Grafana at http://localhost:3000 (admin/admin)
```

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — README + GitHub

### 14.3 README (60 min)

**目标文件:** `README.md` (替换 Day 1 占位版)

**🤖 LLM PROMPT:**
```
Write a README.md for github.com/zhangchuqi1998/velosearch — an HNSW-based vector
search engine in Go.

Audience: senior engineers / recruiters, ~60 seconds to scan.

Sections:
1. Title + shields.io badges: Go 1.22+, License MIT, CI status placeholder
2. One-line pitch: "Open-source vector search engine built around an HNSW index
   implemented from scratch in Go. v0.1 — for benchmarking and learning."
3. Highlights (4-5 bullets covering: HNSW from scratch, gRPC API, WAL durability,
   sub-Xms P99 on SIFT-1M, ann-benchmarks compatible)
4. Quickstart:
   ```bash
   docker run -p 50051:50051 ghcr.io/zhangchuqi1998/velosearch:v0.1.0
   grpcurl -plaintext -d '{"name":"test","dim":128,"metric":"METRIC_L2","m":16,"ef_construction":200}' \
     localhost:50051 velosearch.v1.VectorSearch/CreateCollection
   ```
5. Architecture — mermaid diagram:
   Client → gRPC server → Collection Manager → HNSW Index → vectorStore + Graph
                                        ↓
                                     WAL writer (fsync per op)
6. Benchmarks — short table from BENCHMARK.md (SIFT row + GIST row), link to docs/BENCHMARK.md
7. FAQ:
   - "Why HNSW (vs IVF, vs flat)?"
   - "How does delete work?" (tombstone + future rebuild)
   - "Is this production-ready?" (Honest: "v0.1 — for learning and benchmarking. Use Qdrant or Pinecone for production.")
8. Roadmap — bullets: PQ compression, replication, metadata filtering
9. License: MIT

Tone: technical, honest, no marketing. Code blocks with language tags. Mermaid for diagrams.
```

**⚠️ "v0.1 — for learning" 这个 disclaimer 是关键 — 别假装 production-ready。**

### 14.4 LICENSE (5 min)

**📋 COPY → `LICENSE`** (从 https://opensource.org/license/mit/ 抄,把 `[year]` 改成 `2026`,`[fullname]` 改成 `Chuqi Zhang`)

### 14.5 GitHub Release (25 min)

**▶️ RUN:**
```powershell
# 1. In a browser: github.com/zhangchuqi1998 -> New repository -> name=velosearch, public, do not init README/license
# 2. Command line:
cd F:\app\job\velosearch
git remote add origin git@github.com:zhangchuqi1998/velosearch.git
git branch -M main
git push -u origin main

# 3. Tag + release
git tag -a v0.1.0 -m "v0.1.0: HNSW vector search, gRPC API, WAL persistence, ann-benchmarks compatible"
git push origin v0.1.0

# 4. In a browser:
#    - GitHub repo -> Releases -> Edit v0.1.0 -> write release notes (highlight SIFT numbers + image size)
#    - Profile -> Customize pins -> select velosearch
```

---

## 🏁 Wrap-up (20 min)

```powershell
git add .
git commit -m "feat(day-14): Prometheus metrics, docker-compose, README, LICENSE"
git push
```

🚀 **Project v0.1 shipped.**

---
---

# Day 15 — 2026-05-25 (周一 / Mon) — 简历回填 + 面试准备

## Block 1 (90 min) — Resume

### 15.1 Backfill 实测数字 (30 min)

打开 `F:\app\job\generate_resume.py`,找 VeloSearch 部分,把这几个 placeholder 替换:

| 占位 | 替换为 |
|------|------|
| `Recall@10 of 95%` | Day 7 实测 |
| `P99 search latency under 15ms` | Day 7 实测 |
| `3K+ inserts/sec` | Day 6 实测(warmed-up,单线程)|

**▶️ RUN:**
```powershell
cd F:\app\job
python generate_resume.py
# Check the PDF at C:\Users\Aaron\Downloads\AaronZhang_Resume_v3.pdf
```

### 15.2 LinkedIn / GitHub profile (30 min)
- LinkedIn About 加一段 "Built VeloSearch..."
- LinkedIn Featured: pin velosearch repo 链接
- GitHub profile README (`zhangchuqi1998/zhangchuqi1998` 这个特殊 repo): 1 段简介

### 15.3 自查 (30 min)
- 简历数字 = BENCHMARK.md 数字
- GitHub 链接打开能进项目
- README 没拼写错误
- LinkedIn URL 简历一致

---

## ☕ Break (20 min)

---

## Block 2 (90 min) — Whiteboard Practice

### 15.4 10 题口头答题 (60 min)

**📚 REF — 录音 + 答 + 回放:**

1. Walk me through how HNSW works.
2. Why HNSW over IVF for vector search?
3. What does `efSearch` control? How did you pick it?
4. Recall dropped from SIFT to GIST — why?
5. How does delete work in HNSW? Why is it hard?
6. Walk me through what happens when a client calls Search.
7. How does your WAL handle a crash mid-write?
8. If you had two more weeks, what would you build next?
9. What's the worst bug you hit?
10. What did you learn?

每个 60-90 秒,**不要看笔记**。

**🤖 LLM PROMPT (打磨答案):**
```
Here's my 60-second answer to "<question>":
<transcript of your recording>

Critique it like a senior engineer hiring for backend roles:
- Did I miss any key technical point?
- Did I oversell/undersell?
- Is the structure clear (problem → approach → tradeoff)?
- One concrete improvement.
```

### 15.5 Mock 系统设计 (30 min)

**🤖 LLM PROMPT:**
```
Give me a system design interview question that builds on the VeloSearch project I built
(HNSW vector search engine). The question should ask me to scale it 10× (10M vectors,
10× QPS) and add features I didn't build (replication, sharding, metadata filters).

After my answer (which I'll paste), critique it.
```

口头答 15 分钟,粘贴 transcript 让 LLM critique。

---

## 🏁 Wrap-up (20 min)

```powershell
cd F:\app\job
git add generate_resume.py
git commit -m "docs: update VeloSearch resume bullets with measured metrics"
```

**今日 + 整体 checklist:**
- [ ] 简历数字全是实测
- [ ] GitHub repo public, README 完整, pinned to profile
- [ ] LinkedIn / GitHub profile 更新
- [ ] 10 个问题都口头答过一遍
- [ ] 明天 Day 16 投出第一批

---
---

# 通用资源 / Resources

## Papers + Blogs
- HNSW: https://arxiv.org/abs/1603.09320
- Pinecone HNSW: https://www.pinecone.io/learn/series/faiss/hnsw/
- Reference C++ impl: https://github.com/nmslib/hnswlib
- Vector DB landscape: https://thedataquarry.com/posts/vector-db-4/

## Datasets
- SIFT-1M: http://corpus-texmex.irisa.fr/
- GIST-1M: 同上
- ann-benchmarks: https://github.com/erikbern/ann-benchmarks

## Go tools
- SIMD: github.com/viterin/vek
- gRPC Go: https://grpc.io/docs/languages/go/quickstart/
- Profiling: `go tool pprof`, `go test -bench`

---

# 常见坑 / Common Pitfalls

1. **别用 interface{} 的优先队列** — GC 压力毁掉 latency。用 `[]Candidate`。
2. **先正确,后优化** — Recall first, latency second。
3. **改一处测一次** — 一半"优化"是 no-op 或倒退。没数据别 commit。
4. **别跳 Algorithm 4** — 简单选邻居 recall 卡 80%。
5. **测试固定 RNG seed** — `rand.NewSource(42)`,否则 flaky test。
6. **别每 insert fsync** — 吞吐崩。批量或 async。
7. **不要在不同机器上对比** — ann-benchmarks 在统一环境跑。
8. **`go test -race`** — 加 race flag,并发 bug 早发现。
9. **简历数字必须是 main 分支跑出来的** — 不要"调试版本"数字混进去。
10. **每天 push** — 别攒一晚上提交。

---

# WEEKLY_LOG.md 模板

```markdown
## Day N — YYYY-MM-DD
- 完成 / Done:
  - ...
- 卡点 / Blockers:
  - ...
- 当前指标 / Current metrics:
  - Recall@10: ...
  - P99: ...
  - Insert throughput: ...
- 明天先做 / Tomorrow first task:
  - ...
```
