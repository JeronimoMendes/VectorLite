package api

import (
	"VectorLite/internal/algorithms"
	"VectorLite/internal/algorithms/bruteforce"
	"VectorLite/internal/algorithms/hnsw"
	"VectorLite/internal/state"
	"log"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateDatabaseRequest struct {
	Name      string                 `json:"name" binding:"required"`
	Algorithm string                 `json:"algorithm" binding:"required"`
	Settings  map[string]interface{} `json:"settings,omitempty"`
}

func CreateDatabase(c *gin.Context) {
	var req CreateDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var algorithm algorithms.SearchAlgorithm

	switch req.Algorithm {
	case "bruteforce":
		algorithm = bruteforce.New()
	case "hnsw":
		// Default HNSW parameters
		M := 16
		efConstruction := 200
		mL := 1.0 / math.Log(2.0)
		algorithm = hnsw.New(M, efConstruction, mL)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported algorithm: " + req.Algorithm})
		return
	}

	err := state.State.DatabaseManager.CreateDatabase(req.Name, algorithm)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Created database %s with algorithm %s\n", req.Name, req.Algorithm)
	c.JSON(http.StatusCreated, gin.H{
		"message":   "database created successfully",
		"name":      req.Name,
		"algorithm": req.Algorithm,
	})
}

func ListDatabases(c *gin.Context) {
	databases := state.State.DatabaseManager.ListDatabases()
	log.Printf("Listing %d databases\n", len(databases))
	
	c.JSON(http.StatusOK, gin.H{"databases": databases})
}

func DeleteDatabase(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name is required"})
		return
	}

	err := state.State.DatabaseManager.DeleteDatabase(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Deleted database %s\n", name)
	c.JSON(http.StatusOK, gin.H{"message": "database deleted successfully"})
}