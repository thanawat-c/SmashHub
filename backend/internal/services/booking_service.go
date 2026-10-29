package services

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"smashhub/internal/models"
)

var (
	ErrCourtAlreadyBooked  = errors.New("the selected court was just booked by another user")
	ErrInsufficientStock   = errors.New("insufficient inventory in stock")
	ErrBookingNotFound     = errors.New("booking not found")
	ErrInvalidBookingState = errors.New("invalid booking status for this operation")
)

type CourtSlotInput struct {
	CourtID   uint `json:"court_id" binding:"required"`
	StartHour int  `json:"start_hour" binding:"required,min=10,max=21"`
}

type BookingItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Qty       int  `json:"qty" binding:"required,min=1"`
}

type CreateBookingRequest struct {
	CustomerName  string             `json:"customer_name" binding:"required"`
	CustomerPhone string             `json:"customer_phone" binding:"required"`
	LineUserID    *string            `json:"line_user_id"`
	PlayDate      string             `json:"play_date" binding:"required"` // YYYY-MM-DD
	CourtSlots    []CourtSlotInput   `json:"court_slots" binding:"required,min=1"`
	Items         []BookingItemInput `json:"items"`
}

type TimeSlotInfo struct {
	StartHour  int     `json:"start_hour"`
	EndHour    int     `json:"end_hour"`
	HourlyRate float64 `json:"hourly_rate"`
	Status     string  `json:"status"` // AVAILABLE, LOCKED_PENDING, BOOKED, MAINTENANCE
	BookingNo  *string `json:"booking_no,omitempty"`
}

type CourtAvailability struct {
	CourtID     uint           `json:"court_id"`
	CourtNumber int            `json:"court_number"`
	CourtName   string         `json:"court_name"`
	CourtType   string         `json:"court_type"`
	NormalRate  float64        `json:"normal_rate"`
	PeakRate    float64        `json:"peak_rate"`
	CourtStatus string         `json:"court_status"`
	Slots       []TimeSlotInfo `json:"slots"`
}

type AvailabilityResponse struct {
	PlayDate string              `json:"play_date"`
	Courts   []CourtAvailability `json:"courts"`
}

type BookingResponse struct {
	Booking        *models.Booking `json:"booking"`
	PromptPayQR    string          `json:"promptpay_qr"`
	PromptPayPhone string          `json:"promptpay_phone"`
	PromptPayName  string          `json:"promptpay_name"`
}

type BookingService struct {
	db *gorm.DB
}

func NewBookingService(db *gorm.DB) *BookingService {
	return &BookingService{db: db}
}

// GetCourtAvailability returns the time slot matrix (10:00 - 22:00) for all courts on a given date
func (s *BookingService) GetCourtAvailability(dateStr string) (*AvailabilityResponse, error) {
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	playDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format (must be YYYY-MM-DD): %w", err)
	}

	var courts []models.Court
	if err := s.db.Order("court_number asc").Find(&courts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch courts: %w", err)
	}

	// Fetch active locked court slots for that date
	var lockedSlots []models.BookingCourt
	err = s.db.Preload("Court").
		Joins("JOIN bookings ON bookings.id = booking_courts.booking_id").
		Where("booking_courts.play_date = ? AND booking_courts.is_slot_locked = ?", playDate, true).
		Select("booking_courts.*, bookings.status as booking_status, bookings.expires_at as booking_expires_at, bookings.booking_no as booking_no").
		Find(&lockedSlots).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locked court slots: %w", err)
	}

	// Map locked slots by court_id and start_hour
	type SlotState struct {
		status    string
		bookingNo string
	}
	slotMap := make(map[string]SlotState)

	now := time.Now()
	for _, ls := range lockedSlots {
		var booking models.Booking
		s.db.Select("status, expires_at, booking_no").First(&booking, ls.BookingID)

		key := fmt.Sprintf("%d-%d", ls.CourtID, ls.StartHour)
		if booking.Status == models.BookingStatusPendingPayment {
			if booking.ExpiresAt != nil && booking.ExpiresAt.After(now) {
				slotMap[key] = SlotState{status: "LOCKED_PENDING", bookingNo: booking.BookingNo}
			}
		} else if booking.Status == models.BookingStatusConfirmed || booking.Status == models.BookingStatusManualReview {
			slotMap[key] = SlotState{status: "BOOKED", bookingNo: booking.BookingNo}
		}
	}

	resultCourts := make([]CourtAvailability, 0, len(courts))
	for _, court := range courts {
		ca := CourtAvailability{
			CourtID:     court.ID,
			CourtNumber: court.CourtNumber,
			CourtName:   court.Name,
			CourtType:   court.CourtType,
			NormalRate:  court.NormalRate,
			PeakRate:    court.PeakRate,
			CourtStatus: court.Status,
			Slots:       make([]TimeSlotInfo, 0, 12),
		}

		for hour := 10; hour < 22; hour++ {
			rate := court.NormalRate
			if hour >= 17 {
				rate = court.PeakRate
			}

			status := "AVAILABLE"
			var bookingNo *string

			if court.Status == "MAINTENANCE" {
				status = "MAINTENANCE"
			} else {
				key := fmt.Sprintf("%d-%d", court.ID, hour)
				if val, exists := slotMap[key]; exists {
					status = val.status
					bNo := val.bookingNo
					bookingNo = &bNo
				}
			}

			ca.Slots = append(ca.Slots, TimeSlotInfo{
				StartHour:  hour,
				EndHour:    hour + 1,
				HourlyRate: rate,
				Status:     status,
				BookingNo:  bookingNo,
			})
		}
		resultCourts = append(resultCourts, ca)
	}

	return &AvailabilityResponse{
		PlayDate: dateStr,
		Courts:   resultCourts,
	}, nil
}

// GetProducts returns all active mini-store products
func (s *BookingService) GetProducts() ([]models.Product, error) {
	var products []models.Product
	err := s.db.Where("is_active = ?", true).Order("category asc, id asc").Find(&products).Error
	return products, err
}

// CreateBooking creates a new booking atomically: locks stock, locks court slots, calculates server-side totals
func (s *BookingService) CreateBooking(req CreateBookingRequest) (*BookingResponse, error) {
	playDate, err := time.Parse("2006-01-02", req.PlayDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format (YYYY-MM-DD): %w", err)
	}

	var createdBooking *models.Booking

	// Execute within a single database transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Find or create customer user
		var user models.User
		userQuery := tx.Where("phone = ?", req.CustomerPhone)
		if req.LineUserID != nil && *req.LineUserID != "" {
			userQuery = tx.Where("phone = ? OR line_user_id = ?", req.CustomerPhone, *req.LineUserID)
		}

		if err := userQuery.First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				user = models.User{
					DisplayName: req.CustomerName,
					Phone:       req.CustomerPhone,
					LineUserID:  req.LineUserID,
					Role:        models.RoleCustomer,
				}
				if err := tx.Create(&user).Error; err != nil {
					return fmt.Errorf("failed to create user record: %w", err)
				}
			} else {
				return err
			}
		}

		// 2. Lock and deduct product inventory using SELECT ... FOR UPDATE
		var bookingItems []models.BookingItem
		var itemTotal float64 = 0

		for _, itemInput := range req.Items {
			if itemInput.Qty <= 0 {
				continue
			}
			var product models.Product
			// Row-level lock on product row
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, itemInput.ProductID).Error; err != nil {
				return fmt.Errorf("product not found with ID %d: %w", itemInput.ProductID, err)
			}

			if product.StockQty < itemInput.Qty {
				return fmt.Errorf("product %s has insufficient stock (%d remaining)", product.Name, product.StockQty)
			}

			// Deduct stock
			product.StockQty -= itemInput.Qty
			if err := tx.Save(&product).Error; err != nil {
				return fmt.Errorf("failed to deduct stock for %s: %w", product.Name, err)
			}

			subtotal := product.Price * float64(itemInput.Qty)
			itemTotal += subtotal

			bookingItems = append(bookingItems, models.BookingItem{
				ProductID: product.ID,
				Qty:       itemInput.Qty,
				UnitPrice: product.Price,
				Subtotal:  subtotal,
			})
		}

		// 3. Verify and calculate court slot rates
		var bookingCourts []models.BookingCourt
		var courtTotal float64 = 0

		courtCache := make(map[uint]models.Court)
		for _, slot := range req.CourtSlots {
			court, exists := courtCache[slot.CourtID]
			if !exists {
				if err := tx.First(&court, slot.CourtID).Error; err != nil {
					return fmt.Errorf("court not found with ID %d: %w", slot.CourtID, err)
				}
				if court.Status == "MAINTENANCE" {
					return fmt.Errorf("court %s is currently under maintenance and cannot be booked", court.Name)
				}
				courtCache[slot.CourtID] = court
			}

			hourlyRate := court.NormalRate
			if slot.StartHour >= 17 {
				hourlyRate = court.PeakRate
			}
			courtTotal += hourlyRate

			bookingCourts = append(bookingCourts, models.BookingCourt{
				CourtID:      court.ID,
				PlayDate:     playDate,
				StartHour:    slot.StartHour,
				EndHour:      slot.StartHour + 1,
				HourlyRate:   hourlyRate,
				IsSlotLocked: true,
			})
		}

		netTotal := courtTotal + itemTotal

		// 4. Generate unique Booking No: BK-YYYYMMDD-XXXX
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randSuffix := r.Intn(9000) + 1000
		bookingNo := fmt.Sprintf("BK-%s-%04d", time.Now().Format("20060102"), randSuffix)

		// 5. Create Booking record with 10-minute lock expiration
		expiresAt := time.Now().Add(10 * time.Minute)
		booking := models.Booking{
			BookingNo:     bookingNo,
			UserID:        &user.ID,
			Status:        models.BookingStatusPendingPayment,
			CourtTotal:    courtTotal,
			ItemTotal:     itemTotal,
			NetTotal:      netTotal,
			CheckInStatus: models.CheckInStatusWaiting,
			ExpiresAt:     &expiresAt,
		}

		if err := tx.Create(&booking).Error; err != nil {
			return fmt.Errorf("failed to create booking record: %w", err)
		}

		// 6. Insert booking courts (Enforces Partial Unique Index: idx_unique_active_court_slot)
		for i := range bookingCourts {
			bookingCourts[i].BookingID = booking.ID
			if err := tx.Create(&bookingCourts[i]).Error; err != nil {
				// Detect duplicate key constraint violation
				if strings.Contains(err.Error(), "idx_unique_active_court_slot") ||
					strings.Contains(err.Error(), "duplicate key") {
					return ErrCourtAlreadyBooked
				}
				return fmt.Errorf("failed to lock court slot: %w", err)
			}
		}

		// 7. Insert booking items
		for i := range bookingItems {
			bookingItems[i].BookingID = booking.ID
			if err := tx.Create(&bookingItems[i]).Error; err != nil {
				return fmt.Errorf("failed to record booking items: %w", err)
			}
		}

		var completeBooking models.Booking
		if err := tx.Preload("User").Preload("Courts.Court").Preload("Items.Product").First(&completeBooking, "id = ?", booking.ID).Error; err == nil {
			createdBooking = &completeBooking
		} else {
			booking.Courts = bookingCourts
			booking.Items = bookingItems
			booking.User = &user
			createdBooking = &booking
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Generate PromptPay QR details
	storePhone := "0987654321"
	storeName := "SmashHub Badminton Arena Co., Ltd."
	qrPayload := GeneratePromptPayPayload(storePhone, createdBooking.NetTotal)

	// Broadcast SSE event
	Hub.Broadcast("BOOKING_CREATED", map[string]interface{}{
		"booking_id":  createdBooking.ID,
		"booking_no":  createdBooking.BookingNo,
		"customer":    createdBooking.User.DisplayName,
		"net_total":   createdBooking.NetTotal,
		"status":      createdBooking.Status,
		"created_at":  createdBooking.CreatedAt,
	})

	return &BookingResponse{
		Booking:        createdBooking,
		PromptPayQR:    qrPayload,
		PromptPayPhone: "098-7-65432-1",
		PromptPayName:  storeName,
	}, nil
}

// AddItemsToBooking allows Admin POS to add snacks, drinks, or shuttlecocks to an existing booking
func (s *BookingService) AddItemsToBooking(bookingID uuid.UUID, items []BookingItemInput) (*models.Booking, error) {
	var updatedBooking models.Booking

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var booking models.Booking
		if err := tx.Preload("Items").First(&booking, "id = ?", bookingID).Error; err != nil {
			return ErrBookingNotFound
		}

		if booking.Status == models.BookingStatusExpired || booking.Status == models.BookingStatusCancelled {
			return ErrInvalidBookingState
		}

		var additionalItemTotal float64 = 0

		for _, itemInput := range items {
			if itemInput.Qty <= 0 {
				continue
			}

			var product models.Product
			// Lock product row
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, itemInput.ProductID).Error; err != nil {
				return fmt.Errorf("product not found with ID %d", itemInput.ProductID)
			}

			if product.StockQty < itemInput.Qty {
				return fmt.Errorf("product %s has insufficient stock (%d remaining)", product.Name, product.StockQty)
			}

			// Deduct stock
			product.StockQty -= itemInput.Qty
			if err := tx.Save(&product).Error; err != nil {
				return err
			}

			subtotal := product.Price * float64(itemInput.Qty)
			additionalItemTotal += subtotal

			// Create booking item
			newItem := models.BookingItem{
				BookingID: booking.ID,
				ProductID: product.ID,
				Qty:       itemInput.Qty,
				UnitPrice: product.Price,
				Subtotal:  subtotal,
			}
			if err := tx.Create(&newItem).Error; err != nil {
				return err
			}
		}

		// Update booking totals
		booking.ItemTotal += additionalItemTotal
		booking.NetTotal += additionalItemTotal
		if err := tx.Save(&booking).Error; err != nil {
			return err
		}

		// Reload with relations
		tx.Preload("User").Preload("Courts.Court").Preload("Items.Product").Preload("PaymentSlip").First(&updatedBooking, "id = ?", booking.ID)
		return nil
	})

	if err != nil {
		return nil, err
	}

	Hub.Broadcast("ITEMS_ADDED", map[string]interface{}{
		"booking_id": updatedBooking.ID,
		"booking_no": updatedBooking.BookingNo,
		"net_total":  updatedBooking.NetTotal,
	})

	return &updatedBooking, nil
}

// GetMyBookings retrieves booking history for a specific customer
func (s *BookingService) GetMyBookings(userID, phone, lineUserID string) ([]models.Booking, error) {
	query := s.db.Model(&models.Booking{}).
		Preload("User").
		Preload("Courts.Court").
		Preload("Items.Product").
		Preload("PaymentSlip")

	hasFilter := false
	if userID != "" {
		if uid, err := uuid.Parse(userID); err == nil {
			query = query.Where("bookings.user_id = ?", uid)
			hasFilter = true
		}
	}

	if !hasFilter && phone != "" {
		query = query.Joins("JOIN users ON users.id = bookings.user_id").
			Where("users.phone = ?", phone)
		hasFilter = true
	} else if !hasFilter && lineUserID != "" {
		query = query.Joins("JOIN users ON users.id = bookings.user_id").
			Where("users.line_user_id = ?", lineUserID)
		hasFilter = true
	}

	var bookings []models.Booking
	if err := query.Order("bookings.created_at DESC").Limit(50).Find(&bookings).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch user bookings: %w", err)
	}

	return bookings, nil
}

// CancelBooking cancels a booking, unlocks court slots, and restores product inventory
func (s *BookingService) CancelBooking(bookingID string, reason string) (*models.Booking, error) {
	uid, err := uuid.Parse(bookingID)
	if err != nil {
		return nil, fmt.Errorf("invalid booking ID")
	}

	var cancelledBooking models.Booking

	err = s.db.Transaction(func(tx *gorm.DB) error {
		var booking models.Booking
		if err := tx.Preload("Courts").Preload("Items").First(&booking, "id = ?", uid).Error; err != nil {
			return fmt.Errorf("booking not found: %w", err)
		}

		if booking.Status == models.BookingStatusCancelled {
			return fmt.Errorf("this booking has already been cancelled")
		}
		if booking.Status == models.BookingStatusExpired {
			return fmt.Errorf("this booking has already expired and cannot be cancelled")
		}

		// 1. Update Booking Status to CANCELLED
		booking.Status = models.BookingStatusCancelled
		if err := tx.Save(&booking).Error; err != nil {
			return fmt.Errorf("failed to update booking status: %w", err)
		}

		// 2. Unlock court slots (is_slot_locked = FALSE)
		if err := tx.Model(&models.BookingCourt{}).
			Where("booking_id = ?", booking.ID).
			Update("is_slot_locked", false).Error; err != nil {
			return fmt.Errorf("failed to release court slot: %w", err)
		}

		// 3. Restore inventory stock for each item
		for _, item := range booking.Items {
			if err := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				UpdateColumn("stock_qty", gorm.Expr("stock_qty + ?", item.Qty)).Error; err != nil {
				return fmt.Errorf("failed to restore product inventory: %w", err)
			}
		}

		// Reload with relations
		tx.Preload("User").Preload("Courts.Court").Preload("Items.Product").Preload("PaymentSlip").First(&cancelledBooking, "id = ?", booking.ID)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Broadcast real-time event to Admin and Demo viewers
	Hub.Broadcast("BOOKING_CANCELLED", map[string]interface{}{
		"booking_id": cancelledBooking.ID,
		"booking_no": cancelledBooking.BookingNo,
		"reason":     reason,
	})

	return &cancelledBooking, nil
}


