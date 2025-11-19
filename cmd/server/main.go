package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taqiyeddinedj/daar-project3/pkg/graph"
	"github.com/taqiyeddinedj/daar-project3/pkg/indexer"
	"github.com/taqiyeddinedj/daar-project3/pkg/ranking"
	"github.com/taqiyeddinedj/daar-project3/pkg/search"
	"github.com/taqiyeddinedj/daar-project3/pkg/storage"
)

var (
	idx          *indexer.Indexer
	jaccardGraph *graph.JaccardGraph
	pageRank     map[int]float64
)

func main() {
	fmt.Println("=== Starting Search Engine Server ===\n")

	// Load data
	fmt.Println("Loading index...")
	var err error
	idx, err = storage.LoadFromFile("data/index.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" Index: %d books\n", len(idx.Books))

	fmt.Println("Loading Jaccard graph...")
	jaccardGraph, err = graph.LoadGraphFromFile("data/jaccard_graph.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" Graph: %d edges\n", jaccardGraph.EdgeCount)

	fmt.Println("Calculating PageRank...")
	pageRank = ranking.CalculatePageRank(jaccardGraph, 20, 0.85)
	fmt.Println(" PageRank calculated")

	fmt.Println("\n Server ready! Starting API on :8080\n")

	// Setup routes
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	// Serve static files
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	// Routes
	r.GET("/", homeHandler)
	r.GET("/api/search", searchHandler)
	r.GET("/api/search/regex", regexSearchHandler)
	r.GET("/api/book/:id", bookDetailHandler)
	r.GET("/api/recommend/:id", recommendHandler)

	log.Fatal(r.Run(":8080"))
}

func homeHandler(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func searchHandler(c *gin.Context) {
	keyword := c.Query("q")

	results := search.Search(idx, keyword)
	results = ranking.RankResults(results, pageRank)

	c.JSON(200, gin.H{"results": results, "count": len(results)})
}

func regexSearchHandler(c *gin.Context) {
	pattern := c.Query("pattern")

	results, err := search.RegexSearch(idx, pattern)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	results = ranking.RankResults(results, pageRank)
	c.JSON(200, gin.H{"results": results, "count": len(results)})
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

	recommendations := graph.GetRecommendations(jaccardGraph, id, 5)

	results := []map[string]interface{}{}
	for _, edge := range recommendations {
		book := idx.Books[edge.Target]
		results = append(results, map[string]interface{}{
			"book":       book,
			"similarity": edge.Similarity,
		})
	}

	c.JSON(200, results)
}
