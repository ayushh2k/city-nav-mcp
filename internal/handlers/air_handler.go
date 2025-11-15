package handlers

import (
	"mcp-server/internal/dtos"
	"mcp-server/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AirHandler struct {
	AirService *services.AirService
}

func NewAirHandler(service *services.AirService) *AirHandler {
	return &AirHandler{
		AirService: service,
	}
}

func (h *AirHandler) GetAQI(c *gin.Context) {
	var params dtos.AirQualityRequestParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.AirService.GetAQI(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
