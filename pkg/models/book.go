package models

type Book struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	FilePath  string `json:"file_path"`
	WordCount int    `json:"word_count"`
}

type SearchResult struct {
	Book        Book    `json:"book"`
	Occurrences int     `json:"occurrences"`
	Relevance   float64 `json:"relevance"`
}
