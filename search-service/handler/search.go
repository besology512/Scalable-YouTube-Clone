package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"search-service/elastic"
	"search-service/model"

	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query cannot be empty"})
		return
	}

	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title", "description", "tags"},
			},
		},
	}

	body, _ := json.Marshal(searchBody)
	res, err := elastic.ES.Search(
		elastic.ES.Search.WithContext(context.Background()),
		elastic.ES.Search.WithIndex("videos"),
		elastic.ES.Search.WithBody(io.NopCloser(bytes.NewReader(body))),
		elastic.ES.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Printf("Search error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	json.NewDecoder(res.Body).Decode(&r)

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	var results []model.Video
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		data, _ := json.Marshal(source)
		var v model.Video
		json.Unmarshal(data, &v)
		results = append(results, v)
	}

	c.JSON(http.StatusOK, results)
}
