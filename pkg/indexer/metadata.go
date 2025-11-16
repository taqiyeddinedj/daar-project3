package indexer

import (
	"bufio"
	"os"
	"strings"
)

// ExtractTitle tries to extract book title from first 30 lines
func ExtractTitle(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		return "Unknown"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() && lineCount < 30 {
		line := scanner.Text()
		lineCount++

		if strings.HasPrefix(line, "Title:") {
			title := strings.TrimSpace(strings.TrimPrefix(line, "Title:"))
			if title != "" {
				return title
			}
		}
	}

	return "Unknown"
}

// ExtractAuthor tries to extract author from first 30 lines
func ExtractAuthor(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		return "Unknown"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() && lineCount < 30 {
		line := scanner.Text()
		lineCount++

		if strings.HasPrefix(line, "Author:") {
			author := strings.TrimSpace(strings.TrimPrefix(line, "Author:"))
			if author != "" {
				return author
			}
		}
	}

	return "Unknown"
}
