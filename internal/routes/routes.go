package routes

import (
	"mcp-server/internal/handlers"
	"mcp-server/internal/middleware"
	"mcp-server/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.SetTrustedProxies(nil)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "MCP Server is running"})
	})

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	limitRate := rate.Every(time.Minute / 30)
	limitBurst := 30

	mcpGroup := router.Group("/mcp")
	mcpGroup.Use(middleware.JsonLoggingMiddleware())
	mcpGroup.Use(middleware.AuthMiddleware())
	{
		geoService := services.NewGeoService(httpClient)
		geoHandler := handlers.NewGeoHandler(geoService)

		geoGroup := mcpGroup.Group("/geo")
		geoGroup.Use(middleware.PerIPRateLimiter(limitRate, limitBurst))
		{
			geoGroup.GET("/geocode", geoHandler.Geocode)
			geoGroup.GET("/nearby", geoHandler.Nearby)
		}

		weatherService := services.NewWeatherService(httpClient)
		weatherHandler := handlers.NewWeatherHandler(weatherService)

		weatherGroup := mcpGroup.Group("/weather")
		weatherGroup.Use(middleware.PerIPRateLimiter(limitRate, limitBurst))
		{
			weatherGroup.GET("/forecast", weatherHandler.Forecast)
		}
	}

	return router
}
