package services

import (
	"fmt"
	"strings"
	"testing"
)

func TestCRC16CCITT(t *testing.T) {
	// Standard test vector for CRC-16/CCITT-FALSE
	// "123456789" with init 0xFFFF, poly 0x1021 -> 0x29B1
	input := []byte("123456789")
	expected := uint16(0x29B1)
	got := crc16CCITT(input)

	if got != expected {
		t.Errorf("crc16CCITT(\"123456789\") = %04X; want %04X", got, expected)
	}
}

func TestGeneratePromptPayPayload(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		amount      float64
		wantPrefix  string
		wantTag54   string
	}{
		{
			name:        "Standard 10-digit mobile with amount",
			phoneNumber: "0812345678",
			amount:      350.00,
			wantPrefix:  "00020101021229370016A000000677010111011300668123456785303764",
			wantTag54:   "5406350.00",
		},
		{
			name:        "Mobile with hyphens and spaces",
			phoneNumber: "081-234-5678",
			amount:      180.50,
			wantPrefix:  "00020101021229370016A000000677010111011300668123456785303764",
			wantTag54:   "5406180.50",
		},
		{
			name:        "Mobile with zero amount (dynamic merchant promptpay)",
			phoneNumber: "0987654321",
			amount:      0,
			wantPrefix:  "00020101021229370016A000000677010111011300669876543215303764",
			wantTag54:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload := GeneratePromptPayPayload(tc.phoneNumber, tc.amount)

			// 1. Must start with standard EMVCo prefix
			if !strings.HasPrefix(payload, tc.wantPrefix) {
				t.Errorf("GeneratePromptPayPayload prefix mismatch\ngot:  %s\nwant: %s...", payload, tc.wantPrefix)
			}

			// 2. Must contain tag 54 (amount) if amount > 0
			if tc.wantTag54 != "" && !strings.Contains(payload, tc.wantTag54) {
				t.Errorf("GeneratePromptPayPayload missing amount tag\ngot:  %s\nwant to contain: %s", payload, tc.wantTag54)
			}

			// 3. Must end with 6304 + 4-char CRC16 checksum
			if len(payload) < 8 {
				t.Fatalf("Payload too short: %s", payload)
			}

			rawWithoutChecksum := payload[:len(payload)-4]
			checksumInPayload := payload[len(payload)-4:]

			calculatedChecksum := fmt.Sprintf("%04X", crc16CCITT([]byte(rawWithoutChecksum)))
			if checksumInPayload != calculatedChecksum {
				t.Errorf("Checksum mismatch: got %s, calculated %s", checksumInPayload, calculatedChecksum)
			}
		})
	}
}
