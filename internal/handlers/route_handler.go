package handlers

import (
	"mcp-server/internal/dtos"
	"mcp-server/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouteHandler struct {
	RouteService *services.RouteService
}

func NewRouteHandler(service *services.RouteService) *RouteHandler {
	return &RouteHandler{
		RouteService: service,
	}
}

func (h *RouteHandler) GetEta(c *gin.Context) {
	var params dtos.EtaRequest

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.RouteService.GetEta(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
