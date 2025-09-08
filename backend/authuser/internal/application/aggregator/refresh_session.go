package aggregator

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
)

func (s *AuthAggregatorService) RefreshSession(ctx context.Context, sessionID uuid.UUID) (*domain.RefreshSessionResponse, error) {
	return s.session.RefreshSession(ctx, sessionID)
}
