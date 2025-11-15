package handlers

import (
	"mcp-server/internal/dtos"
	"mcp-server/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CalendarHandler struct {
	CalendarService *services.CalendarService
}

func NewCalendarHandler(service *services.CalendarService) *CalendarHandler {
	return &CalendarHandler{
		CalendarService: service,
	}
}

func (h *CalendarHandler) GetHolidays(c *gin.Context) {
	var params dtos.HolidayRequestParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.CalendarService.GetHolidays(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
