# VeloSearch

> Open-source vector search engine built around an **HNSW** index implemented from scratch in Go.

🚧 **v0.1 in development.** Target release: 2026-05-24.

## Status

| Component | Day | Status |
|-----------|-----|--------|
| Project scaffold | 1 | ✅ |
| Distance functions + heaps + Index struct | 2 | ⏳ |
| HNSW `searchLayer` + `Insert` | 3 | ⏳ |
| `Search` + brute-force baseline (Recall ≥ 90% on 10K) | 4 | ⏳ |
| SIFT-1M loader + 1M-vector build | 5 | ⏳ |
| Baseline benchmark + parameter sweep | 6 | ⏳ |
| Optimization (P99 < 15ms at recall ≥ 95%) | 7 | ⏳ |
| Protobuf schema + Collection Manager | 8 | ⏳ |
| gRPC API + integration tests | 9 | ⏳ |
| WAL persistence + replay | 10 | ⏳ |
| Tombstone delete + crash recovery | 11 | ⏳ |
| ann-benchmarks adapter + SIFT/GIST runs | 12-13 | ⏳ |
| Dockerfile + docker-compose + final README | 13-14 | ⏳ |
| GitHub v0.1.0 release | 14 | ⏳ |

## Quickstart

```bash
# After v0.1.0 release:
docker run -p 50051:50051 ghcr.io/zhangchuqi1998/velosearch:v0.1.0
```

## License

MIT (to be added in Day 14).
