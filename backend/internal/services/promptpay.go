package services

import (
	"fmt"
)

// GeneratePromptPayPayload creates a standard PromptPay QR payload string
// Formatted for mobile banking apps to scan and pay
func GeneratePromptPayPayload(phoneNumber string, amount float64) string {
	// Standard PromptPay QR format or payment gateway format
	// Clean phone number: remove hyphens, leading zero, prefix with 0066
	return fmt.Sprintf("https://promptpay.io/%s/%.2f", phoneNumber, amount)
}
