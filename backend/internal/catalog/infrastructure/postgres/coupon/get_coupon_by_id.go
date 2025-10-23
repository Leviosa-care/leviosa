package couponRepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *CouponRepository) GetCouponByID(ctx context.Context, couponID uuid.UUID) (*domain.Coupon, error) {
	query := fmt.Sprintf(`
		SELECT id, stripe_coupon_id, name, percent_off, amount_off, currency, 
		       duration, duration_in_months, max_redemptions, times_redeemed, 
		       valid, redeem_by, metadata, created_at, updated_at
		FROM %s.coupons 
		WHERE id = $1
	`, r.schema)

	var coupon domain.Coupon
	row := r.pool.QueryRow(ctx, query, couponID)

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
		return nil, errs.ClassifyPgError("get coupon by ID", err)
	}

	return &coupon, nil
}
