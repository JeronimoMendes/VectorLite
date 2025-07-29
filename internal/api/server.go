package api

import (
	api "VectorLite/internal/api/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Serve(port int) {
	r := gin.Default()
	
	// Database management endpoints
	r.POST("/databases", api.CreateDatabase)
	r.GET("/databases", api.ListDatabases)
	r.DELETE("/databases/:name", api.DeleteDatabase)
	
	// Entry and query endpoints
	r.POST("/entries", api.AddEntries)
	r.GET("/entries", api.ListEntries)
	r.POST("/query", api.Query)

	addr := fmt.Sprintf(":%d", port)
	r.Run(addr)
}
