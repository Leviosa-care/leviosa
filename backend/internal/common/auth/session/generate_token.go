package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
)

// GenerateToken generates a secure random session ID
func GenerateToken() (string, error) {
	// length = number of raw bytes, before encoding
	b := make([]byte, cookies.TokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}
	// Base64 encode to make it URL-safe (can also use hex encoding)
	return base64.URLEncoding.EncodeToString(b), nil
}
