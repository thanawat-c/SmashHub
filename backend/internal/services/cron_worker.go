package services

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"

	"smashhub/internal/models"
)

type CronWorker struct {
	db     *gorm.DB
	ticker *time.Ticker
	stopCh chan struct{}
}

func NewCronWorker(db *gorm.DB) *CronWorker {
	return &CronWorker{
		db:     db,
		stopCh: make(chan struct{}),
	}
}

// Start begins background polling every 30 seconds
func (w *CronWorker) Start(ctx context.Context) {
	w.ticker = time.NewTicker(30 * time.Second)
	log.Println("[CronWorker] Started booking expiration worker (polling every 30s)")

	// Run once immediately on startup
	w.ExpireOverdueBookings()

	go func() {
		for {
			select {
			case <-w.ticker.C:
				w.ExpireOverdueBookings()
			case <-w.stopCh:
				w.ticker.Stop()
				log.Println("[CronWorker] Stopped")
				return
			case <-ctx.Done():
				w.ticker.Stop()
				log.Println("[CronWorker] Context cancelled, worker stopped")
				return
			}
		}
	}()
}

// Stop stops the background ticker
func (w *CronWorker) Stop() {
	close(w.stopCh)
}

// ExpireOverdueBookings scans for expired PENDING_PAYMENT bookings, unlocks court slots, and restores product inventory
func (w *CronWorker) ExpireOverdueBookings() (int, error) {
	now := time.Now()
	var overdueBookings []models.Booking

	// Find overdue bookings
	err := w.db.Preload("Courts").Preload("Items").
		Where("status = ? AND expires_at < ?", models.BookingStatusPendingPayment, now).
		Find(&overdueBookings).Error

	if err != nil {
		log.Printf("[CronWorker] Error querying overdue bookings: %v", err)
		return 0, err
	}

	if len(overdueBookings) == 0 {
		return 0, nil
	}

	log.Printf("[CronWorker] Found %d overdue bookings to expire", len(overdueBookings))
	expiredCount := 0

	for _, booking := range overdueBookings {
		txErr := w.db.Transaction(func(tx *gorm.DB) error {
			// 1. Mark booking as EXPIRED
			if err := tx.Model(&models.Booking{}).Where("id = ?", booking.ID).
				Updates(map[string]interface{}{
					"status":     models.BookingStatusExpired,
					"updated_at": time.Now(),
				}).Error; err != nil {
				return err
			}

			// 2. Unlock court slots (is_slot_locked = FALSE)
			// This frees the slots for other customers!
			if err := tx.Model(&models.BookingCourt{}).Where("booking_id = ?", booking.ID).
				Update("is_slot_locked", false).Error; err != nil {
				return err
			}

			// 3. Restock inventory for items in this booking
			for _, item := range booking.Items {
				if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
					Update("stock_qty", gorm.Expr("stock_qty + ?", item.Qty)).Error; err != nil {
					return err
				}
			}

			return nil
		})

		if txErr != nil {
			log.Printf("[CronWorker] Failed to expire booking %s: %v", booking.BookingNo, txErr)
			continue
		}

		expiredCount++
		log.Printf("[CronWorker] Successfully expired booking %s, unlocked court slots, and restored stock", booking.BookingNo)

		// Broadcast SSE event
		Hub.Broadcast("BOOKING_EXPIRED", map[string]interface{}{
			"booking_id": booking.ID,
			"booking_no": booking.BookingNo,
			"status":     models.BookingStatusExpired,
		})
	}

	return expiredCount, nil
}
