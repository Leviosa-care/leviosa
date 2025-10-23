package couponRepository

import (
	"fmt"
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/jackc/pgx/v5"
)

func (r *CouponRepository) GetCouponByStripeID(ctx context.Context, stripeCouponID string) (*domain.Coupon, error) {
	query := fmt.Sprintf(`
		SELECT id, stripe_coupon_id, name, percent_off, amount_off, currency, 
		       duration, duration_in_months, max_redemptions, times_redeemed, 
		       valid, redeem_by, metadata, created_at, updated_at
		FROM %s.coupons 
		WHERE stripe_coupon_id = $1
	`, r.schema)

	var coupon domain.Coupon
	row := r.pool.QueryRow(ctx, query, stripeCouponID)

	err := row.Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewRepositoryNotFoundErr(err, "coupon")
		}
		return nil, errs.ClassifyPgError("get coupon by Stripe ID", err)
	}

	return &coupon, nil
}