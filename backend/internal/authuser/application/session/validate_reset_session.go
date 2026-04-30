package session

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/hengadev/encx"
)

func (s *SessionService) ValidateResetSession(ctx context.Context, token string) (string, error) {
	// Hash the token for lookup
	tokenBytes, err := encx.SerializeValue(token)
	if err != nil {
		return "", errs.NewInvalidValueErr(err.Error())
	}
	tokenHash := s.crypto.HashBasic(ctx, tokenBytes)

	userEmail, err := s.repo.ValidateResetSession(ctx, tokenHash)
	if err != nil {
		return "", fmt.Errorf("validate reset session: %w", err)
	}

	return userEmail, nil
}
