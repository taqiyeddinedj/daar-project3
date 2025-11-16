package search

import (
	"strings"

	"github.com/taqiyeddinedj/daar-project3/pkg/indexer"
	"github.com/taqiyeddinedj/daar-project3/pkg/models"
)

// Search finds books containing a keyword
func Search(idx *indexer.Indexer, keyword string) []models.SearchResult {
	keyword = strings.ToLower(keyword)

	bookOccurrences, found := idx.WordToBooks[keyword]
	if !found {
		return []models.SearchResult{}
	}

	results := make([]models.SearchResult, 0, len(bookOccurrences))

	for bookID, count := range bookOccurrences {
		book, exists := idx.Books[bookID]
		if !exists {
			continue
		}

		results = append(results, models.SearchResult{
			Book:        book,
			Occurrences: count,
			Relevance:   float64(count),
		})
	}

	return results
}
