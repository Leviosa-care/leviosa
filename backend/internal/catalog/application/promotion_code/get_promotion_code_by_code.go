package promotionCode

import (
	"context"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *PromotionCodeService) GetPromotionCodeByCode(ctx context.Context, code string) (*domain.PromotionCode, error) {
	if code == "" {
		return nil, errs.NewInvalidValueErr("promotion code cannot be empty")
	}

	code = strings.ToUpper(strings.TrimSpace(code))

	promotionCode, err := s.repo.GetPromotionCodeByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get promotion code by code %s: %w", code, err)
	}

	return promotionCode, nil
}
