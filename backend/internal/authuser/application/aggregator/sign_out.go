package aggregator

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
)

func (s *AuthAggregatorService) SignOut(ctx context.Context, sessionInfo *session.SessionInfo) error {
	// Remove the session using the session ID from the SessionInfo
	if err := s.session.RemoveSession(ctx, sessionInfo.ID); err != nil {
		return fmt.Errorf("failed to sign out user: %w", err)
	}

	return nil
}
