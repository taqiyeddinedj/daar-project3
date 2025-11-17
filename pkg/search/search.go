package search

import (
	"regexp"
	"sort"
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
	// sorting the book was missing, we got an output of random mactched book which is not efficient!!
	sort.Slice(results, func(i, j int) bool {
		return results[i].Occurrences > results[j].Occurrences
	})

	return results
}

func RegexSearch(idx *indexer.Indexer, pattern string) ([]models.SearchResult, error) {
	// we need to have an engine that treats the regex ??
	// wha* ==> ? how to guess it to whale
	// run egrep on the index.json, capture the output then reutrn the the bookoccurences with that word
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	var matchingWords []string
	for word := range idx.WordToBooks {
		if re.MatchString(word) {
			matchingWords = append(matchingWords, word)
		}
	}

	bookOccurrences := make(map[int]int)
	for _, word := range matchingWords {
		for bookID, count := range idx.WordToBooks[word] {
			bookOccurrences[bookID] += count
		}
	}

	results := []models.SearchResult{}
	for bookID, totalCount := range bookOccurrences {
		book := idx.Books[bookID]
		results = append(results, models.SearchResult{
			Book:        book,
			Occurrences: totalCount,
			Relevance:   float64(totalCount),
		})
	}
	// sorting the book was missing, we got an output of random mactched book which is not efficient!!
	sort.Slice(results, func(i, j int) bool {
		return results[i].Occurrences > results[j].Occurrences
	})
	return results, nil
}
