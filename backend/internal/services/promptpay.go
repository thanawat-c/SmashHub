package services

import (
	"fmt"
	"strings"
)

// crc16CCITT calculates the CRC-16/CCITT-FALSE checksum (polynomial 0x1021, initial 0xFFFF)
func crc16CCITT(data []byte) uint16 {
	var crc uint16 = 0xFFFF
	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc = crc << 1
			}
		}
	}
	return crc
}

// GeneratePromptPayPayload creates an official EMVCo Thai QR Payment payload for PromptPay
func GeneratePromptPayPayload(phoneNumber string, amount float64) string {
	// Sanitize phone number (e.g. "098-7-65432-1" -> "0987654321")
	cleaned := strings.ReplaceAll(phoneNumber, "-", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	// Prefix with country code 0066 without leading 0
	if strings.HasPrefix(cleaned, "0") {
		cleaned = "0066" + cleaned[1:]
	}

	phoneTag := fmt.Sprintf("01%02d%s", len(cleaned), cleaned)
	merchantInfo := "0016A000000677010111" + phoneTag
	tag29 := fmt.Sprintf("29%02d%s", len(merchantInfo), merchantInfo)

	tag54 := ""
	if amount > 0 {
		amtStr := fmt.Sprintf("%.2f", amount)
		tag54 = fmt.Sprintf("54%02d%s", len(amtStr), amtStr)
	}

	raw := fmt.Sprintf("000201010212%s5303764%s5802TH6304", tag29, tag54)
	checksum := fmt.Sprintf("%04X", crc16CCITT([]byte(raw)))
	return raw + checksum
}
