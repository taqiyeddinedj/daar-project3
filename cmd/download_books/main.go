package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	minWords     = 10000 // requirement from the project
	targetBooks  = 1664  // minimum number of books
	baseDir      = "data/books"
	startBookID  = 1
	maxBookIDTry = 50000 // safety limit so we do not loop forever
)

func downloadBook(id int) error {
	url := fmt.Sprintf("https://www.gutenberg.org/cache/epub/%d/pg%d.txt", id, id)
	filepath := fmt.Sprintf("data/books/book_%d.txt", id)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("book %d not found", id)
	}

	// Check for book size, if word count > 1000
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body for ID %d: %w", id, err)
	}
	text := string(body)
	words := strings.Fields(text)
	wordCount := len(words)
	if wordCount < minWords {
		fmt.Printf("ID %d skipped: only %d words (need at least %d).\n", id, wordCount, minWords)
		return nil
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.Write(body)
	return err
}

func main() {
	savedcount := 0
	for id := startBookID; id <= maxBookIDTry && savedcount < targetBooks; id++ {
		// Check if the file already exist
		filepath := fmt.Sprintf("data/books/book_%d.txt", id)
		if _, err := os.Stat(filepath); err == nil {
			fmt.Printf("Book %d already exists, skipping.\n", id)
			savedcount++
			continue
		}

		fmt.Printf("Trying book ID %d (saved: %d/%d)\n", id, savedcount, targetBooks)
		err := downloadBook(id)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		if _, err := os.Stat(filepath); err == nil {
			savedcount++
		}
	}
	fmt.Printf("Finished. Total saved books: %d\n", savedcount)
}
