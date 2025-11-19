package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taqiyeddinedj/daar-project3/pkg/graph"
	"github.com/taqiyeddinedj/daar-project3/pkg/indexer"
	"github.com/taqiyeddinedj/daar-project3/pkg/models"
	"github.com/taqiyeddinedj/daar-project3/pkg/ranking"
	"github.com/taqiyeddinedj/daar-project3/pkg/search"
	"github.com/taqiyeddinedj/daar-project3/pkg/storage"
)

var (
	idx          *indexer.Indexer
	jaccardGraph *graph.JaccardGraph
	pageRank     map[int]float64
)

type SearchResponse struct {
	Books      []models.Book         `json:"books"`
	Results    []models.SearchResult `json:"results"`
	TotalCount int                   `json:"total_count"`
	Page       int                   `json:"page"`
	PerPage    int                   `json:"per_page"`
	TotalPages int                   `json:"total_pages"`
}

func main() {
	fmt.Println("=== Starting Search Engine Server ===\n")

	fmt.Println("Loading index...")
	var err error
	idx, err = storage.LoadFromFile("data/index.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Index: %d books\n", len(idx.Books))

	fmt.Println("Loading Jaccard graph...")
	jaccardGraph, err = graph.LoadGraphFromFile("data/jaccard_graph.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Graph: %d edges\n", jaccardGraph.EdgeCount)

	fmt.Println("Calculating PageRank...")
	pageRank = ranking.CalculatePageRank(jaccardGraph, 20, 0.85)
	fmt.Println("✓ PageRank calculated")

	fmt.Println("\n🚀 Server: http://localhost:8080\n")

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	})

	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	r.GET("/", homeHandler)
	r.GET("/api/search", searchHandler)
	r.GET("/api/book/:id", bookDetailHandler)
	r.GET("/api/recommendations/:id", recommendHandler)
	r.GET("/api/content/:id", contentHandler)

	log.Fatal(r.Run(":8080"))
}

func homeHandler(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func searchHandler(c *gin.Context) {
	query := c.Query("q")
	searchType := c.Query("type")
	pageStr := c.DefaultQuery("page", "1")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	perPage := 20

	var results []models.SearchResult
	if searchType == "regex" {
		var err error
		results, err = search.RegexSearch(idx, query)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	} else {
		results = search.Search(idx, query)
	}

	results = ranking.RankResults(results, pageRank)

	// Extract books for pagination
	var books []models.Book
	for _, r := range results {
		books = append(books, r.Book)
	}

	totalCount := len(books)
	totalPages := (totalCount + perPage - 1) / perPage
	start := (page - 1) * perPage
	end := start + perPage

	var paginatedBooks []models.Book

	if start >= totalCount {
		paginatedBooks = []models.Book{}
	} else {
		if end > totalCount {
			end = totalCount
		}
		paginatedBooks = books[start:end]
	}

	response := SearchResponse{
		Books:      paginatedBooks,
		Results:    results,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}

	c.JSON(200, response)
}

func bookDetailHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	book, exists := idx.Books[id]
	if !exists {
		c.JSON(404, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(200, book)
}

func recommendHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	recommendations := graph.GetRecommendations(jaccardGraph, id, 10)

	books := []models.Book{}
	for _, edge := range recommendations {
		book := idx.Books[edge.Target]
		books = append(books, book)
	}

	c.JSON(200, books)
}

func contentHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	book, exists := idx.Books[id]
	if !exists {
		c.JSON(404, gin.H{"error": "Book not found"})
		return
	}

	possiblePaths := []string{
		book.FilePath,
		filepath.Join("data/books", book.FilePath),
		filepath.Join("books", book.FilePath),
	}

	var content []byte
	var err error
	var successPath string

	for _, path := range possiblePaths {
		content, err = os.ReadFile(path)
		if err == nil {
			successPath = path
			break
		}
	}

	if err != nil {
		log.Printf("Failed to read book %d. Tried paths: %v", id, possiblePaths)
		c.JSON(500, gin.H{
			"error":       "Failed to read book content",
			"details":     fmt.Sprintf("Book ID: %d, FilePath: %s", id, book.FilePath),
			"tried_paths": possiblePaths,
		})
		return
	}

	log.Printf("✓ Successfully read book %d from: %s", id, successPath)

	c.JSON(200, gin.H{
		"book_id": book.ID,
		"title":   book.Title,
		"author":  book.Author,
		"content": string(content),
	})
}
