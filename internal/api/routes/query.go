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

	vector := vector.NewVector(rb.QueryVector...)
	log.Println(fmt.Sprintf("k=%d, metric=%s", rb.K, rb.Metric))
	results := state.State.Database.Query(vector, rb.K, rb.Metric)

	log.Println(fmt.Sprintf("Got %d results", len(results)))
	serializedEntries := make([]gin.H, len(results))
	for i, entry := range results {
		serializedEntries[i] = gin.H{
			"vector":   entry.Vector.Values, // Assuming entry.Vector returns a slice of float64
			"metadata": entry.Metadata,      // Assuming entry.Metadata returns a map[string]string
			"id":       entry.Id,
		}
	}

	c.JSON(http.StatusOK, gin.H{"entries": serializedEntries})
}
