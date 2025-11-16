package indexer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/taqiyeddinedj/daar-project3/pkg/models"
)

type Indexer struct {
	WordToBooks map[string]map[int]int `json:"word_to_books"`
	Books       map[int]models.Book    `json:"books"`
	TotalWords  int                    `json:"total_words"`
	UniqueWords int                    `json:"unique_words"`
}

// NewIndexer creates a new empty indexer
func NewIndexer() *Indexer {
	return &Indexer{
		WordToBooks: make(map[string]map[int]int),
		Books:       make(map[int]models.Book),
	}
}

// IndexBook reads a book file and adds it to the index
func (idx *Indexer) IndexBook(bookID int, filepath string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read book %d: %w", bookID, err)
	}

	words := Tokenize(string(content))
	if len(words) == 0 {
		return fmt.Errorf("book %d has no valid words", bookID)
	}

	title := ExtractTitle(filepath)
	author := ExtractAuthor(filepath)

	idx.Books[bookID] = models.Book{
		ID:        bookID,
		Title:     title,
		Author:    author,
		FilePath:  filepath,
		WordCount: len(words),
	}

	wordCount := make(map[string]int)
	for _, word := range words {
		wordCount[word]++
	}

	for word, count := range wordCount {
		if idx.WordToBooks[word] == nil {
			idx.WordToBooks[word] = make(map[int]int)
		}
		idx.WordToBooks[word][bookID] = count
	}

	idx.TotalWords += len(words)
	return nil
}

// BuildIndexFromDirectory scans all books in a directory
func (idx *Indexer) BuildIndexFromDirectory(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "book_*.txt"))
	if err != nil {
		return err
	}

	fmt.Printf("Found %d book files to index...\n", len(files))

	for i, file := range files {
		var bookID int
		_, err := fmt.Sscanf(filepath.Base(file), "book_%d.txt", &bookID)
		if err != nil {
			fmt.Printf("Skipping invalid filename: %s\n", file)
			continue
		}

		if (i+1)%100 == 0 || i == 0 {
			fmt.Printf("Indexing book %d/%d (ID: %d)...\n", i+1, len(files), bookID)
		}

		err = idx.IndexBook(bookID, file)
		if err != nil {
			fmt.Printf("Error indexing book %d: %v\n", bookID, err)
			continue
		}
	}

	idx.UniqueWords = len(idx.WordToBooks)

	fmt.Printf("\n✓ Indexing complete!\n")
	fmt.Printf("  Total books: %d\n", len(idx.Books))
	fmt.Printf("  Total words: %d\n", idx.TotalWords)
	fmt.Printf("  Unique words: %d\n", idx.UniqueWords)

	return nil
}

// GetStats returns index statistics
func (idx *Indexer) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"total_books":   len(idx.Books),
		"total_words":   idx.TotalWords,
		"unique_words":  idx.UniqueWords,
		"index_size_mb": float64(idx.TotalWords*8) / (1024 * 1024),
	}
}
