package ranking

import (
	"sort"

	"github.com/taqiyeddinedj/daar-project3/pkg/graph"
	"github.com/taqiyeddinedj/daar-project3/pkg/models"
)

// CalculatePageRank computes PageRank scores for all books
func CalculatePageRank(jaccardGraph *graph.JaccardGraph, iterations int, dampingFactor float64) map[int]float64 {
	bookIDs := make([]int, 0, len(jaccardGraph.Edges))
	for bookID := range jaccardGraph.Edges {
		bookIDs = append(bookIDs, bookID)
	}

	n := len(bookIDs)
	if n == 0 {
		return map[int]float64{}
	}

	// Initialize PageRank scores (equal distribution)
	pageRank := make(map[int]float64)
	for _, id := range bookIDs {
		pageRank[id] = 1.0 / float64(n)
	}

	for iter := 0; iter < iterations; iter++ {
		newPageRank := make(map[int]float64)

		for _, bookID := range bookIDs {
			rank := (1.0 - dampingFactor) / float64(n)

			// Add contributions from neighbors
			for _, edge := range jaccardGraph.Edges[bookID] {
				neighborID := edge.Target
				numOutLinks := len(jaccardGraph.Edges[neighborID])

				if numOutLinks > 0 {
					rank += dampingFactor * pageRank[neighborID] / float64(numOutLinks)
				}
			}

			newPageRank[bookID] = rank
		}

		pageRank = newPageRank
	}

	return pageRank
}

// RankResults sorts search results by PageRank
func RankResults(results []models.SearchResult, pageRank map[int]float64) []models.SearchResult {
	// Update relevance with PageRank
	for i := range results {
		bookID := results[i].Book.ID
		pr := pageRank[bookID]

		// Combine occurrence count with PageRank
		results[i].Relevance = float64(results[i].Occurrences) * (1.0 + pr*10)
	}

	// Sort by new relevance
	sort.Slice(results, func(i, j int) bool {
		return results[i].Relevance > results[j].Relevance
	})

	return results
}
