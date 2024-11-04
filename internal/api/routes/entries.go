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
	Vector   []float64         `json:"vector" binding:"required"`
	Metadata map[string]string `json:"metadata" binding:"required"`
}

func AddEntries(c *gin.Context) {
	var rb []EntryRequest
	if err := c.ShouldBindJSON(&rb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Adding %d entries\n", len(rb))
	for _, entry := range rb {
		vector := vector.NewVector(entry.Vector...)
		engine.AddEntry(state.State.Database, *vector, entry.Metadata)
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
