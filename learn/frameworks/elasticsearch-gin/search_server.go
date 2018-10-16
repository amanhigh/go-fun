package elasticsearch_gin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	. "github.com/amanhigh/go-fun/models/learn/frameworks"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"github.com/teris-io/shortid"
)

const (
	elasticIndexName = "learn"
	elasticTypeName  = "document"
	endpoint         = "http://192.168.99.100:9200"
)

var (
	elasticClient *elastic.Client
)

/*
	Add document- curl -X POST http://localhost:8080/documents -d '[{ "Title":"Aman", "Content": "Preet"}]'
	Query Elastic Search - http://192.168.99.100:9200/learn/_search?pretty
	Query API - http://localhost:8080/search?query=aman&skip=1&top=2
*/

func StartSearchServer() {
	var err error
	// Create Elastic client and wait for Elasticsearch to be ready
	for {
		elasticClient, err = elastic.NewClient(
			elastic.SetURL(endpoint),
			elastic.SetSniff(false),
		)
		if err != nil {
			log.Println(err)
			// Retry every 3 seconds
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	// Start HTTP server
	r := gin.Default()
	r.POST("/documents", createDocumentsEndpoint)
	r.GET("/search", searchEndpoint)
	if err = r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func createDocumentsEndpoint(c *gin.Context) {
	// Parse request
	var docs []DocumentRequest
	if err := c.BindJSON(&docs); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		return
	}
	// Insert documents in bulk
	bulk := elasticClient.Bulk().Index(elasticIndexName).Type(elasticTypeName)
	for _, d := range docs {
		doc := Document{
			ID:        shortid.MustGenerate(),
			Title:     d.Title,
			CreatedAt: time.Now().UTC(),
			Content:   d.Content,
		}
		bulk.Add(elastic.NewBulkIndexRequest().Id(doc.ID).Doc(doc))
	}
	if _, err := bulk.Do(c.Request.Context()); err != nil {
		log.Println(err)
		errorResponse(c, http.StatusInternalServerError, "Failed to create documents")
		return
	}
	c.Status(http.StatusOK)
}

func searchEndpoint(c *gin.Context) {
	// Parse request
	query := c.Query("query")
	if query == "" {
		errorResponse(c, http.StatusBadRequest, "Query not specified")
		return
	}
	skip := 0
	top := 10
	if i, err := strconv.Atoi(c.Query("skip")); err == nil {
		skip = i
	}
	if i, err := strconv.Atoi(c.Query("top")); err == nil {
		top = i
	}

	// Perform search
	esQuery := elastic.NewMultiMatchQuery(query, "title", "content").
		Fuzziness("2").
		MinimumShouldMatch("2")

	result, err := elasticClient.Search().
		Index(elasticIndexName).
		Query(esQuery).
		From(skip).Size(top).
		Do(c.Request.Context())

	if err == nil {
		c.JSON(http.StatusOK, resultToResponse(result))
	} else {
		log.Println(err)
		errorResponse(c, http.StatusInternalServerError, "Something went wrong")
	}
}

func resultToResponse(result *elastic.SearchResult) SearchResponse {
	res := SearchResponse{
		Time: fmt.Sprintf("%d", result.TookInMillis),
		Hits: fmt.Sprintf("%d", result.Hits.TotalHits),
	}
	// Transform search results before returning them
	docs := make([]DocumentResponse, 0)
	for _, hit := range result.Hits.Hits {
		var doc DocumentResponse
		if err := json.Unmarshal(*hit.Source, &doc); err == nil {
			docs = append(docs, doc)
		} else {
			log.Println(err)
		}
	}
	res.Documents = docs
	return res
}

func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"error": err,
	})
}
