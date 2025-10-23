package couponRepository

import (
	"fmt"
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *CouponRepository) CouponExistsByName(ctx context.Context, name string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(SELECT 1 FROM %s.coupons WHERE name = $1)
	`, r.schema)

	var exists bool
	err := r.pool.QueryRow(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, errs.ClassifyPgError("check coupon exists by name", err)
	}

	return exists, nil
}