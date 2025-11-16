package promotionCode

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *PromotionCodeService) GetPromotionCodeByStripeID(ctx context.Context, stripePromotionID string) (*domain.PromotionCode, error) {
	if stripePromotionID == "" {
		return nil, errs.NewInvalidValueErr("stripe promotion ID cannot be empty")
	}

	promotionCode, err := s.repo.GetPromotionCodeByStripeID(ctx, stripePromotionID)
	if err != nil {
		return nil, fmt.Errorf("get promotion code by Stripe ID %s: %w", stripePromotionID, err)
	}

	return promotionCode, nil
}
