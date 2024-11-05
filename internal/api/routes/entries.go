package api

import (
	"VectorLite/internal/engine"
	"VectorLite/internal/state"
	"VectorLite/internal/vector"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EntryRequest struct {
	Vectors   [][]float64         `json:"vectors" binding:"required"`
	Metadatas []map[string]string `json:"metadatas" binding:"required"`
}

func AddEntries(c *gin.Context) {
	var rb EntryRequest
	if err := c.ShouldBindJSON(&rb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Adding %d entries\n", len(rb.Vectors))
	for i, vec := range rb.Vectors {
		vector := vector.NewVector(vec...)
		engine.AddEntry(state.State.Database, *vector, rb.Metadatas[i])
	}
}

func ListEntries(c *gin.Context) {
	entries := engine.ListEntries(state.State.Database)
	log.Println("Listing entries")

	serializedEntries := make([]gin.H, len(entries))
	for i, entry := range entries {
		serializedEntries[i] = gin.H{
			"vector":   entry.Vector.Values, // Assuming entry.Vector returns a slice of float64
			"metadata": entry.Metadata,      // Assuming entry.Metadata returns a map[string]string
			"id":       entry.Id,
		}
	}

	c.JSON(http.StatusOK, gin.H{"entries": serializedEntries})
}
