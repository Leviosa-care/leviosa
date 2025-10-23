package couponRepository

import (
	"fmt"
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *CouponRepository) CreateCoupon(ctx context.Context, coupon *domain.Coupon) (string, error) {
	// Generate UUID if not provided
	if coupon.ID == uuid.Nil {
		coupon.ID = uuid.New()
	}

	query := fmt.Sprintf(`
		INSERT INTO %s.coupons (
			id, stripe_coupon_id, name, percent_off, amount_off, currency, 
			duration, duration_in_months, max_redemptions, times_redeemed, 
			valid, redeem_by, metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		) RETURNING id
	`, r.schema)

	var createdID uuid.UUID
	err := r.pool.QueryRow(ctx, query,
		coupon.ID,
		coupon.StripeCouponID,
		coupon.Name,
		coupon.PercentOff,
		coupon.AmountOff,
		coupon.Currency,
		coupon.Duration,
		coupon.DurationInMonths,
		coupon.MaxRedemptions,
		coupon.TimesRedeemed,
		coupon.IsValid,
		coupon.RedeemBy,
		coupon.Metadata,
		coupon.CreatedAt,
		coupon.UpdatedAt,
	).Scan(&createdID)

	if err != nil {
		return "", errs.ClassifyPgError("create coupon", err)
	}

	return createdID.String(), nil
}