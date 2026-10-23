package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"smashhub/internal/services"
)

type CourtHandler struct {
	bookingService *services.BookingService
}

func NewCourtHandler(bookingService *services.BookingService) *CourtHandler {
	return &CourtHandler{bookingService: bookingService}
}

// GetCourtAvailability handles GET /api/v1/courts/availability?date=YYYY-MM-DD
func (h *CourtHandler) GetCourtAvailability(c *gin.Context) {
	date := c.Query("date")
	availability, err := h.bookingService.GetCourtAvailability(date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, availability)
}

// GetProducts handles GET /api/v1/products
func (h *CourtHandler) GetProducts(c *gin.Context) {
	products, err := h.bookingService.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  products,
		"total": len(products),
	})
}
