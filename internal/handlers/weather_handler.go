package handlers

import (
	"mcp-server/internal/dtos"
	"mcp-server/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	WeatherService *services.WeatherService
}

func NewWeatherHandler(service *services.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		WeatherService: service,
	}
}

func (h *WeatherHandler) Forecast(c *gin.Context) {
	var params dtos.ForecastRequestParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.WeatherService.GetForecast(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
