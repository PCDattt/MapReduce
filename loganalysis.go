package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// generateLogFiles writes multiple log files with fake access logs
func generateLogFiles(
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
