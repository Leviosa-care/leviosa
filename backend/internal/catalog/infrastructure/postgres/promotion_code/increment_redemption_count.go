package promotionCodeRepository

import (
	"fmt"
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *PromotionCodeRepository) IncrementRedemptionCount(ctx context.Context, promotionCodeID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.promotion_codes 
		SET times_redeemed = times_redeemed + 1,
		    updated_at = $2
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, promotionCodeID, time.Now())
	if err != nil {
		return errs.ClassifyPgError("increment promotion code redemption count", err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "promotion code")
	}

	return nil
}