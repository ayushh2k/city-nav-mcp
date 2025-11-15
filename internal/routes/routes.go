package routes

import (
	"log"
	"mcp-server/internal/handlers"
	"mcp-server/internal/middleware"
	"mcp-server/internal/services"
	"net/http"
	"os"
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

	openAQKey := os.Getenv("OPENAQ_API_KEY")
	if openAQKey == "" {
		log.Fatal("FATAL: OPENAQ_API_KEY environment variable not set.")
	}

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

		airService := services.NewAirService(httpClient, openAQKey)
		airHandler := handlers.NewAirHandler(airService)
		airGroup := mcpGroup.Group("/air")
		airGroup.Use(middleware.PerIPRateLimiter(limitRate, limitBurst))
		{
			airGroup.GET("/aqi", airHandler.GetAQI)
		}

		routeService := services.NewRouteService(httpClient)
		routeHandler := handlers.NewRouteHandler(routeService)
		routeGroup := mcpGroup.Group("/route")
		routeGroup.Use(middleware.PerIPRateLimiter(limitRate, limitBurst))
		{
			routeGroup.POST("/eta", routeHandler.GetEta)
		}
	}

	return router
}
