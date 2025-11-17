package aggregator

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
)

func (s *AuthAggregatorService) SignOut(ctx context.Context, sessionInfo *session.SessionInfo) error {
	return s.session.RemoveSession(ctx, sessionInfo.ID)
}
