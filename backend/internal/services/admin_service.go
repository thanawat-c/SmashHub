package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smashhub/internal/models"
)

type AdminFilterParams struct {
	Date    string `form:"date"`     // YYYY-MM-DD
	Status  string `form:"status"`   // PENDING_PAYMENT, CONFIRMED, MANUAL_REVIEW, EXPIRED, CANCELLED
	CourtID *uint  `form:"court_id"` // Filter by specific court
	Search  string `form:"search"`   // Search by booking_no, customer name, phone
}

type ManualApproveRequest struct {
	Action string `json:"action" binding:"required"` // APPROVE, REJECT
	Note   string `json:"note"`
}

type AdminService struct {
	db *gorm.DB
}

func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{db: db}
}

// GetBookings retrieves bookings with multiple filter criteria for Admin Dashboard & POS
func (s *AdminService) GetBookings(filter AdminFilterParams) ([]models.Booking, error) {
	query := s.db.Model(&models.Booking{}).
		Preload("User").
		Preload("Courts.Court").
		Preload("Items.Product").
		Preload("PaymentSlip")

	if filter.Status != "" {
		query = query.Where("bookings.status = ?", filter.Status)
	}

	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Joins("LEFT JOIN users ON users.id = bookings.user_id").
			Where("bookings.booking_no ILIKE ? OR users.display_name ILIKE ? OR users.phone ILIKE ?",
				searchTerm, searchTerm, searchTerm)
	}

	if filter.CourtID != nil && *filter.CourtID > 0 {
		query = query.Where("bookings.id IN (SELECT booking_id FROM booking_courts WHERE court_id = ?)", *filter.CourtID)
	}

	if filter.Date != "" {
		query = query.Where("bookings.id IN (SELECT booking_id FROM booking_courts WHERE play_date = ?)", filter.Date)
	}

	var bookings []models.Booking
	err := query.Order("bookings.created_at desc").Find(&bookings).Error
	return bookings, err
}

// ManualApprove processes admin review decisions on MANUAL_REVIEW bookings
func (s *AdminService) ManualApprove(bookingID uuid.UUID, req ManualApproveRequest) (*models.Booking, error) {
	var booking models.Booking

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Preload("Items").Preload("PaymentSlip").First(&booking, "id = ?", bookingID).Error; err != nil {
			return ErrBookingNotFound
		}

		if booking.Status != models.BookingStatusManualReview && booking.Status != models.BookingStatusPendingPayment {
			return fmt.Errorf("ไม่สามารถดำเนินการได้: สถานะปัจจุบันคือ %s", booking.Status)
		}

		now := time.Now()

		switch req.Action {
		case "APPROVE":
			booking.Status = models.BookingStatusConfirmed
			if err := tx.Save(&booking).Error; err != nil {
				return err
			}

			// Update slip status to VERIFIED
			if booking.PaymentSlip != nil {
				booking.PaymentSlip.VerificationStatus = models.SlipStatusVerified
				booking.PaymentSlip.VerifiedAt = &now
				if err := tx.Save(booking.PaymentSlip).Error; err != nil {
					return err
				}
			}

		case "REJECT":
			booking.Status = models.BookingStatusCancelled
			if err := tx.Save(&booking).Error; err != nil {
				return err
			}

			// Unlock court slots
			if err := tx.Model(&models.BookingCourt{}).Where("booking_id = ?", booking.ID).
				Update("is_slot_locked", false).Error; err != nil {
				return err
			}

			// Restock inventory
			for _, item := range booking.Items {
				if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
					Update("stock_qty", gorm.Expr("stock_qty + ?", item.Qty)).Error; err != nil {
					return err
				}
			}

		default:
			return errors.New("Action ไม่ถูกต้อง (ต้องเป็น APPROVE หรือ REJECT)")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload with full relations
	s.db.Preload("User").Preload("Courts.Court").Preload("Items.Product").Preload("PaymentSlip").First(&booking, "id = ?", bookingID)

	Hub.Broadcast("BOOKING_STATUS_CHANGED", map[string]interface{}{
		"booking_id": booking.ID,
		"booking_no": booking.BookingNo,
		"status":     booking.Status,
		"action":     req.Action,
	})

	return &booking, nil
}

// UpdateCheckInStatus updates check-in status (WAITING, CHECKED_IN, COMPLETED)
func (s *AdminService) UpdateCheckInStatus(bookingID uuid.UUID, status string) (*models.Booking, error) {
	if status != models.CheckInStatusWaiting && status != models.CheckInStatusCheckedIn && status != models.CheckInStatusCompleted {
		return nil, errors.New("สถานะการเช็คอินไม่ถูกต้อง")
	}

	var booking models.Booking
	if err := s.db.First(&booking, "id = ?", bookingID).Error; err != nil {
		return nil, ErrBookingNotFound
	}

	booking.CheckInStatus = status
	if err := s.db.Save(&booking).Error; err != nil {
		return nil, err
	}

	s.db.Preload("User").Preload("Courts.Court").Preload("Items.Product").Preload("PaymentSlip").First(&booking, "id = ?", bookingID)

	Hub.Broadcast("CHECKIN_UPDATED", map[string]interface{}{
		"booking_id":      booking.ID,
		"booking_no":      booking.BookingNo,
		"check_in_status": status,
	})

	return &booking, nil
}

// UpdateCourtStatus toggles court status (ACTIVE or MAINTENANCE) and hourly rate
func (s *AdminService) UpdateCourtStatus(courtID uint, status string, normalRate, peakRate float64) (*models.Court, error) {
	var court models.Court
	if err := s.db.First(&court, courtID).Error; err != nil {
		return nil, fmt.Errorf("ไม่พบคอร์ทรหัส %d", courtID)
	}

	if status != "" {
		court.Status = status
	}
	if normalRate > 0 {
		court.NormalRate = normalRate
	}
	if peakRate > 0 {
		court.PeakRate = peakRate
	}

	if err := s.db.Save(&court).Error; err != nil {
		return nil, err
	}

	Hub.Broadcast("COURT_UPDATED", map[string]interface{}{
		"court_id":    court.ID,
		"court_number": court.CourtNumber,
		"status":      court.Status,
	})

	return &court, nil
}

// RestockProduct updates stock quantity for mini-store products
func (s *AdminService) RestockProduct(productID uint, addQty int, newPrice float64) (*models.Product, error) {
	var product models.Product
	if err := s.db.First(&product, productID).Error; err != nil {
		return nil, fmt.Errorf("ไม่พบสินค้ารหัส %d", productID)
	}

	if addQty != 0 {
		product.StockQty += addQty
		if product.StockQty < 0 {
			product.StockQty = 0
		}
	}
	if newPrice > 0 {
		product.Price = newPrice
	}

	if err := s.db.Save(&product).Error; err != nil {
		return nil, err
	}

	Hub.Broadcast("PRODUCT_RESTOCKED", map[string]interface{}{
		"product_id": product.ID,
		"name":       product.Name,
		"stock_qty":  product.StockQty,
		"price":      product.Price,
	})

	return &product, nil
}

