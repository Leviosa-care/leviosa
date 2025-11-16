package promotionCode

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *PromotionCodeService) GetPromotionCodeByID(ctx context.Context, promotionCodeID string) (*domain.PromotionCode, error) {
	id, err := uuid.Parse(promotionCodeID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("invalid promotion code ID format")
	}

	promotionCode, err := s.repo.GetPromotionCodeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get promotion code by ID: %w", err)
	}

	return promotionCode, nil
}
