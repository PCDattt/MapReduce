package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
	"bufio"
	"strings"
	"sync"
)

type AnalysisResult struct {
	IPCount map[string]int
	MethodCount map[string]int
	StatusCount map[string]int
	EndpointCount map[string]int
	SizeBucketCount map[string]int
}

func NewAnalysisResult() AnalysisResult {
	return AnalysisResult {
		IPCount: make(map[string]int),
		MethodCount: make(map[string]int),
		StatusCount: make(map[string]int),
		EndpointCount: make(map[string]int),
		SizeBucketCount: make(map[string]int),
	}
}

func CompareAnalysisResult(a, b AnalysisResult) bool {
	return compareMap(a.IPCount, b.IPCount) &&
			compareMap(a.EndpointCount, b.EndpointCount) &&
			compareMap(a.MethodCount, b.MethodCount) &&
			compareMap(a.StatusCount, b.StatusCount) &&
			compareMap(a.SizeBucketCount, b.SizeBucketCount)
}

func LogAnalysis(n int) {
	numFiles := 100
	linesPerFile := 10000
	numIPs := 5
	numEndpoints := 5
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	statusCodes := []int{100, 200, 300, 400}
	minSize, maxSize := 200, 1600
	outputDir := "logs"

	var totalSynchronous time.Duration
	var totalConcurrent time.Duration

	for i := 0; i < n; i++ {
		fmt.Printf("Run #%d\n", i+1)

		err := GenerateLogFiles(numFiles, linesPerFile, numIPs, numEndpoints, methods, statusCodes, minSize, maxSize, outputDir)
		if err != nil {
			panic(err)
		}

		start := time.Now()
		result1 := AnalysisSynchronous(numFiles, outputDir)
		elapsed := time.Since(start)
		totalSynchronous += elapsed

		start = time.Now()
		result2 := AnalysisConcurrent(numFiles, outputDir)
		elapsed = time.Since(start)
		totalConcurrent += elapsed

		isEqual := CompareAnalysisResult(result1, result2)
		
		if isEqual {
			fmt.Println("Results are equal")
		} else {
			fmt.Println("Results are not equal")
		}
		fmt.Println("---")
	}

	fmt.Printf("\n==== AVERAGE RESULTS AFTER %v RUNS====\n",n)
	fmt.Printf("Synchronous:                                           %v\n", totalSynchronous/time.Duration(n))
	fmt.Printf("Concurrent:                                            %v\n", totalConcurrent/time.Duration(n))
	
	// fmt.Printf("\n📄 Analyzed\n")
	// fmt.Println("IP Count:           ", result1.IPCount)
	// fmt.Println("Method Count:       ", result1.MethodCount)
	// fmt.Println("Status Code Count:  ", result1.StatusCount)
	// fmt.Println("Endpoint Count:     ", result1.EndpointCount)
	// fmt.Println("Size Bucket Count:  ", result1.SizeBucketCount)
	// fmt.Println("\nIP Count:           ", result2.IPCount)
	// fmt.Println("Method Count:       ", result2.MethodCount)
	// fmt.Println("Status Code Count:  ", result2.StatusCount)
	// fmt.Println("Endpoint Count:     ", result2.EndpointCount)
	// fmt.Println("Size Bucket Count:  ", result2.SizeBucketCount)
}

// generateLogFiles writes multiple log files with fake access logs
func GenerateLogFiles(
	numFiles int,
	linesPerFile int,
	numIPs int,
	numEndpoints int,
	methods []string,
	statusCodes []int,
	minSize int,
	maxSize int,
	outputDir string,
) error {
	os.MkdirAll(outputDir, os.ModePerm)

	for fileIndex := 1; fileIndex <= numFiles; fileIndex++ {
		filename := filepath.Join(outputDir, fmt.Sprintf("log%d.txt", fileIndex))
		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", filename, err)
		}
		defer f.Close()

		for i := 0; i < linesPerFile; i++ {
			ip := fmt.Sprintf("192.168.1.%d", rand.Intn(numIPs)+1)
			timestamp := time.Now().Add(-time.Duration(rand.Intn(30*24)) * time.Hour).Format("2006-01-02T15:04:05")
			method := methods[rand.Intn(len(methods))]
			endpoint := fmt.Sprintf("/api/v1/resource%d", rand.Intn(numEndpoints)+1)
			status := statusCodes[rand.Intn(len(statusCodes))]
			size := rand.Intn(maxSize-minSize+1) + minSize

			line := fmt.Sprintf("%s %s %s %s %d %d\n", ip, timestamp, method, endpoint, status, size)
			if _, err := f.WriteString(line); err != nil {
				return fmt.Errorf("failed to write line: %w", err)
			}
		}
	}

	return nil
}

func Analysis(filename string) AnalysisResult {
	result := NewAnalysisResult()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: ", err)
		return NewAnalysisResult()
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		ip := parts[0]
		method := parts[2]
		endpoint := parts[3]
		status := parts[4]
		size := parts[5]

		result.IPCount[ip]++
		result.MethodCount[method]++
		result.EndpointCount[endpoint]++
		result.StatusCount[status]++

		var sizeInt int
		fmt.Sscanf(size, "%d", &sizeInt)
		switch {
		case sizeInt < 500:
			result.SizeBucketCount["small"]++
		case sizeInt < 1000:
			result.SizeBucketCount["medium"]++
		default:
			result.SizeBucketCount["large"]++
		}
	}
	return result
}

func AnalysisSynchronous(numFiles int, outputDir string) AnalysisResult {
	result := NewAnalysisResult()
	for fileIndex := 1; fileIndex <= numFiles; fileIndex++ {
		filename := filepath.Join(outputDir, fmt.Sprintf("log%d.txt",fileIndex))
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error: ", err)
			return AnalysisResult{}
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Fields(line)
			ip := parts[0]
			method := parts[2]
			endpoint := parts[3]
			status := parts[4]
			size := parts[5]

			result.IPCount[ip]++
			result.MethodCount[method]++
			result.EndpointCount[endpoint]++
			result.StatusCount[status]++

			var sizeInt int
			fmt.Sscanf(size, "%d", &sizeInt)
			switch {
			case sizeInt < 500:
				result.SizeBucketCount["small"]++
			case sizeInt < 1000:
				result.SizeBucketCount["medium"]++
			default:
				result.SizeBucketCount["large"]++
			}
		}
	}
	return result
}

func AnalysisConcurrent(numFiles int, outputDir string) AnalysisResult {
	resultChan := make(chan AnalysisResult)
	var wg sync.WaitGroup

	for fileIndex := 1; fileIndex <= numFiles; fileIndex++ {
		filename := filepath.Join(outputDir, fmt.Sprintf("log%d.txt", fileIndex))
		wg.Add(1)
		go func() {
			defer wg.Done()
			temp := Analysis(filename)
			resultChan <- temp
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return MergeAnalysisResult(resultChan)
}

func MergeAnalysisResult(ch <-chan AnalysisResult) AnalysisResult {
	final := NewAnalysisResult()

	for analysis := range ch {
		for k, v := range analysis.IPCount {
			final.IPCount[k] += v
		}
		for k, v := range analysis.EndpointCount {
			final.EndpointCount[k] += v
		}
		for k, v := range analysis.MethodCount {
			final.MethodCount[k] += v
		}
		for k, v := range analysis.StatusCount {
			final.StatusCount[k] += v
		}
		for k, v := range analysis.SizeBucketCount {
			final.SizeBucketCount[k] += v
		}
	}

	return final
}