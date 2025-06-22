package main

import (
	"strings"
	"sync"
	"math/rand"
	"time"
	"fmt"
	"runtime"
)

func mapFunc(line string) map[string]int {
	result := make(map[string]int)
	words := strings.Fields(line)
	for _, word := range words {
		result[word]++
	}
	return result
}

func reduce(ch <-chan map[string]int) map[string]int {
	final := make(map[string]int)
	for partial := range ch {
		for word, count := range partial {
			final[word] += count
		}
	}
	return final
}

func reduceMutex(ch <-chan map[string]int) map[string]int {
	final := make(map[string]int)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for partial := range ch {
		wg.Add(1)
		go func(part map[string]int) {
			defer wg.Done()
			for word, count := range part {
				mu.Lock()
				final[word] += count
				mu.Unlock()
			}
		}(partial)
	}
	wg.Wait()
	return final
}

func concurrent(lines []string) map[string]int{
	var wg sync.WaitGroup
	mapChan := make(chan map[string]int)

	for _, line := range lines {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			mapChan <- mapFunc(l)
		}(line)
	}

	go func() {
		wg.Wait()
		close(mapChan)
	}()

	return reduce(mapChan)
}

func concurrentBufferedChan(lines []string) map[string]int{
	var wg sync.WaitGroup
	mapChan := make(chan map[string]int, len(lines))

	for _, line := range lines {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			mapChan <- mapFunc(l)
		}(line)
	}

	go func() {
		wg.Wait()
		close(mapChan)
	}()

	return reduce(mapChan)
}

func concurrentMutexReduce(lines []string) map[string]int{
	var wg sync.WaitGroup
	mapChan := make(chan map[string]int, len(lines))

	for _, line := range lines {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			mapChan <- mapFunc(l)
		}(line)
	}

	go func() {
		wg.Wait()
		close(mapChan)
	}()

	return reduceMutex(mapChan)
}

func synchronous(lines []string) map[string]int {
	result := make(map[string]int)
	for _, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			result[word]++
		}
	}
	return result
}

func concurrentWorkerPool(lines []string, workers, chunkSize int) map[string]int{
	mapChan := make(chan map[string]int, 1024)
	chunks := chunkLines(lines, chunkSize)
	jobChan := make(chan []string, len(chunks))
	var wg sync.WaitGroup

	for _, chunk := range chunks {
		jobChan <- chunk
	}
	close(jobChan)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for lines := range jobChan {
				result := make(map[string]int)
				for _, line := range lines {
					for word, count := range mapFunc(line) {
						result[word] += count
					}
				}
				mapChan <- result
			}
		}()
	}

	go func() {
		wg.Wait()
		close(mapChan)
	}()

	return reduce(mapChan)
}

func generateLines(lineCount, wordsPerLine int) []string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	lines := make([]string, lineCount)

	for i := 0; i < lineCount; i++ {
		var words []string
		for j := 0; j < wordsPerLine; j++ {
			wordLen := 3
			var sb strings.Builder
			for k := 0; k < wordLen; k++ {
				sb.WriteByte(letters[rand.Intn(len(letters))])
			}
			words = append(words, sb.String())
		}
		lines[i] = strings.Join(words, " ")
	}

	return lines
}

func chunkLines(lines []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(lines); i += chunkSize {
		end := i + chunkSize
		if end > len(lines) {
			end = len(lines)
		}
		chunks = append(chunks, lines[i:end])
	}
	return chunks
}

func compareMap(a, b map[string]int) bool {
	if len(a) != len(b){
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func calculateAverage(n, lines, wordInLine int) {
	var totalSynchronous time.Duration
	var totalConcurrent time.Duration
	var totalBuffered time.Duration
	var totalWorker time.Duration
	var totalMutex time.Duration

	for i := 0; i < n; i++ {
		fmt.Printf("Run #%d\n", i+1)
		lines := generateLines(lines, wordInLine)

		start := time.Now()
		result1 := synchronous(lines)
		elapsed := time.Since(start)
		totalSynchronous += elapsed

		start = time.Now()
		result2 := concurrent(lines)
		elapsed = time.Since(start)
		totalConcurrent += elapsed

		start = time.Now()
		result3 := concurrentBufferedChan(lines)
		elapsed = time.Since(start)
		totalBuffered += elapsed

		start = time.Now()
		result4 := concurrentWorkerPool(lines, runtime.NumCPU(), 25)
		elapsed = time.Since(start)
		totalWorker += elapsed

		start = time.Now()
		result5 := concurrentMutexReduce(lines)
		elapsed = time.Since(start)
		totalMutex += elapsed

		isEqual := compareMap(result1, result2) &&
			compareMap(result1, result3) &&
			compareMap(result1, result4) &&
			compareMap(result1, result5)

		if isEqual {
			fmt.Println("All results are equal")
		} else {
			fmt.Println("Results are not equal")
		}

		fmt.Println("---")
	}

	fmt.Printf("\n==== AVERAGE RESULTS AFTER %v RUNS ====\n", n)
	fmt.Printf("Synchronous:                 %v\n", totalSynchronous/time.Duration(n))
	fmt.Printf("Concurrent:                  %v\n", totalConcurrent/time.Duration(n))
	fmt.Printf("Concurrent Buffered Channel: %v\n", totalBuffered/time.Duration(n))
	fmt.Printf("Concurrent with Worker Pool: %v\n", totalWorker/time.Duration(n))
	fmt.Printf("Concurrent Mutex Reduce:     %v\n", totalMutex/time.Duration(n))
}
