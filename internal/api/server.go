package api

import (
	api "VectorLite/internal/api/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Serve(port int) {
	r := gin.Default()
	r.POST("/entries", api.AddEntries)
	r.GET("/entries", api.ListEntries)
	r.POST("/query", api.Query)

	addr := fmt.Sprintf(":%d", port)
	r.Run(addr)
}
