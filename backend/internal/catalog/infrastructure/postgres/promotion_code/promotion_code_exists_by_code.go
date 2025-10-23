package promotionCodeRepository

import (
	"fmt"
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *PromotionCodeRepository) PromotionCodeExistsByCode(ctx context.Context, code string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(SELECT 1 FROM %s.promotion_codes WHERE code = $1)
	`, r.schema)

	var exists bool
	err := r.pool.QueryRow(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, errs.ClassifyPgError("check promotion code exists by code", err)
	}

	return exists, nil
}