package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"smashhub/internal/models"
)

var (
	ErrDuplicateTransRef = errors.New("สลิปนี้เคยถูกใช้งานไปแล้ว (รหัสอ้างอิงซ้ำ)")
	ErrUnderpaid         = errors.New("ยอดเงินในสลิปไม่ตรงกับยอดที่ต้องชำระ (โอนขาด)")
	ErrWrongReceiver     = errors.New("บัญชีปลายทางไม่ถูกต้อง กรุณาโอนเข้าบัญชีร้านสแมชฮับ")
	ErrExternalAPIDown   = errors.New("ระบบตรวจสลิปปลายทางขัดข้อง กำลังส่งให้แอดมินตรวจสอบ")
)

type VerifySlipRequest struct {
	MockScenario string  `json:"mock_scenario" form:"mock_scenario"` // MOCK_VALID, MOCK_DUPLICATE, MOCK_UNDERPAID, MOCK_API_DOWN
	TransRef     string  `json:"trans_ref" form:"trans_ref"`
	Amount       float64 `json:"amount" form:"amount"`
	SlipImageURL string  `json:"slip_image_url" form:"slip_image_url"`
}

type SlipVerificationResult struct {
	Success            bool                `json:"success"`
	BookingStatus      models.BookingStatus `json:"booking_status"`
	VerificationStatus models.SlipStatus   `json:"verification_status"`
	Message            string              `json:"message"`
	Slip               *models.PaymentSlip `json:"slip,omitempty"`
	Booking            *models.Booking     `json:"booking,omitempty"`
}

type SlipService struct {
	db *gorm.DB
}

func NewSlipService(db *gorm.DB) *SlipService {
	return &SlipService{db: db}
}

// VerifySlip executes the 4-step verification flow with portfolio mock scenarios
func (s *SlipService) VerifySlip(bookingID uuid.UUID, req VerifySlipRequest) (*SlipVerificationResult, error) {
	var booking models.Booking
	if err := s.db.Preload("User").Preload("Courts.Court").Preload("Items.Product").First(&booking, "id = ?", bookingID).Error; err != nil {
		return nil, ErrBookingNotFound
	}

	if booking.Status == models.BookingStatusConfirmed {
		return &SlipVerificationResult{
			Success:            true,
			BookingStatus:      booking.Status,
			VerificationStatus: models.SlipStatusVerified,
			Message:            "การจองนี้ได้รับการยืนยันการชำระเงินเรียบร้อยแล้ว",
			Booking:            &booking,
		}, nil
	}

	if booking.Status == models.BookingStatusExpired || booking.Status == models.BookingStatusCancelled {
		return nil, errors.New("ใบจองนี้หมดอายุหรือถูกยกเลิกแล้ว ไม่สามารถแนบสลิปได้")
	}

	// 1. Prepare Slip Data & handle Mock Scenarios
	transRef := req.TransRef
	amount := req.Amount
	receiverAccount := "098-7-65432-1 (บจก. สแมชฮับ แบดมินตัน)"
	senderName := "ลูกค้าสแมชฮับ"
	if booking.User != nil {
		senderName = booking.User.DisplayName
	}
	slipURL := req.SlipImageURL
	if slipURL == "" {
		slipURL = "https://placehold.co/600x800/004d40/ffffff?text=Uploaded+Slip+" + booking.BookingNo
	}

	// Mock mode parameters
	switch req.MockScenario {
	case "MOCK_DUPLICATE":
		// Use known existing trans_ref from seed data
		transRef = "KBANK-20261024-9876543210"
		amount = booking.NetTotal

	case "MOCK_UNDERPAID":
		transRef = fmt.Sprintf("MOCK-TXN-%d", time.Now().UnixNano()%10000000)
		amount = booking.NetTotal - 100
		if amount <= 0 {
			amount = 1.00
		}

	case "MOCK_API_DOWN":
		transRef = fmt.Sprintf("MOCK-TXN-%d", time.Now().UnixNano()%10000000)
		amount = booking.NetTotal

	case "MOCK_VALID", "":
		if transRef == "" {
			transRef = fmt.Sprintf("KBANK-%s-%06d", time.Now().Format("20060102"), time.Now().UnixNano()%1000000)
		}
		if amount <= 0 {
			amount = booking.NetTotal
		}
	}

	// Step 1: Duplicate check (trans_ref)
	var existingSlip models.PaymentSlip
	if err := s.db.Where("trans_ref = ?", transRef).First(&existingSlip).Error; err == nil {
		// Duplicate found!
		rawPayload, _ := json.Marshal(map[string]interface{}{
			"error":        "Duplicate transaction reference found in database",
			"trans_ref":    transRef,
			"duplicate_of": existingSlip.ID,
		})

		rejectedSlip := models.PaymentSlip{
			BookingID:          booking.ID,
			SlipImageURL:       slipURL,
			TransRef:           &transRef,
			SenderName:         senderName,
			ReceiverAccount:    receiverAccount,
			Amount:             amount,
			VerificationStatus: models.SlipStatusRejectedDuplicate,
			RawPayload:         datatypes.JSON(rawPayload),
		}
		s.db.Create(&rejectedSlip)

		Hub.Broadcast("SLIP_REJECTED", map[string]interface{}{
			"booking_id": booking.ID,
			"reason":     "REJECTED_DUPLICATE",
		})

		return &SlipVerificationResult{
			Success:            false,
			BookingStatus:      booking.Status,
			VerificationStatus: models.SlipStatusRejectedDuplicate,
			Message:            "ตรวจสอบไม่ผ่าน: สลิปนี้เคยถูกใช้งานไปแล้ว (รหัสอ้างอิงซ้ำ)",
			Slip:               &rejectedSlip,
			Booking:            &booking,
		}, ErrDuplicateTransRef
	}

	// Step 2: Receiver account check
	// (In production, checks receiver promptpay ID / company tax ID)

	// Step 3: Check amount >= booking.net_total
	if amount < booking.NetTotal {
		rawPayload, _ := json.Marshal(map[string]interface{}{
			"error":            "Amount mismatch",
			"expected_amount":  booking.NetTotal,
			"transferred_amount": amount,
			"difference":       booking.NetTotal - amount,
		})

		rejectedSlip := models.PaymentSlip{
			BookingID:          booking.ID,
			SlipImageURL:       slipURL,
			TransRef:           &transRef,
			SenderName:         senderName,
			ReceiverAccount:    receiverAccount,
			Amount:             amount,
			VerificationStatus: models.SlipStatusRejectedAmountMismatch,
			RawPayload:         datatypes.JSON(rawPayload),
		}
		s.db.Create(&rejectedSlip)

		Hub.Broadcast("SLIP_REJECTED", map[string]interface{}{
			"booking_id": booking.ID,
			"reason":     "REJECTED_AMOUNT_MISMATCH",
		})

		return &SlipVerificationResult{
			Success:            false,
			BookingStatus:      booking.Status,
			VerificationStatus: models.SlipStatusRejectedAmountMismatch,
			Message:            fmt.Sprintf("ตรวจสอบไม่ผ่าน: ยอดเงินในสลิป (%.2f บาท) น้อยกว่ายอดที่ต้องชำระ (%.2f บาท)", amount, booking.NetTotal),
			Slip:               &rejectedSlip,
			Booking:            &booking,
		}, ErrUnderpaid
	}

	// Step 4: External API Gateway check
	if req.MockScenario == "MOCK_API_DOWN" {
		rawPayload, _ := json.Marshal(map[string]interface{}{
			"error":       "OpenSlipVerify Gateway Timeout (HTTP 504)",
			"description": "ธนาคารปลายทางตอบสนองช้ากว่ากำหนด ส่งเข้าคิวตรวจโดยเจ้าหน้าที่",
			"timestamp":   time.Now().Format(time.RFC3339),
		})

		slip := models.PaymentSlip{
			BookingID:          booking.ID,
			SlipImageURL:       slipURL,
			TransRef:           &transRef,
			SenderName:         senderName,
			ReceiverAccount:    receiverAccount,
			Amount:             amount,
			VerificationStatus: models.SlipStatusAPIErrorPendingAdmin,
			RawPayload:         datatypes.JSON(rawPayload),
		}
		s.db.Create(&slip)

		// Transition booking to MANUAL_REVIEW
		booking.Status = models.BookingStatusManualReview
		s.db.Save(&booking)

		Hub.Broadcast("BOOKING_MANUAL_REVIEW", map[string]interface{}{
			"booking_id": booking.ID,
			"booking_no": booking.BookingNo,
			"customer":   senderName,
			"amount":     amount,
		})

		return &SlipVerificationResult{
			Success:            true,
			BookingStatus:      models.BookingStatusManualReview,
			VerificationStatus: models.SlipStatusAPIErrorPendingAdmin,
			Message:            "ระบบตรวจสลิปอัตโนมัติขัดข้องชั่วคราว เจ้าหน้าที่กำลังตรวจสอบสลิปให้ท่านโดยเร็ว",
			Slip:               &slip,
			Booking:            &booking,
		}, nil
	}

	// Step 5: All checks passed! Mark as VERIFIED & CONFIRMED
	now := time.Now()
	rawPayload, _ := json.Marshal(map[string]interface{}{
		"status":      "SUCCESS",
		"verified_by": "OpenSlipVerify Gateway",
		"bank":        "KBANK",
		"timestamp":   now.Format(time.RFC3339),
	})

	verifiedSlip := models.PaymentSlip{
		BookingID:          booking.ID,
		SlipImageURL:       slipURL,
		TransRef:           &transRef,
		SenderName:         senderName,
		ReceiverAccount:    receiverAccount,
		Amount:             amount,
		VerificationStatus: models.SlipStatusVerified,
		RawPayload:         datatypes.JSON(rawPayload),
		VerifiedAt:         &now,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&verifiedSlip).Error; err != nil {
			return err
		}
		booking.Status = models.BookingStatusConfirmed
		return tx.Save(&booking).Error
	})

	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถบันทึกผลการยืนยันได้: %w", err)
	}

	Hub.Broadcast("BOOKING_CONFIRMED", map[string]interface{}{
		"booking_id": booking.ID,
		"booking_no": booking.BookingNo,
		"customer":   senderName,
		"amount":     amount,
	})

	return &SlipVerificationResult{
		Success:            true,
		BookingStatus:      models.BookingStatusConfirmed,
		VerificationStatus: models.SlipStatusVerified,
		Message:            "ยืนยันการชำระเงินเรียบร้อยแล้ว คอร์ทของท่านพร้อมใช้งาน!",
		Slip:               &verifiedSlip,
		Booking:            &booking,
	}, nil
}
