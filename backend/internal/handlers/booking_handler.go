package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"smashhub/internal/models"
	"smashhub/internal/services"
)

type BookingHandler struct {
	bookingService *services.BookingService
	slipService    *services.SlipService
	db             *gorm.DB
}

func NewBookingHandler(bookingService *services.BookingService, slipService *services.SlipService, db *gorm.DB) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		slipService:    slipService,
		db:             db,
	}
}

// CreateBooking handles POST /api/v1/bookings
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req services.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลการจองไม่ถูกต้อง: " + err.Error()})
		return
	}

	result, err := h.bookingService.CreateBooking(req)
	if err != nil {
		if errors.Is(err, services.ErrCourtAlreadyBooked) {
			// CRITICAL: Return 409 Conflict when partial unique index rejects overlapping slot
			c.JSON(http.StatusConflict, gin.H{
				"error":   "คอร์ทที่เลือกเพิ่งถูกจองไปเมื่อสักครู่",
				"details": "มีผู้ใช้อื่นจองสล็อตเวลานี้ไปก่อนหน้า กรุณาเลือกช่วงเวลาหรือคอร์ทใหม่อีกครั้ง",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "สร้างการจองสำเร็จ กรุณาชำระเงินและแนบสลิปภายใน 10 นาที",
		"data":    result,
	})
}

// GetBookingByID handles GET /api/v1/bookings/:id
func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รหัสการจองไม่ถูกต้อง"})
		return
	}

	var booking models.Booking
	if err := h.db.Preload("User").
		Preload("Courts.Court").
		Preload("Items.Product").
		Preload("PaymentSlip").
		First(&booking, "id = ?", bookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลการจอง"})
		return
	}

	promptPayQR := services.GeneratePromptPayPayload("0987654321", booking.NetTotal)

	c.JSON(http.StatusOK, gin.H{
		"data":            booking,
		"promptpay_qr":    promptPayQR,
		"promptpay_phone": "098-7-65432-1",
		"promptpay_name":  "บจก. สแมชฮับ แบดมินตัน (SmashHub)",
	})
}

// VerifySlip handles POST /api/v1/bookings/:id/verify-slip
func (h *BookingHandler) VerifySlip(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รหัสการจองไม่ถูกต้อง"})
		return
	}

	var req services.VerifySlipRequest
	// Support both JSON body and multipart form (file upload / text)
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลสลิปไม่ถูกต้อง: " + err.Error()})
		return
	}

	// Handle optional uploaded slip file if present
	file, fileErr := c.FormFile("slip_image")
	if fileErr == nil && file != nil {
		req.SlipImageURL = "https://placehold.co/600x800/004d40/ffffff?text=Uploaded+" + file.Filename
	}

	result, err := h.slipService.VerifySlip(bookingID, req)
	if err != nil {
		if errors.Is(err, services.ErrDuplicateTransRef) {
			c.JSON(http.StatusConflict, gin.H{
				"error":   err.Error(),
				"details": result,
			})
			return
		}
		if errors.Is(err, services.ErrUnderpaid) {
			c.JSON(http.StatusPaymentRequired, gin.H{
				"error":   err.Error(),
				"details": result,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "details": result})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": result.Message,
		"data":    result,
	})
}

// GetMyBookings handles GET /api/v1/bookings/my-bookings
func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-Id")
	}
	phone := c.Query("phone")
	lineUserID := c.Query("line_user_id")

	bookings, err := h.bookingService.GetMyBookings(userID, phone, lineUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  bookings,
		"count": len(bookings),
	})
}

