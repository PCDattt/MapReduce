package main

import (
	"strings"
	"sync"
	"math/rand"
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

func wordCountConcurrency(lines []string) {
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

	reduce(mapChan)
}

func wordCountConcurrencyWithMutexReduce(lines []string) {
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

	reduceMutex(mapChan)

	// for word, count := range finalCount {
	// 	fmt.Printf("%s: %d\n", word, count)
	// }
}

func wordCount(lines []string) {
	result := make(map[string]int)
	for _, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			result[word]++
		}
	}
	// for word, count := range result {
	// 	fmt.Printf("%s: %d\n", word, count)
	// }
}

func wordCountWithWorkerPool(lines []string, workers, chunkSize int) {
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

	reduce(mapChan)
}



func generateLines(lineCount, wordsPerLine int) []string {
	rand.Seed(time.Now().UnixNano())
	const letters = "abcdefghijklmnopqrstuvwxyz"
	lines := make([]string, lineCount)

	for i := 0; i < lineCount; i++ {
		var words []string
		for j := 0; j < wordsPerLine; j++ {
			wordLen := rand.Intn(5) + 3 // word length: 3–7
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
