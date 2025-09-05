package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateToken generates a secure random session ID
func GenerateToken() (string, error) {
	// length = number of raw bytes, before encoding
	b := make([]byte, TokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}
	// Base64 encode to make it URL-safe (can also use hex encoding)
	return base64.URLEncoding.EncodeToString(b), nil
}
