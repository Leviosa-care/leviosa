package promotionCode

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

func (s *PromotionCodeService) GetAllPromotionCodes(ctx context.Context) ([]*domain.PromotionCode, error) {
	promotionCodes, err := s.repo.GetAllPromotionCodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all promotion codes: %w", err)
	}

	return promotionCodes, nil
}
