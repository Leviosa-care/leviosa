package promotionCodeRepository

import (
	"fmt"
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *PromotionCodeRepository) DeletePromotionCodesByCouponID(ctx context.Context, couponID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.promotion_codes 
		WHERE coupon_id = $1
	`, r.schema)

	_, err := r.pool.Exec(ctx, query, couponID)
	if err != nil {
		return errs.ClassifyPgError("delete promotion codes by coupon ID", err)
	}

	// Note: We don't check RowsAffected here because it's valid to delete 0 rows
	// (a coupon might not have any promotion codes)

	return nil
}