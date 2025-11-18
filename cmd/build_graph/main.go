package main

import (
	"fmt"
	"os"
	"time"

	"github.com/taqiyeddinedj/daar-project3/pkg/graph"
	"github.com/taqiyeddinedj/daar-project3/pkg/storage"
)

func main() {
	fmt.Println("=== Building Jaccard Graph ===\n")

	// Load index
	fmt.Println("Loading index...")
	idx, err := storage.LoadFromFile("data/index.json")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf(" Loaded %d books\n\n", len(idx.Books))

	// Build graph
	startTime := time.Now()
	threshold := 0.1
	jaccardGraph := graph.BuildJaccardGraph(idx, threshold)
	elapsed := time.Since(startTime)

	fmt.Printf("\n Graph built in %v\n", elapsed)
	fmt.Println("\nSaving graph to disk...")
	err = jaccardGraph.SaveToFile("data/jaccard_graph.json")
	if err != nil {
		fmt.Printf("Error saving: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n=== Graph Statistics ===")
	stats := jaccardGraph.GetStats()
	for key, value := range stats {
		fmt.Printf("%s: %v\n", key, value)
	}
}
