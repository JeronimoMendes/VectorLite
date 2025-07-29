package api

import (
	"VectorLite/internal/state"
	"VectorLite/internal/vector"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EntryRequest struct {
	Database  string              `json:"database" binding:"required"`
	Vectors   [][]float64         `json:"vectors" binding:"required"`
	Metadatas []map[string]string `json:"metadatas" binding:"required"`
}

func AddEntries(c *gin.Context) {
	var rb EntryRequest
	if err := c.ShouldBindJSON(&rb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database, err := state.State.DatabaseManager.GetDatabase(rb.Database)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Adding %d entries to database %s\n", len(rb.Vectors), rb.Database)
	for i, vec := range rb.Vectors {
		vector := vector.NewVector(vec...)
		database.AddEntry(*vector, rb.Metadatas[i])
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "entries added successfully"})
}

func ListEntries(c *gin.Context) {
	databaseName := c.Query("database")
	if databaseName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database parameter is required"})
		return
	}

	database, err := state.State.DatabaseManager.GetDatabase(databaseName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	entries := database.ListEntries()
	log.Printf("Listing entries from database %s\n", databaseName)

	serializedEntries := make([]gin.H, len(entries))
	for i, entry := range entries {
		serializedEntries[i] = gin.H{
			"vector":   entry.Vector.Values,
			"metadata": entry.Metadata,
			"id":       entry.Id,
		}
	}

	c.JSON(http.StatusOK, gin.H{"entries": serializedEntries})
}
