package graph

import (
	"fmt"

	"github.com/taqiyeddinedj/daar-project3/pkg/indexer"
)

// Build book to words,
// each book wil hold an arrray of words.
func BuildBookToWords(idx *indexer.Indexer) map[int][]string {
	booksToWords := make(map[int][]string)

	for word, books := range idx.WordToBooks {
		for bookID := range books {
			booksToWords[bookID] = append(booksToWords[bookID], word)
		}
	}
	return booksToWords
}

// Return the corresponding intersection devided by union.
// The function is going to be used by BuildBookToWord().
func CalculateSimularity(wordsA, wordsB []string) float64 {
	setA := make(map[string]bool)
	for _, word := range wordsA {
		setA[word] = true
	}
	setB := make(map[string]bool)
	for _, word := range wordsA {
		setB[word] = true
	}

	intersection := 0
	for word := range setA {
		if setB[word] {
			intersection++
		}
	}

	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 0.0 // Avoid division by zero
	}

	return float64(intersection) / float64(union)
}

// ------------------------------------------------------
// Need a grap structure
// I think the source key:value can be removed
/// But for the sake of code visibility i am keeping it
/* graph.Edges[2701] = [
    {Source: 2701, Target: 164, Similarity: 0.45},
    {Source: 2701, Target: 46, Similarity: 0.32}
					]
*/
// ------------------------------------------------------

type Edge struct {
	Source     int
	Target     int
	Similarity float64
}

type JackardGraph struct {
	Edge map[int][]Edge
}

func BuildJackardGraph(idx indexer.Indexer, threshold float64) *JackardGraph {
	fmt.Println("Building Jaccard graph...")
	bookToWords := BuildBookToWords(&idx)

}
