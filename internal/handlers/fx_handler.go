package handlers

import (
	"mcp-server/internal/dtos"
	"mcp-server/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FxHandler struct {
	FxService *services.FxService
}

func NewFxHandler(service *services.FxService) *FxHandler {
	return &FxHandler{
		FxService: service,
	}
}

func (h *FxHandler) Convert(c *gin.Context) {
	var params dtos.FxRequestParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.FxService.ConvertCurrency(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
