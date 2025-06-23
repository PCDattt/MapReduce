package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"time"
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
	var mu sync.Mutex
	var wg sync.WaitGroup

	workerCount := 10

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for partial := range ch {
				for word, count := range partial {
					mu.Lock()
					final[word] += count
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	return final
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

func concurrent(lines []string) map[string]int {
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

func concurrentWorkerPoolMap(lines []string, workers, chunkSize int) map[string]int {
	mapChan := make(chan map[string]int)
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

func concurrentBufferedChanMutexReduce(lines []string) map[string]int {
	mapChan := make(chan map[string]int, 10)
	lineChan := make(chan string, len(lines))
	var wg sync.WaitGroup
	var wgLine sync.WaitGroup
	worker := 10

	wgLine.Add(1)
	go func() {
		wgLine.Done()
		for _, line := range lines {
			lineChan <- line
		}
	}()
	
	go func() {
		wgLine.Wait()
		close(lineChan)
	}()

	for i := 0; i < worker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lineChan {
				mapChan <- mapFunc(line)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(mapChan)
	}()

	return reduceMutex(mapChan)
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
	if len(a) != len(b) {
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
		result3 := concurrentWorkerPoolMap(lines, runtime.NumCPU(), 20)
		elapsed = time.Since(start)
		totalWorker += elapsed

		start = time.Now()
		result4 := concurrentBufferedChanMutexReduce(lines)
		elapsed = time.Since(start)
		totalMutex += elapsed

		isEqual := compareMap(result1, result2) &&
			compareMap(result1, result3) &&
			compareMap(result1, result4)

		if isEqual {
			fmt.Println("All results are equal")
		} else {
			fmt.Println("Results are not equal")
		}

		fmt.Println("---")
	}

	fmt.Printf("\n==== AVERAGE RESULTS AFTER %v RUNS, %v LINES, %v WORDS IN LINE ====\n", n, lines, wordInLine)
	fmt.Printf("Synchronous:                                           %v\n", totalSynchronous/time.Duration(n))
	fmt.Printf("Concurrent:                                            %v\n", totalConcurrent/time.Duration(n))
	fmt.Printf("Concurrent with Worker Pool Map:                       %v\n", totalWorker/time.Duration(n))
	fmt.Printf("Concurrent with Worker Pool, Buffered Channel, Mutex:  %v\n", totalMutex/time.Duration(n))
}
