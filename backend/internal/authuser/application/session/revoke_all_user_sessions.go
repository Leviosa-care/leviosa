package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

func (s *SessionService) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	userIDBytes, err := encx.SerializeValue(userID)
	if err != nil {
		return errs.NewInvalidValueErr(fmt.Sprintf("failed to serialize userID: %v", err))
	}
	userIDHash := s.crypto.HashBasic(ctx, userIDBytes)

	// Call repository with the hashed userID
	if err := s.repo.RevokeAllUserSessions(ctx, userIDHash); err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			// No sessions found for user - this is actually success for revocation
			return nil
		}
		return fmt.Errorf("revoke all user sessions: %w", err)
	}

	return nil
}
