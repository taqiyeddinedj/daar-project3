package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/taqiyeddinedj/daar-project3/pkg/graph"
	"github.com/taqiyeddinedj/daar-project3/pkg/indexer"
	"github.com/taqiyeddinedj/daar-project3/pkg/ranking"
	"github.com/taqiyeddinedj/daar-project3/pkg/search"
	"github.com/taqiyeddinedj/daar-project3/pkg/storage"
)

type BenchmarkResult struct {
	QueryType   string    `json:"query_type"`
	Query       string    `json:"query"`
	Times       []float64 `json:"times_ms"`
	Mean        float64   `json:"mean_ms"`
	StdDev      float64   `json:"stddev_ms"`
	ResultCount int       `json:"result_count"`
}

type AllResults struct {
	SearchSimple    []BenchmarkResult `json:"search_simple"`
	SearchRegex     []BenchmarkResult `json:"search_regex"`
	Recommendations []BenchmarkResult `json:"recommendations"`
}

func main() {
	fmt.Println("=== DAAR Project 3 - Performance Benchmarks ===\n")
	fmt.Println("Loading index...")
	idx, err := storage.LoadFromFile("data/index.json")
	if err != nil {
		panic(err)
	}
	fmt.Printf(" Loaded %d books\n", len(idx.Books))
	fmt.Println("Loading Jaccard graph...")
	jaccardGraph, err := graph.LoadGraphFromFile("data/jaccard_graph.json")
	if err != nil {
		panic(err)
	}
	fmt.Printf(" Loaded graph with %d edges\n", jaccardGraph.EdgeCount)

	fmt.Println("Calculating PageRank...")
	pageRank := ranking.CalculatePageRank(jaccardGraph, 20, 0.85)
	fmt.Println(" PageRank calculated\n")

	results := AllResults{
		SearchSimple:    []BenchmarkResult{},
		SearchRegex:     []BenchmarkResult{},
		Recommendations: []BenchmarkResult{},
	}

	// ===== SIMPLE SEARCH BENCHMARKS =====
	fmt.Println("=== Testing Simple Search ===")

	simpleQueries := []string{
		"love",    // Very common word
		"king",    // Common word
		"whale",   // Medium frequency
		"zephyr",  // Rare word
		"quantum", // Very rare/absent
	}

	for _, query := range simpleQueries {
		result := benchmarkSimpleSearch(idx, pageRank, query, 100)
		results.SearchSimple = append(results.SearchSimple, result)
		fmt.Printf("  %s: %.2f ms (±%.2f), %d results\n",
			query, result.Mean, result.StdDev, result.ResultCount)
	}

	// ===== REGEX SEARCH BENCHMARKS =====
	fmt.Println("\n=== Testing Regex Search ===")

	regexQueries := []string{
		"wha.*",        // Simple wildcard
		"(king|queen)", // Alternation
		"[a-z]{10,}",   // Long words
		"^the.*end$",   // Start and end anchors
		".*love.*",     // Contains pattern
	}

	for _, query := range regexQueries {
		result := benchmarkRegexSearch(idx, pageRank, query, 50)
		results.SearchRegex = append(results.SearchRegex, result)
		fmt.Printf("  %s: %.2f ms (±%.2f), %d results\n",
			query, result.Mean, result.StdDev, result.ResultCount)
	}

	// ===== RECOMMENDATIONS BENCHMARKS =====
	fmt.Println("\n=== Testing Recommendations ===")

	// Sample 10 random books
	bookIDs := []int{}
	for id := range idx.Books {
		bookIDs = append(bookIDs, id)
		if len(bookIDs) >= 10 {
			break
		}
	}

	for _, bookID := range bookIDs {
		result := benchmarkRecommendations(jaccardGraph, bookID, 100)
		results.Recommendations = append(results.Recommendations, result)
	}

	avgTime := 0.0
	for _, r := range results.Recommendations {
		avgTime += r.Mean
	}
	avgTime /= float64(len(results.Recommendations))
	fmt.Printf("  Average: %.2f ms across %d books\n", avgTime, len(results.Recommendations))

	// Save results
	fmt.Println("\n=== Saving Results ===")
	file, _ := os.Create("data/benchmark_results.json")
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(results)
	file.Close()
	fmt.Println(" Results saved to benchmark_results.json")
}

func benchmarkSimpleSearch(idx *indexer.Indexer, pageRank map[int]float64, query string, iterations int) BenchmarkResult {
	times := []float64{}
	var resultCount int

	for i := 0; i < iterations; i++ {
		start := time.Now()
		results := search.Search(idx, query)
		results = ranking.RankResults(results, pageRank)
		elapsed := time.Since(start)

		times = append(times, float64(elapsed.Microseconds())/1000.0)
		resultCount = len(results)
	}

	mean, stddev := calculateStats(times)

	return BenchmarkResult{
		QueryType:   "simple",
		Query:       query,
		Times:       times,
		Mean:        mean,
		StdDev:      stddev,
		ResultCount: resultCount,
	}
}

func benchmarkRegexSearch(idx *indexer.Indexer, pageRank map[int]float64, query string, iterations int) BenchmarkResult {
	times := []float64{}
	var resultCount int

	for i := 0; i < iterations; i++ {
		start := time.Now()
		results, _ := search.RegexSearch(idx, query)
		results = ranking.RankResults(results, pageRank)
		elapsed := time.Since(start)

		times = append(times, float64(elapsed.Microseconds())/1000.0)
		resultCount = len(results)
	}

	mean, stddev := calculateStats(times)

	return BenchmarkResult{
		QueryType:   "regex",
		Query:       query,
		Times:       times,
		Mean:        mean,
		StdDev:      stddev,
		ResultCount: resultCount,
	}
}

func benchmarkRecommendations(jaccardGraph *graph.JaccardGraph, bookID int, iterations int) BenchmarkResult {
	times := []float64{}
	var resultCount int

	for i := 0; i < iterations; i++ {
		start := time.Now()
		recommendations := graph.GetRecommendations(jaccardGraph, bookID, 10)
		elapsed := time.Since(start)

		times = append(times, float64(elapsed.Microseconds())/1000.0)
		resultCount = len(recommendations)
	}

	mean, stddev := calculateStats(times)

	return BenchmarkResult{
		QueryType:   "recommendation",
		Query:       fmt.Sprintf("book_%d", bookID),
		Times:       times,
		Mean:        mean,
		StdDev:      stddev,
		ResultCount: resultCount,
	}
}

func calculateStats(values []float64) (mean, stddev float64) {
	if len(values) == 0 {
		return 0, 0
	}

	// Calculate mean
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean = sum / float64(len(values))

	// Calculate standard deviation
	sumSquaredDiff := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}
	variance := sumSquaredDiff / float64(len(values))
	stddev = sqrt(variance)

	return mean, stddev
}

func sqrt(x float64) float64 {
	if x < 0 {
		return 0
	}
	z := 1.0
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}
