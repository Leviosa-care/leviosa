package promotionCode

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

func (s *PromotionCodeService) GetActivePromotionCodes(ctx context.Context) ([]*domain.PromotionCode, error) {
	promotionCodes, err := s.repo.GetActivePromotionCodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("get active promotion codes: %w", err)
	}

	return promotionCodes, nil
}
