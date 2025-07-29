package api

import (
	"VectorLite/internal/state"
	"VectorLite/internal/vector"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueryRequest struct {
	Database    string    `json:"database" binding:"required"`
	QueryVector []float64 `json:"vector" binding:"required"`
	K           int       `json:"k" binding:"required"`
	Metric      string    `json:"metric" binding:"required"`
}

func Query(c *gin.Context) {
	var rb QueryRequest
	if err := c.ShouldBindJSON(&rb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database, err := state.State.DatabaseManager.GetDatabase(rb.Database)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	vector := vector.NewVector(rb.QueryVector...)
	log.Println(fmt.Sprintf("database=%s, k=%d, metric=%s", rb.Database, rb.K, rb.Metric))
	results := database.Query(vector, rb.K, rb.Metric)

	log.Println(fmt.Sprintf("Got %d results", len(results)))
	serializedEntries := make([]gin.H, len(results))
	for i, entry := range results {
		serializedEntries[i] = gin.H{
			"vector":   entry.Vector.Values,
			"metadata": entry.Metadata,
			"id":       entry.Id,
		}
	}

	c.JSON(http.StatusOK, gin.H{"entries": serializedEntries})
}
