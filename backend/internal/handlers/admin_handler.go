package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"smashhub/internal/services"
)

type AdminHandler struct {
	adminService   *services.AdminService
	bookingService *services.BookingService
	cronWorker     *services.CronWorker
}

func NewAdminHandler(
	adminService *services.AdminService,
	bookingService *services.BookingService,
	cronWorker *services.CronWorker,
) *AdminHandler {
	return &AdminHandler{
		adminService:   adminService,
		bookingService: bookingService,
		cronWorker:     cronWorker,
	}
}

// GetBookings handles GET /api/v1/admin/bookings
func (h *AdminHandler) GetBookings(c *gin.Context) {
	var filter services.AdminFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	bookings, err := h.adminService.GetBookings(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  bookings,
		"total": len(bookings),
	})
}

// ManualApprove handles PATCH /api/v1/admin/bookings/:id/manual-approve
func (h *AdminHandler) ManualApprove(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var req services.ManualApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload (action must be APPROVE or REJECT)"})
		return
	}

	booking, err := h.adminService.ManualApprove(bookingID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := "Booking payment successfully approved"
	if req.Action == "REJECT" {
		message = "Slip rejected and booking cancelled. Court released and inventory restored."
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    booking,
	})
}

// AddItems handles POST /api/v1/admin/bookings/:id/add-items (POS feature)
func (h *AdminHandler) AddItems(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var req struct {
		Items []services.BookingItemInput `json:"items" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid items payload: " + err.Error()})
		return
	}

	booking, err := h.bookingService.AddItemsToBooking(bookingID, req.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Items successfully added and inventory deducted",
		"data":    booking,
	})
}

// UpdateCheckIn handles PATCH /api/v1/admin/bookings/:id/check-in
func (h *AdminHandler) UpdateCheckIn(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"` // WAITING, CHECKED_IN, COMPLETED
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide check-in status (WAITING, CHECKED_IN, COMPLETED)"})
		return
	}

	booking, err := h.adminService.UpdateCheckInStatus(bookingID, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Check-in status updated successfully",
		"data":    booking,
	})
}

// TriggerCron handles POST /api/v1/admin/cron/trigger (for testing / demo)
func (h *AdminHandler) TriggerCron(c *gin.Context) {
	expiredCount, err := h.cronWorker.ExpireOverdueBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Expired bookings check completed successfully",
		"expired_count": expiredCount,
	})
}

// UpdateCourt handles PATCH /api/v1/admin/courts/:id
func (h *AdminHandler) UpdateCourt(c *gin.Context) {
	idParam := c.Param("id")
	var courtID uint
	if _, err := fmt.Sscanf(idParam, "%d", &courtID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid court ID"})
		return
	}

	var req struct {
		Status     string  `json:"status"` // ACTIVE, MAINTENANCE
		NormalRate float64 `json:"normal_rate"`
		PeakRate   float64 `json:"peak_rate"`
	}
	_ = c.ShouldBindJSON(&req)

	court, err := h.adminService.UpdateCourtStatus(courtID, req.Status, req.NormalRate, req.PeakRate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Court status updated successfully",
		"data":    court,
	})
}

// RestockProduct handles POST /api/v1/admin/products/:id/restock
func (h *AdminHandler) RestockProduct(c *gin.Context) {
	idParam := c.Param("id")
	var productID uint
	if _, err := fmt.Sscanf(idParam, "%d", &productID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req struct {
		AddQty int     `json:"add_qty"`
		Price  float64 `json:"price"`
	}
	_ = c.ShouldBindJSON(&req)

	product, err := h.adminService.RestockProduct(productID, req.AddQty, req.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory restocked successfully",
		"data":    product,
	})
}

