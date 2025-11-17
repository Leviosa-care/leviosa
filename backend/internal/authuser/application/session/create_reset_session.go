package session

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/hengadev/encx"
)

func (s *SessionService) CreateResetSession(ctx context.Context, token, userEmail string, ttl time.Duration) error {
	// Hash the token for storage
	tokenBytes, err := encx.SerializeValue(token)
	if err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}
	tokenHash := s.crypto.HashBasic(ctx, tokenBytes)

	// Store plaintext email directly (no hashing needed)
	if err := s.repo.StoreResetSession(ctx, tokenHash, userEmail, ttl); err != nil {
		return fmt.Errorf("store reset session: %w", err)
	}

	return nil
}
