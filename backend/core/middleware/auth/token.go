package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/hengadev/errsx"
)

const TokenLength = 32

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

func ValidateToken(token string) error {
	var errs errsx.Map

	if token == "" {
		errs.Set("token missing", "token is required")
	}
	if _, err := base64.URLEncoding.DecodeString(token); err != nil {
		errs.Set("invalid format", "token must be a valid base64 string")
	}
	decoded, err := base64.URLEncoding.DecodeString(token)
	if err == nil && len(decoded) != TokenLength {
		errs.Set("token invalid", "token has invalid length")
	}

	return errs.AsError()
}