package price

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *PriceService) GetAllPrices(ctx context.Context) ([]*domain.Price, error) {
	prices, err := s.repo.GetAllPrices(ctx)
	if err != nil {
		return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed to get all prices: %w", err))
	}
	return prices, nil
}
