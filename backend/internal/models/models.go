package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// BookingStatus enum type
type BookingStatus string

const (
	BookingStatusPendingPayment BookingStatus = "PENDING_PAYMENT"
	BookingStatusConfirmed      BookingStatus = "CONFIRMED"
	BookingStatusManualReview   BookingStatus = "MANUAL_REVIEW"
	BookingStatusExpired        BookingStatus = "EXPIRED"
	BookingStatusCancelled      BookingStatus = "CANCELLED"
)

// SlipStatus enum type
type SlipStatus string

const (
	SlipStatusVerified             SlipStatus = "VERIFIED"
	SlipStatusRejectedDuplicate    SlipStatus = "REJECTED_DUPLICATE"
	SlipStatusRejectedAmountMismatch SlipStatus = "REJECTED_AMOUNT_MISMATCH"
	SlipStatusRejectedWrongAccount SlipStatus = "REJECTED_WRONG_ACCOUNT"
	SlipStatusAPIErrorPendingAdmin SlipStatus = "API_ERROR_PENDING_ADMIN"
)

// User role
const (
	RoleCustomer = "CUSTOMER"
	RoleAdmin    = "ADMIN"
)

// Check-in status
const (
	CheckInStatusWaiting   = "WAITING"
	CheckInStatusCheckedIn = "CHECKED_IN"
	CheckInStatusCompleted = "COMPLETED"
)

// User represents customer or admin
type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	LineUserID  *string   `gorm:"type:varchar(100);unique" json:"line_user_id,omitempty"`
	DisplayName string    `gorm:"type:varchar(255);not null" json:"display_name"`
	PictureURL  string    `gorm:"type:text" json:"picture_url,omitempty"`
	Phone       string    `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Role        string    `gorm:"type:varchar(20);default:'CUSTOMER';not null" json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Court represents a badminton court
type Court struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourtNumber int       `gorm:"unique;not null" json:"court_number"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	CourtType   string    `gorm:"type:varchar(20);not null" json:"court_type"` // PVC, ACRYLIC
	NormalRate  float64   `gorm:"type:numeric(10,2);not null" json:"normal_rate"`
	PeakRate    float64   `gorm:"type:numeric(10,2);not null" json:"peak_rate"` // After 17:00
	Status      string    `gorm:"type:varchar(20);default:'ACTIVE';not null" json:"status"` // ACTIVE, MAINTENANCE
	CreatedAt   time.Time `json:"created_at"`
}

// Product represents mini store merchandise or rental item
type Product struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SKU       string    `gorm:"type:varchar(50);unique;not null" json:"sku"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Category  string    `gorm:"type:varchar(50);not null" json:"category"` // SHUTTLECOCK, DRINK, RENTAL
	Price     float64   `gorm:"type:numeric(10,2);not null" json:"price"`
	StockQty  int       `gorm:"not null" json:"stock_qty"`
	ImageURL  string    `gorm:"type:text" json:"image_url,omitempty"`
	IsActive  bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Booking represents a reservation transaction
type Booking struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BookingNo     string         `gorm:"type:varchar(50);unique;not null" json:"booking_no"`
	UserID        *uuid.UUID     `gorm:"type:uuid" json:"user_id,omitempty"`
	User          *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status        BookingStatus  `gorm:"type:booking_status;default:'PENDING_PAYMENT';not null" json:"status"`
	CourtTotal    float64        `gorm:"type:numeric(10,2);default:0.00;not null" json:"court_total"`
	ItemTotal     float64        `gorm:"type:numeric(10,2);default:0.00;not null" json:"item_total"`
	NetTotal      float64        `gorm:"type:numeric(10,2);default:0.00;not null" json:"net_total"`
	CheckInStatus string         `gorm:"type:varchar(20);default:'WAITING';not null" json:"check_in_status"`
	ExpiresAt     *time.Time     `json:"expires_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Courts        []BookingCourt `gorm:"foreignKey:BookingID" json:"courts,omitempty"`
	Items         []BookingItem  `gorm:"foreignKey:BookingID" json:"items,omitempty"`
	PaymentSlip   *PaymentSlip   `gorm:"foreignKey:BookingID" json:"payment_slip,omitempty"`
}

// BookingCourt represents a 1-hour court reservation slot
type BookingCourt struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BookingID    uuid.UUID `gorm:"type:uuid;not null;index" json:"booking_id"`
	CourtID      uint      `gorm:"not null" json:"court_id"`
	Court        *Court    `gorm:"foreignKey:CourtID" json:"court,omitempty"`
	PlayDate     time.Time `gorm:"type:date;not null" json:"play_date"`
	StartHour    int       `gorm:"not null" json:"start_hour"` // 10 - 21
	EndHour      int       `gorm:"not null" json:"end_hour"`   // 11 - 22
	HourlyRate   float64   `gorm:"type:numeric(10,2);not null" json:"hourly_rate"`
	IsSlotLocked bool      `gorm:"default:true;not null" json:"is_slot_locked"`
	CreatedAt    time.Time `json:"created_at"`
}

// BookingItem represents mini-store products purchased with a booking
type BookingItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BookingID uuid.UUID `gorm:"type:uuid;not null;index" json:"booking_id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Qty       int       `gorm:"not null" json:"qty"`
	UnitPrice float64   `gorm:"type:numeric(10,2);not null" json:"unit_price"`
	Subtotal  float64   `gorm:"type:numeric(10,2);not null" json:"subtotal"`
	CreatedAt time.Time `json:"created_at"`
}

// PaymentSlip represents the payment slip verification audit record
type PaymentSlip struct {
	ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BookingID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"booking_id"`
	SlipImageURL       string         `gorm:"type:text" json:"slip_image_url,omitempty"`
	TransRef           *string        `gorm:"type:varchar(100);unique" json:"trans_ref,omitempty"`
	SenderName         string         `gorm:"type:varchar(255)" json:"sender_name,omitempty"`
	ReceiverAccount    string         `gorm:"type:varchar(100)" json:"receiver_account,omitempty"`
	Amount             float64        `gorm:"type:numeric(10,2);not null" json:"amount"`
	VerificationStatus SlipStatus     `gorm:"type:slip_status;not null" json:"verification_status"`
	RawPayload         datatypes.JSON `gorm:"type:jsonb" json:"raw_payload,omitempty"`
	VerifiedAt         *time.Time     `json:"verified_at,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
}
