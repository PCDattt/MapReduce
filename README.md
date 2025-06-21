# 🗂 Go MapReduce Project

A collection of high-performance, idiomatic Go implementations of the **MapReduce** pattern — designed to demonstrate mastery of concurrency, channel communication, and real-world data processing in Go.

This project is crafted as part of a portfolio to showcase backend and systems programming skills using Go — with a focus on correctness, performance, and idiomatic code structure.

---

## 🧠 Versions Included

This repository contains two primary versions of the MapReduce pattern:

- **🔹 Synchronous version:**  
  A straightforward, single-threaded implementation with sequential map and reduce logic. Useful as a baseline for benchmarking and correctness.

- **🔹 Concurrent version:**  
  A fully concurrent, idiomatic Go implementation using goroutines, channels, `sync.WaitGroup`, and `sync.Mutex` — scalable and efficient on large inputs.

These versions allow direct performance comparison and highlight how Go handles concurrency with elegance and minimal overhead.

---

## 🚀 Features

- ✅ **Pure Go** — built entirely with the standard library, no third-party dependencies
- ✅ **Idiomatic concurrency** — uses goroutines, channels, and synchronization patterns idiomatically
- ✅ **Real-world use cases** — word counting, log analysis, sales data aggregation, text indexing
- ✅ **Benchmark-ready** — compare sync vs concurrent performance on large files
- ✅ **Flexible input** — supports file-based or in-memory data
- ✅ **Buffered channel support** — experiment with backpressure and throughput tuning

---
