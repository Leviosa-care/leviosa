package domain

import (
	"fmt"
	"regexp"
	"strings"
)

// OTP SMS
type OTPSMSRequest struct {
	PhoneNumber string `json:"phone_number"`
	OTP         string `json:"otp"`
}

// Generic SMS
type GenericSMSRequest struct {
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
}

// SMSRequest is the internal SMS structure
type SMSRequest struct {
	Phone   string
	Message string
}

// NewSMSRequest creates a validated SMS request with normalized phone number
func NewSMSRequest(phone, message string) (*SMSRequest, error) {
	cleanPhone, err := normalizePhoneNumber(phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	return &SMSRequest{
		Phone:   cleanPhone,
		Message: strings.TrimSpace(message),
	}, nil
}

// normalizePhoneNumber validates and normalizes phone numbers to E.164 format
func normalizePhoneNumber(phone string) (string, error) {
	phone = strings.TrimSpace(phone)

	if phone == "" {
		return "", fmt.Errorf("phone number cannot be empty")
	}

	// E.164 format validation: +[country code][number] (max 15 digits)
	re := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !re.MatchString(phone) {
		return "", fmt.Errorf("phone number format is invalid")
	}

	// Add + prefix if missing
	if !strings.HasPrefix(phone, "+") {
		phone = "+" + phone
	}

	return phone, nil
}
