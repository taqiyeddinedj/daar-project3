package main

import (
	"fmt"
	"os"

	"github.com/taqiyeddinedj/daar-project3/pkg/search"
	"github.com/taqiyeddinedj/daar-project3/pkg/storage"
)

func main() {
	fmt.Println("=== Testing Search Engine ===")

	// Load index
	fmt.Println("Loading index...")
	idx, err := storage.LoadFromFile("data/index.json")
	if err != nil {
		fmt.Printf("Error loading index: %v\n", err)
		fmt.Println("Run 'go run cmd/build_index/main.go' first!")
		os.Exit(1)
	}

	fmt.Printf(" Index loaded: %d books, %d unique words\n\n",
		len(idx.Books), idx.UniqueWords)

	// Test searches
	testQueries := []string{"whale", "captain", "treasure", "love", "dragon"}

	for _, query := range testQueries {
		fmt.Printf("Searching for '%s'...\n", query)
		results := search.Search(idx, query)

		if len(results) == 0 {
			fmt.Printf("  No books found\n\n")
			continue
		}

		fmt.Printf("  Found in %d books:\n", len(results))

		// Show top 5
		max := 5
		if len(results) < max {
			max = len(results)
		}

		for i := 0; i < max; i++ {
			r := results[i]
			fmt.Printf("    %d. [Book %d] %s by %s - %d occurrences\n",
				i+1, r.Book.ID, r.Book.Title, r.Book.Author, r.Occurrences)
		}

		if len(results) > 5 {
			fmt.Printf("    ... and %d more books\n", len(results)-5)
		}
		fmt.Println()
	}
}
