package domain

import (
	"fmt"
	"regexp"
	"strings"
)

type SMSRequest struct {
	Phone   string
	Message string
}

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

func normalizePhoneNumber(phone string) (string, error) {
	phone = strings.TrimSpace(phone)
	
	if phone == "" {
		return "", fmt.Errorf("phone number cannot be empty")
	}

	re := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !re.MatchString(phone) {
		return "", fmt.Errorf("phone number format is invalid")
	}

	if !strings.HasPrefix(phone, "+") {
		phone = "+" + phone
	}

	return phone, nil
}