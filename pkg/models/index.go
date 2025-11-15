package models

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// List of words to avoid indexing later on
var stopWords = map[string]bool{
	"the": true, "a": true, "an": true, "and": true, "or": true,
	"but": true, "in": true, "on": true, "at": true, "to": true,
	"for": true, "of": true, "as": true, "by": true, "is": true,
	"was": true, "are": true, "were": true, "been": true, "be": true,
	"have": true, "has": true, "had": true, "do": true, "does": true,
	"did": true, "will": true, "would": true, "could": true, "should": true,
	"this": true, "that": true, "these": true, "those": true, "it": true,
	"he": true, "she": true, "they": true, "we": true, "you": true,
	"i": true, "me": true, "my": true, "mine": true, "your": true,
	"his": true, "her": true, "their": true, "its": true, "our": true,
}

type Indexer struct {
	WordToBook map[string]map[int]int
	Books      map[int]Book

	TotalWords  int
	UniqueWords int
}

func NewIndexer() *Indexer {
	return &Indexer{
		WordToBook: make(map[string]map[int]int),
		Books:      make(map[int]Book),
	}
}
func Tokenizer(text string) []string {
	var words []string
	var current strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			current.WriteRune(unicode.ToLower(r))
		} else if current.Len() > 0 {
			word := current.String()
			if len(word) > 2 && !stopWords[word] {
				words = append(words, word)
			}
		}
	}
	if current.Len() > 0 {
		word := current.String()
		if len(word) > 2 && !stopWords[word] {
			words = append(words, word)
		}
	}

	return words
}

func ExtractTitle(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		return "couldn't open file, or book file dosen't exists"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	linecount := 0
	for scanner.Scan() && linecount <= 30 {
		line := scanner.Text()
		linecount++

		if strings.HasPrefix(line, "Title") {
			title := strings.TrimSpace(strings.TrimPrefix(line, "Title:"))
			if title != "" {
				return title
			}
		}
	}
	return "uknown"
}

func ExtractAuthor(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("File couln't be opened or book dosent exists")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	linecount := 0

	for scanner.Scan() && linecount < 30 {
		line := scanner.Text()
		linecount++
		if strings.HasPrefix(line, "Author") {
			author := strings.TrimSpace(strings.TrimPrefix(line, "Author:"))
			if author != "" {
				return author
			}
		}
	}
	return "unkown"
}

func (idx *Indexer) IndexBook(bookID int, filepath string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read book %d: %w", bookID, err)
	}
	words := Tokenizer(string(content))
	if len(words) == 0 {
		return fmt.Errorf("book %d has no valid words", bookID)
	}
	title := ExtractTitle(filepath)
	author := ExtractAuthor(filepath)

	idx.Books[bookID] = Book{
		ID:        bookID,
		Title:     title,
		Author:    author,
		FilePath:  filepath,
		WordCount: len(words),
	}

	// WordToBook map[string]map[int]int:
	/*
		"whale" : {
			"blue_whale_book" : {
				500
			}
		}
	*/
	wordCount := make(map[string]int)
	for _, word := range words {
		wordCount[word]++
	}
	for word, count := range wordCount {
		if idx.WordToBook[word] == nil {
			idx.WordToBook[word] = make(map[int]int)
		}
		idx.WordToBook[word][bookID] = count
	}
	idx.TotalWords += len(words)
	return nil
}
