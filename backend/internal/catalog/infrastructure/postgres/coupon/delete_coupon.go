package couponRepository

import (
	"fmt"
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *CouponRepository) DeleteCoupon(ctx context.Context, couponID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.coupons 
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, couponID)
	if err != nil {
		return errs.ClassifyPgError("delete coupon", err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "coupon")
	}

	return nil
}