package handlers

import (
	"mcp-server/internal/dtos"
	"mcp-server/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WikidataHandler struct {
	WikidataService *services.WikidataService
}

func NewWikidataHandler(service *services.WikidataService) *WikidataHandler {
	return &WikidataHandler{
		WikidataService: service,
	}
}

func (h *WikidataHandler) Query(c *gin.Context) {
	var params dtos.WikidataRequest

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.WikidataService.Query(params.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
