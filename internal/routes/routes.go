package routes

import (
	"mcp-server/internal/handlers"
	"mcp-server/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "MCP Server is running"})
	})

	geoService := services.NewGeoService()
	geoHandler := handlers.NewGeoHandler(geoService)

	geoGroup := router.Group("/mcp-geo")
	{
		geoGroup.GET("/geocode", geoHandler.Geocode)

		geoGroup.GET("/nearby", geoHandler.Nearby)
	}

	return router
}
