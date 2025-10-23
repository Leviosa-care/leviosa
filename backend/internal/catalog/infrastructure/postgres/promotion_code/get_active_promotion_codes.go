package promotionCodeRepository

import (
	"fmt"
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *PromotionCodeRepository) GetActivePromotionCodes(ctx context.Context) ([]*domain.PromotionCode, error) {
	query := fmt.Sprintf(`
		SELECT id, stripe_promotion_id, coupon_id, code, active, 
		       max_redemptions, times_redeemed, expires_at, first_time_transaction,
		       minimum_amount, minimum_amount_currency, restrictions, 
		       metadata, created_at, updated_at
		FROM %s.promotion_codes 
		WHERE active = true 
		  AND (expires_at IS NULL OR expires_at > NOW())
		  AND (max_redemptions IS NULL OR times_redeemed < max_redemptions)
		ORDER BY created_at DESC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("get active promotion codes", err)
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
			return nil, errs.ClassifyPgError("scan active promotion code", err)
		}

		promotionCodes = append(promotionCodes, &promotionCode)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate over active promotion codes", err)
	}

	return promotionCodes, nil
}