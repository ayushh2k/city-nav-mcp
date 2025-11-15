package routes

import (
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
	exchangeRateKey := os.Getenv("EXCHANGERATE_API_KEY")

	mcpGroup := router.Group("/mcp")
	mcpGroup.Use(middleware.JsonLoggingMiddleware())
	mcpGroup.Use(middleware.AuthMiddleware())
	{
		// --- Geo Tool ---
		geoService := services.NewGeoService(httpClient)
		geoHandler := handlers.NewGeoHandler(geoService)
		geoGroup := mcpGroup.Group("/geo")
		geoGroup.Use(middleware.PerIPRateLimiter(limitRate, limitBurst))
		{
			geoGroup.GET("/geocode", geoHandler.Geocode)
			geoGroup.GET("/nearby", geoHandler.Nearby)
		}

		// --- Weather Tool ---
		weatherService := services.NewWeatherService(httpClient)
		weatherHandler := handlers.NewWeatherHandler(weatherService)
		weatherGroup := mcpGroup.Group("/weather")
		weatherGroup.Use(middleware.PerIPRateLimiter(limitRate, limitBurst))
		{
			weatherGroup.GET("/forecast", weatherHandler.Forecast)
		}

		// --- Air Tool ---
		airService := services.NewAirService(httpClient, openAQKey)
		airHandler := handlers.NewAirHandler(airService)
		airGroup := mcpGroup.Group("/air")
		airGroup.Use(middleware.PerIPRateLimiter(limitRate, limitBurst))
		{
			airGroup.GET("/aqi", airHandler.GetAQI)
		}

		// --- Route Tool ---
		routeService := services.NewRouteService(httpClient)
		routeHandler := handlers.NewRouteHandler(routeService)
		routeGroup := mcpGroup.Group("/route")
		routeGroup.Use(middleware.PerIPRateLimiter(limitRate, limitBurst))
		{
			routeGroup.POST("/eta", routeHandler.GetEta)
		}

		// --- Calendar Tool ---
		calendarService := services.NewCalendarService(httpClient)
		calendarHandler := handlers.NewCalendarHandler(calendarService)
		calendarGroup := mcpGroup.Group("/calendar")
		calendarGroup.Use(middleware.PerIPRateLimiter(limitRate, limitBurst))
		{
			calendarGroup.GET("/holidays", calendarHandler.GetHolidays)
		}

		// --- Fx Tool ---
		fxService := services.NewFxService(httpClient, exchangeRateKey)
		fxHandler := handlers.NewFxHandler(fxService)
		fxGroup := mcpGroup.Group("/fx")
		fxGroup.Use(middleware.PerIPRateLimiter(limitRate, limitBurst))
		{
			fxGroup.GET("/convert", fxHandler.Convert)
		}

		// --- Wikidata Tool ---
		wikidataService := services.NewWikidataService(httpClient)
		wikidataHandler := handlers.NewWikidataHandler(wikidataService)
		wikidataGroup := mcpGroup.Group("/wikidata")
		wikidataGroup.Use(middleware.PerIPRateLimiter(limitRate, limitBurst))
		{
			wikidataGroup.POST("/query", wikidataHandler.Query)
		}
	}

	return router
}
