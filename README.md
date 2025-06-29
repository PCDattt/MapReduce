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

The Word Count problem is a classic MapReduce use case where the goal is to count how many times each word appears in a large collection of text. It's a simple but powerful pattern to demonstrate data-parallel processing using mapping (tokenizing lines into words) and reducing (aggregating word frequencies).

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

### 🧪 Log Analysis

The Log Analysis problem involves parsing large volumes of server logs to extract insights like request counts per IP, HTTP method, endpoint, status code, and response size distribution. It simulates real-world observability and monitoring tasks, and serves as a great benchmark for comparing the performance of synchronous vs concurrent data processing pipelines.

**Average of 5 runs**  
**Dataset:** 25 files × 800,000 lines per file, 1000 IPs, 1000 Endpoints, 4 Methods, 4 Status Codes and size between 200 - 1600
| Version                                                  | Average Time       |
|-----------------------------------------------------------|--------------------|
| 🟦 Synchronous                                             | 23.906264377s        |
| 🟩 Concurrent                                              | 10.999082619s        |
| 🟪 Concurrent with Buffered Channel + Mutex | 11.577384654s        |

### 📝 Observations

- ✅ **Concurrent** version provides an average **~2x speedup** over the synchronous baseline.
- ⚙️ The **Concurrent with Buffered Channel and Mutex** version performs **slightly slower** than plain concurrent because:
  - It spawns multiple goroutines per result to merge counts in parallel.
  - It uses multiple `sync.Mutex` locks for different maps, which introduces overhead due to lock contention and goroutine scheduling.
- 🧵 **Synchronous** is the **simplest** and most predictable, but **does not scale well** with larger data volumes or multi-core CPUs.


📌 *Tip: Tune chunk sizes and buffer capacities based on available CPU cores and workload characteristics for optimal performance.*

