package couponRepository

import (
	"fmt"
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *CouponRepository) GetAllCoupons(ctx context.Context) ([]*domain.Coupon, error) {
	query := fmt.Sprintf(`
		SELECT id, stripe_coupon_id, name, percent_off, amount_off, currency, 
		       duration, duration_in_months, max_redemptions, times_redeemed, 
		       valid, redeem_by, metadata, created_at, updated_at
		FROM %s.coupons 
		ORDER BY created_at DESC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("get all coupons", err)
	}
	defer rows.Close()

	var coupons []*domain.Coupon
	for rows.Next() {
		var coupon domain.Coupon
		err := rows.Scan(
			&coupon.ID,
			&coupon.StripeCouponID,
			&coupon.Name,
			&coupon.PercentOff,
			&coupon.AmountOff,
			&coupon.Currency,
			&coupon.Duration,
			&coupon.DurationInMonths,
			&coupon.MaxRedemptions,
			&coupon.TimesRedeemed,
			&coupon.IsValid,
			&coupon.RedeemBy,
			&coupon.Metadata,
			&coupon.CreatedAt,
			&coupon.UpdatedAt,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan coupon", err)
		}
		coupons = append(coupons, &coupon)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate over coupons", err)
	}

	return coupons, nil
}