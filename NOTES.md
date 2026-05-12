# HNSW Reading Notes

> 在 Day 1 Block 2 读完论文后填这份笔记。
> Fill this in after reading the paper on Day 1 Block 2.

## Paper / 论文
- Malkov & Yashunin (2018), *"Efficient and robust approximate nearest neighbor search using Hierarchical Navigable Small World graphs"*
- https://arxiv.org/abs/1603.09320

## Parameters / 参数

| Param | 含义 / Meaning | Typical |
|-------|---------|---------|
| `M` | layer > 0 上每节点最大邻居数 / max neighbors per node on layers > 0 | 16 |
| `M_max` | = M | 16 |
| `M_max0` | layer 0 上每节点最大邻居数 / max neighbors on layer 0 | 2*M = 32 |
| `mL` | `1 / ln(M)` — 层级生成归一化因子 / level normalization | 0.361 (M=16) |
| `efConstruction` | 建图候选集大小 / build-time candidate set size | 200 |
| `efSearch` | 查询候选集大小,**runtime knob** | 50-400 |

## Algorithms (paper Section 4)

### Algorithm 1 — INSERT
- TODO: 自己用一两句话总结 / summarize in your own words

### Algorithm 2 — SEARCH-LAYER ⭐ 核心 / core
- TODO

### Algorithm 3 — SELECT-NEIGHBORS-SIMPLE
- **不要用 / Do not use.** Recall 卡在 ~80% 上不去。

### Algorithm 4 — SELECT-NEIGHBORS-HEURISTIC ⭐ 必须用 / required
- TODO: 启发式规则一句话 — accept candidate c only if no already-accepted neighbor is closer to c than c is to the query.

### Algorithm 5 — K-NN-SEARCH
- TODO

---

## 三个概念问题(不查资料口述)/ Three concept questions (answer without notes)

**1. 为啥 HNSW 是分层图而不是单层?**
TODO:

**2. 搜索时,什么时候从 greedy descent 切到 ef-bounded search?**
TODO:

**3. Algorithm 4 的邻居启发式防止了什么?**
TODO:
