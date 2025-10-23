package promotionCodeRepository

import (
	"fmt"
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *PromotionCodeRepository) GetPromotionCodeByID(ctx context.Context, promotionCodeID uuid.UUID) (*domain.PromotionCode, error) {
	query := fmt.Sprintf(`
		SELECT id, stripe_promotion_id, coupon_id, code, active, 
		       max_redemptions, times_redeemed, expires_at, first_time_transaction,
		       minimum_amount, minimum_amount_currency, restrictions, 
		       metadata, created_at, updated_at
		FROM %s.promotion_codes 
		WHERE id = $1
	`, r.schema)

	var promotionCode domain.PromotionCode
	row := r.pool.QueryRow(ctx, query, promotionCodeID)

	err := row.Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewRepositoryNotFoundErr(err, "promotion code")
		}
		return nil, errs.ClassifyPgError("get promotion code by ID", err)
	}

	return &promotionCode, nil
}