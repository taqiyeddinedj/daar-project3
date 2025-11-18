package graph

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

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
func CalculateSimilarity(wordsA, wordsB []string) float64 {
	setA := make(map[string]bool)
	for _, word := range wordsA {
		setA[word] = true
	}
	setB := make(map[string]bool)
	for _, word := range wordsB {
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

type Edge struct {
	Source     int     `json:"source"`
	Target     int     `json:"target"`
	Similarity float64 `json:"similarity"`
}

type JaccardGraph struct {
	Edges     map[int][]Edge `json:"edges"`
	Threshold float64        `json:"threshold"`
	BookCount int            `json:"book_count"`
	EdgeCount int            `json:"edge_count"`
}

// ------------------------------------------------------
// Need a graph structure
// I think the source key:value can be removed
/// But for the sake of code visibility i am keeping it
/* graph.Edges[2701] = [
    {Source: 2701, Target: 164, Similarity: 0.45},
    {Source: 2701, Target: 46, Similarity: 0.32}
					]
*/
// ------------------------------------------------------

func BuildJaccardGraph(idx *indexer.Indexer, threshold float64) *JaccardGraph {
	fmt.Println("Building Jaccard graph...")
	fmt.Printf("Threshold: %.3f\n", threshold)

	bookToWords := BuildBookToWords(idx)

	// Initialize graph with metadata
	graph := &JaccardGraph{
		Edges:     make(map[int][]Edge),
		Threshold: threshold,
		BookCount: len(idx.Books),
	}

	//	Get all book IDs
	//	Another way is to loop over idx.books twice,
	// 	and add an if to continue on the first sourcebook == targetbook
	bookIDs := make([]int, 0, len(idx.Books))
	for id := range idx.Books {
		bookIDs = append(bookIDs, id)
	}

	fmt.Printf("Comparing %d books...\n", len(bookIDs))

	for i := 0; i < len(bookIDs); i++ {
		bookA := bookIDs[i]
		wordsA := bookToWords[bookA]

		for j := i + 1; j < len(bookIDs); j++ {
			bookB := bookIDs[j]
			wordsB := bookToWords[bookB]

			similarity := CalculateSimilarity(wordsA, wordsB)

			if similarity > threshold {
				// Add edge A → B
				graph.Edges[bookA] = append(graph.Edges[bookA], Edge{
					Source:     bookA,
					Target:     bookB,
					Similarity: similarity,
				})

				// Add edge B → A
				graph.Edges[bookB] = append(graph.Edges[bookB], Edge{
					Source:     bookB,
					Target:     bookA,
					Similarity: similarity,
				})
			}
		}

		if (i+1)%100 == 0 {
			fmt.Printf("  Processed %d/%d books\n", i+1, len(bookIDs))
		}
	}

	graph.EdgeCount = CountEdges(graph)

	fmt.Printf("\n Graph complete!\n")
	fmt.Printf("  Total edges: %d\n", graph.EdgeCount)
	fmt.Printf("  Books with connections: %d\n", len(graph.Edges))

	return graph
}

// this functions is going to be moved out in another package

func GetRecommendations(graph *JaccardGraph, bookID int, topN int) []Edge {
	// Get all similar books
	edges := graph.Edges[bookID]

	// Sort by similarity (highest first)
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Similarity > edges[j].Similarity
	})

	// Return top N
	if len(edges) > topN {
		return edges[:topN]
	}
	return edges
}

// SaveToFile saves the graph to a JSON file
func (g *JaccardGraph) SaveToFile(filename string) error {
	// Calculate edge count before saving
	g.EdgeCount = CountEdges(g)

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print

	err = encoder.Encode(g)
	if err != nil {
		return fmt.Errorf("failed to encode graph: %w", err)
	}

	fmt.Printf("Graph saved to %s\n", filename)
	fmt.Printf("  Edges: %d\n", g.EdgeCount)
	fmt.Printf("  Books with connections: %d\n", len(g.Edges))

	return nil
}

func LoadGraphFromFile(filename string) (*JaccardGraph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var graph JaccardGraph
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&graph)
	if err != nil {
		return nil, fmt.Errorf("failed to decode graph: %w", err)
	}

	fmt.Printf("Graph loaded from %s\n", filename)
	fmt.Printf("  Edges: %d\n", graph.EdgeCount)
	fmt.Printf("  Threshold: %.3f\n", graph.Threshold)

	return &graph, nil
}

func CountEdges(graph *JaccardGraph) int {
	total := 0
	for _, edges := range graph.Edges {
		total += len(edges)
	}
	// Divide by 2 because graph is undirected (each edge stored twice)
	return total / 2
}

func (g *JaccardGraph) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"total_edges":              g.EdgeCount,
		"books_with_connections":   len(g.Edges),
		"threshold":                g.Threshold,
		"avg_connections_per_book": float64(g.EdgeCount*2) / float64(len(g.Edges)),
	}
}
