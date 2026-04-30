package session

import (
	"encoding/base64"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/hengadev/errsx"
)

func ValidateToken(token string) error {
	var errs errsx.Map

	if token == "" {
		errs.Set("token missing", "token is required")
	}
	decoded, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		errs.Set("invalid format", "token must be a valid base64 string")
	}
	if err == nil && len(decoded) != cookies.TokenLength {
		errs.Set("token invalid", "token has invalid length")
	}

	return errs.AsError()
}