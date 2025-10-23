package promotionCodeRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPromotionCodeByStripeID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful retrieval by stripe ID", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Stripe ID Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("STRIPEID", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve promotion code by Stripe ID
		retrievedCode, err := repo.GetPromotionCodeByStripeID(ctx, promotionCode.StripePromotionID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCode)
		assert.Equal(t, promotionCode.ID, retrievedCode.ID)
		assert.Equal(t, promotionCode.Code, retrievedCode.Code)
		assert.Equal(t, promotionCode.CouponID, retrievedCode.CouponID)
		assert.Equal(t, promotionCode.StripePromotionID, retrievedCode.StripePromotionID)
		assert.Equal(t, promotionCode.Active, retrievedCode.Active)
	})

	t.Run("successful retrieval with all optional fields", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Complete Stripe Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with all fields
		promotionCode := helpers.NewValidPromotionCodeWithRestrictions("COMPLETE", coupon.ID, []string{"premium_user"})
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by Stripe ID
		retrievedCode, err := repo.GetPromotionCodeByStripeID(ctx, promotionCode.StripePromotionID)

		require.NoError(t, err)
		assert.Equal(t, []string{"premium_user"}, retrievedCode.Restrictions.CurrencyOptions)
		assert.Equal(t, promotionCode.Metadata, retrievedCode.Metadata)
	})

	t.Run("retrieval with non-existent stripe ID should return not found error", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		nonExistentStripeID := "promo_nonexistent123456789"

		retrievedCode, err := repo.GetPromotionCodeByStripeID(ctx, nonExistentStripeID)

		assert.Error(t, err)
		assert.Nil(t, retrievedCode)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("retrieval with empty stripe ID should return error", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		retrievedCode, err := repo.GetPromotionCodeByStripeID(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, retrievedCode)
	})

	t.Run("multiple promotion codes, retrieve specific one by stripe ID", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Multiple Stripe IDs")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert multiple promotion codes
		code1 := helpers.NewValidPromotionCode("STRIPE1", coupon.ID)
		code2 := helpers.NewValidPromotionCode("STRIPE2", coupon.ID)
		code3 := helpers.NewValidPromotionCode("STRIPE3", coupon.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)
		helpers.InsertPromotionCode(t, ctx, testPool, code3)

		// Retrieve specific code by Stripe ID
		retrievedCode, err := repo.GetPromotionCodeByStripeID(ctx, code2.StripePromotionID)

		require.NoError(t, err)
		assert.Equal(t, code2.ID, retrievedCode.ID)
		assert.Equal(t, "STRIPE2", retrievedCode.Code)
		assert.Equal(t, code2.StripePromotionID, retrievedCode.StripePromotionID)
	})

	t.Run("retrieval of expired promotion code by stripe ID", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Expired Stripe Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert expired promotion code
		promotionCode := helpers.NewExpiredPromotionCode("EXPIREDSTRIPE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by Stripe ID
		retrievedCode, err := repo.GetPromotionCodeByStripeID(ctx, promotionCode.StripePromotionID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCode.ExpiresAt)
		assert.Equal(t, promotionCode.StripePromotionID, retrievedCode.StripePromotionID)
	})

	t.Run("retrieval of inactive promotion code by stripe ID", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Inactive Stripe Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert inactive promotion code
		promotionCode := helpers.NewInactivePromotionCode("INACTIVESTRIPE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by Stripe ID
		retrievedCode, err := repo.GetPromotionCodeByStripeID(ctx, promotionCode.StripePromotionID)

		require.NoError(t, err)
		assert.False(t, retrievedCode.Active)
		assert.Equal(t, promotionCode.StripePromotionID, retrievedCode.StripePromotionID)
	})

	t.Run("retrieval with special characters in stripe ID", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Special Stripe Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with special Stripe ID
		promotionCode := helpers.NewValidPromotionCode("SPECIAL", coupon.ID)
		promotionCode.StripePromotionID = "promo_test-123_special.id"
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by Stripe ID
		retrievedCode, err := repo.GetPromotionCodeByStripeID(ctx, "promo_test-123_special.id")

		require.NoError(t, err)
		assert.Equal(t, "promo_test-123_special.id", retrievedCode.StripePromotionID)
	})

	t.Run("retrieval with redemption limits", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Limited Stripe Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with limits
		promotionCode := helpers.NewValidPromotionCodeWithRedemptionLimits("LIMITEDSTRIPE", coupon.ID, 50)
		promotionCode.TimesRedeemed = 10
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by Stripe ID
		retrievedCode, err := repo.GetPromotionCodeByStripeID(ctx, promotionCode.StripePromotionID)

		require.NoError(t, err)
		require.NotNil(t, retrievedCode.MaxRedemptions)
		assert.Equal(t, 50, *retrievedCode.MaxRedemptions)
		assert.Equal(t, 10, retrievedCode.TimesRedeemed)
	})

	t.Run("retrieval with metadata", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Metadata Stripe Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with metadata
		promotionCode := helpers.NewValidPromotionCode("METASTRIPE", coupon.ID)
		promotionCode.Metadata = map[string]string{
			"stripe_campaign": "holiday2024",
			"integration":     "webhook",
		}
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by Stripe ID
		retrievedCode, err := repo.GetPromotionCodeByStripeID(ctx, promotionCode.StripePromotionID)

		require.NoError(t, err)
		assert.Equal(t, "holiday2024", retrievedCode.Metadata["stripe_campaign"])
		assert.Equal(t, "webhook", retrievedCode.Metadata["integration"])
	})
}

