# List 01 — Go Runtime

Part of the ML Infra Roadmap.

## What I'm learning

- GMP scheduler (G, M, P)
- Work-stealing, syscall, netpoller
- Channels (buffered, unbuffered, select, hchan)
- sync: Mutex, RWMutex, WaitGroup, Once, atomic
- Escape analysis
- GC (tri-color, GOGC)
- sync.Pool
- pprof (CPU, heap, goroutine, block)
- benchmark + benchstat
- trace

## Projects

All projects live in this repository.

- [player-service](./player-service/) — first microservice (gRPC + HTTP + Postgres + GMP + channels)
- `task-service` — coming next (+ sync)

## Flagship

`final-go-runtime` — the flagship of List 1. Combines everything: GMP, channels, sync, escape analysis, GC, pprof, benchmark, trace.

## Status

In progress.
