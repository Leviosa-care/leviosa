package promotionCode

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *PromotionCodeService) DeactivatePromotionCode(ctx context.Context, promotionCodeID string) error {
	id, err := uuid.Parse(promotionCodeID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid promotion code ID format")
	}

	// Check if promotion code exists
	_, err = s.repo.GetPromotionCodeByID(ctx, id)
	if err != nil {
		return fmt.Errorf("validate promotion code: %w", err)
	}

	if err := s.repo.DeactivatePromotionCode(ctx, id); err != nil {
		return fmt.Errorf("deactivate promotion code: %w", err)
	}

	return nil
}
