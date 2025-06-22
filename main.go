package main

import (
	"fmt"
	"time"
	"runtime"
)

func main() {
	lines := generateLines(20_000, 1000)
	start := time.Now()
	wordCount(lines)
	fmt.Println("Synchronous Word Count:", time.Since(start))

	start = time.Now()
	wordCountConcurrency(lines)
	fmt.Println("Concurrency Word Count: ", time.Since(start))

	start = time.Now()
	wordCountWithWorkerPool(lines, runtime.NumCPU(), 500)
	fmt.Println("Concurrency Word Count with worker : ", time.Since(start))

	// start = time.Now()
	// wordCountConcurrencyWithMutexReduce(lines)
	// fmt.Println("Concurrency Word Count with reduce mutex duration: ", time.Since(start))
}
