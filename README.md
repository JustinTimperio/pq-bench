# PQ-Bench

## Background
This repository contains benchmarks for various priority queue implementations in Go, Rust, Zig, and C++. All tests are single threaded and designed to try and reproduce roughly the same workload across all implementations. While these benchmarks are not perfect, they should give a rough idea of the performance of each implementation and how they might work in a real-world scenario. 

### Sister Repositories
- [fibheap (Fibonacci Heaps)](https://github.com/JustinTimperio/fibheap)
- [gpq (Go Priority Queue)](https://github.com/JustinTimperio/gpq)
- [rpq (Rust Priority Queue)](https://github.com/JustinTimperio/rpq)


## Benchmarks
![Benchmark](./docs/Time-Spent-(seconds)-vs-Implementation.png)

## Features
| Feature          | GPQ | RPQ | Go Binary Heap | Zig Binary Heap | C++ Binary Heap | Rust Binary Heap |
|------------------|-----|-----|----------------|-----------------|-----------------|------------------|
| Enqueue          | ✅   | ✅   | ✅              | ✅               | ✅               | ✅                |
| Dequeue          | ✅   | ✅   | ✅              | ✅               | ✅               | ✅                |
| Disk Cache       | ✅   | 🚧  | ❌              | ❌               | ❌               | ❌                |
| Mutable Priority | ✅   | ✅   | ❌              | ❌               | ❌               | ❌                |
| Timeouts         | ✅   | ✅   | ❌              | ❌               | ❌               | ❌                |
