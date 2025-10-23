package promotionCodeRepository

import (
	"fmt"
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *PromotionCodeRepository) GetPromotionCodesByCouponID(ctx context.Context, couponID uuid.UUID) ([]*domain.PromotionCode, error) {
	query := fmt.Sprintf(`
		SELECT id, stripe_promotion_id, coupon_id, code, active, 
		       max_redemptions, times_redeemed, expires_at, first_time_transaction,
		       minimum_amount, minimum_amount_currency, restrictions, 
		       metadata, created_at, updated_at
		FROM %s.promotion_codes 
		WHERE coupon_id = $1
		ORDER BY created_at DESC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, couponID)
	if err != nil {
		return nil, errs.ClassifyPgError("get promotion codes by coupon ID", err)
	}
	defer rows.Close()

	var promotionCodes []*domain.PromotionCode
	for rows.Next() {
		var promotionCode domain.PromotionCode

		err := rows.Scan(
			&promotionCode.ID,
			&promotionCode.StripePromotionID,
			&promotionCode.CouponID,
			&promotionCode.Code,
			&promotionCode.Active,
			&promotionCode.MaxRedemptions,
			&promotionCode.TimesRedeemed,
			&promotionCode.ExpiresAt,
			&promotionCode.FirstTimeTransaction,
			&promotionCode.MinimumAmount,
			&promotionCode.MinimumAmountCurrency,
			&promotionCode.Restrictions,
			&promotionCode.Metadata,
			&promotionCode.CreatedAt,
			&promotionCode.UpdatedAt,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan promotion code", err)
		}

		promotionCodes = append(promotionCodes, &promotionCode)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate over promotion codes", err)
	}

	return promotionCodes, nil
}