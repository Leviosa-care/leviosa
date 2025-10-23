package promotionCodeRepository

import (
	"fmt"
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *PromotionCodeRepository) DeletePromotionCode(ctx context.Context, promotionCodeID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.promotion_codes 
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, promotionCodeID)
	if err != nil {
		return errs.ClassifyPgError("delete promotion code", err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "promotion code")
	}

	return nil
}