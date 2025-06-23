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

## 📊 Benchmark Results

This section provides benchmark data for various MapReduce use cases. All implementations are written in idiomatic Go with concurrency optimizations where applicable.

### 🧪 Word Count

**Average of 5 runs**  
**Dataset:** 500 lines × 600,000 words each

| Version                                                  | Average Time       |
|-----------------------------------------------------------|--------------------|
| 🟦 Synchronous                                             | 12.6483617s        |
| 🟩 Concurrent                                              | 6.63623224s        |
| 🟨 Concurrent with Worker Pool (Map Phase)                 | 6.84312916s        |
| 🟪 Concurrent with Worker Pool + Buffered Channel + Mutex | 5.79071534s        |

---

### 📝 Observations

- 🔁 **Concurrent implementations consistently outperform the synchronous version**, confirming Go’s concurrency model is effective for large-scale data processing.
- ⚙️ **Raw concurrent version** (one goroutine per line) is fast but can lead to too many goroutines, which stresses the scheduler and memory.
- 👷 **Worker pool in map phase** adds control over goroutine count, providing a good trade-off between speed and resource usage.
- 🚀 **Worker pool + buffered channel + mutex reduce phase** is fastest. This setup benefits from:
  - Batching work to limit goroutines,
  - Buffered channels to reduce producer-consumer blocking,
  - Parallel reduction with `sync.Mutex` to safely aggregate results.

---

📌 *Tip: Tune chunk sizes and buffer capacities based on available CPU cores and workload characteristics for optimal performance.*

