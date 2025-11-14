package routes

import (
	"mcp-server/internal/handlers"
	"mcp-server/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "MCP Server is running"})
	})

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	geoService := services.NewGeoService(httpClient)
	geoHandler := handlers.NewGeoHandler(geoService)

	geoGroup := router.Group("/mcp-geo")
	{
		geoGroup.GET("/geocode", geoHandler.Geocode)
		geoGroup.GET("/nearby", geoHandler.Nearby)
	}

	weatherService := services.NewWeatherService(httpClient)
	weatherHandler := handlers.NewWeatherHandler(weatherService)

	weatherGroup := router.Group("/mcp-weather")
	{
		weatherGroup.GET("/forecast", weatherHandler.Forecast)
	}

	return router
}
