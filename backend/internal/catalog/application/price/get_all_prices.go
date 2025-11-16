package price

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

func (s *PriceService) GetAllPrices(ctx context.Context) ([]*domain.Price, error) {
	prices, err := s.repo.GetAllPrices(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all prices: %w", err)
	}
	return prices, nil
}
