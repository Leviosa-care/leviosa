package auth

import (
	"encoding/base64"

	"github.com/hengadev/errsx"
)

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