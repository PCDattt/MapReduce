package main

import (
	"fmt"
)

func main() {
	// calculateAverage(5, 500, 600_000)
	// Configurable parameters
	numFiles := 10
	linesPerFile := 500
	numIPs := 5
	numEndpoints := 5
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	statusCodes := []int{100, 200, 300, 400}
	minSize, maxSize := 200, 1600
	outputDir := "logs"

	err := generateLogFiles(numFiles, linesPerFile, numIPs, numEndpoints, methods, statusCodes, minSize, maxSize, outputDir)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ Generated %d log files in ./%s/\n", numFiles, outputDir)
}
