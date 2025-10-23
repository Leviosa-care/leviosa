package couponRepository

import (
	"fmt"
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *CouponRepository) DeactivateCoupon(ctx context.Context, couponID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.coupons 
		SET valid = false,
		    updated_at = $2
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, couponID, time.Now())
	if err != nil {
		return errs.ClassifyPgError("deactivate coupon", err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "coupon")
	}

	return nil
}