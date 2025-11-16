package indexer

import (
	"strings"
	"unicode"
)

// Stop words to skip during indexing
var StopWords = map[string]bool{
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

// Tokenize converts text into cleaned words
func Tokenize(text string) []string {
	var words []string
	var current strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			current.WriteRune(unicode.ToLower(r))
		} else if current.Len() > 0 {
			word := current.String()
			if len(word) > 2 && !StopWords[word] {
				words = append(words, word)
			}
			current.Reset()
		}
	}

	if current.Len() > 0 {
		word := current.String()
		if len(word) > 2 && !StopWords[word] {
			words = append(words, word)
		}
	}

	return words
}
