package handlers

import (
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "พารามิเตอร์การค้นหาไม่ถูกต้อง"})
		return
	}

	bookings, err := h.adminService.GetBookings(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดในการดึงข้อมูล: " + err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "รหัสการจองไม่ถูกต้อง"})
		return
	}

	var req services.ManualApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง (ต้องระบุ action เป็น APPROVE หรือ REJECT)"})
		return
	}

	booking, err := h.adminService.ManualApprove(bookingID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := "อนุมัติการชำระเงินเรียบร้อยแล้ว"
	if req.Action == "REJECT" {
		message = "ปฏิเสธสลิปและยกเลิกการจองเรียบร้อยแล้ว คืนสต็อกและปลดล็อคคอร์ทแล้ว"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "รหัสการจองไม่ถูกต้อง"})
		return
	}

	var req struct {
		Items []services.BookingItemInput `json:"items" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลสินค้าไม่ถูกต้อง: " + err.Error()})
		return
	}

	booking, err := h.bookingService.AddItemsToBooking(bookingID, req.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "เพิ่มสินค้าลงในบิลและตัดสต็อกเรียบร้อยแล้ว",
		"data":    booking,
	})
}

// UpdateCheckIn handles PATCH /api/v1/admin/bookings/:id/check-in
func (h *AdminHandler) UpdateCheckIn(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รหัสการจองไม่ถูกต้อง"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"` // WAITING, CHECKED_IN, COMPLETED
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณาระบุสถานะ (WAITING, CHECKED_IN, COMPLETED)"})
		return
	}

	booking, err := h.adminService.UpdateCheckInStatus(bookingID, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "อัปเดตสถานะการเช็คอินเรียบร้อยแล้ว",
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
		"message":       "รันระบบตรวจสอบการจองหมดอายุเรียบร้อยแล้ว",
		"expired_count": expiredCount,
	})
}
