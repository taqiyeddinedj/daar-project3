package main

import (
	"fmt"
	"os"
	"time"

	"github.com/taqiyeddinedj/daar-project3/pkg/indexer"
	"github.com/taqiyeddinedj/daar-project3/pkg/storage"
)

func main() {
	fmt.Println("=== Building Search Index ===")
	fmt.Println("This may take 30-60 minutes...")
	fmt.Println()

	startTime := time.Now()

	// Create indexer
	idx := indexer.NewIndexer()

	// Build index
	fmt.Println("Step 1: Reading and indexing all books...")
	err := idx.BuildIndexFromDirectory("data/books")
	if err != nil {
		fmt.Printf("Error building index: %v\n", err)
		os.Exit(1)
	}

	// Print stats
	fmt.Println("\n=== Index Statistics ===")
	stats := idx.GetStats()
	for key, value := range stats {
		fmt.Printf("%s: %v\n", key, value)
	}

	// Save index
	fmt.Println("\nStep 2: Saving index to disk...")
	indexPath := "data/index.json"
	err = storage.SaveToFile(idx, indexPath)
	if err != nil {
		fmt.Printf("Error saving index: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\n Index built successfully in %v\n", elapsed)
	fmt.Printf("Index saved to: %s\n", indexPath)
}
