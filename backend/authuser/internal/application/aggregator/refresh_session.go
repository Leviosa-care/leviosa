package aggregator

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
)

func (s *AuthAggregatorService) RefreshSession(ctx context.Context, request *domain.RefreshSessionRequest) (*domain.RefreshSessionResponse, error) {
	return s.session.RefreshSession(ctx, request)
}