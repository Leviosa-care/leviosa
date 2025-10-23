package promotionCodeRepository

import (
	"fmt"
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *PromotionCodeRepository) CreatePromotionCode(ctx context.Context, promotionCode *domain.PromotionCode) (string, error) {
	// Generate UUID if not provided
	if promotionCode.ID == uuid.Nil {
		promotionCode.ID = uuid.New()
	}

	// Extract currency options from restrictions
	query := fmt.Sprintf(`
		INSERT INTO %s.promotion_codes (
			id, stripe_promotion_id, coupon_id, code, active, 
			max_redemptions, times_redeemed, expires_at, first_time_transaction,
			minimum_amount, minimum_amount_currency, restrictions, 
			metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		) RETURNING id
	`, r.schema)

	var createdID uuid.UUID
	err := r.pool.QueryRow(ctx, query,
		promotionCode.ID,
		promotionCode.StripePromotionID,
		promotionCode.CouponID,
		promotionCode.Code,
		promotionCode.Active,
		promotionCode.MaxRedemptions,
		promotionCode.TimesRedeemed,
		promotionCode.ExpiresAt,
		promotionCode.FirstTimeTransaction,
		promotionCode.MinimumAmount,
		promotionCode.MinimumAmountCurrency,
		promotionCode.Restrictions,
		promotionCode.Metadata,
		promotionCode.CreatedAt,
		promotionCode.UpdatedAt,
	).Scan(&createdID)

	if err != nil {
		return "", errs.ClassifyPgError("create promotion code", err)
	}

	return createdID.String(), nil
}
