package handlers

import (
	"mcp-server/internal/dtos"
	"mcp-server/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GeoHandler struct {
	GeoService *services.GeoService
}

func NewGeoHandler(service *services.GeoService) *GeoHandler {
	return &GeoHandler{
		GeoService: service,
	}
}

func (h *GeoHandler) Geocode(c *gin.Context) {
	var params dtos.GeocodeRequestParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.GeoService.GetGeocode(params)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *GeoHandler) Nearby(c *gin.Context) {
	var params dtos.NearbyRequestParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.GeoService.GetNearby(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
