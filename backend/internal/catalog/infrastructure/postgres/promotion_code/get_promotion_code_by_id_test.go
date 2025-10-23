package promotionCodeRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPromotionCodeByID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful retrieval by ID", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("GETBYID", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve promotion code by ID
		retrievedCode, err := repo.GetPromotionCodeByID(ctx, promotionCode.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCode)
		assert.Equal(t, promotionCode.ID, retrievedCode.ID)
		assert.Equal(t, promotionCode.Code, retrievedCode.Code)
		assert.Equal(t, promotionCode.CouponID, retrievedCode.CouponID)
		assert.Equal(t, promotionCode.StripePromotionID, retrievedCode.StripePromotionID)
		assert.Equal(t, promotionCode.Active, retrievedCode.Active)
		assert.Equal(t, promotionCode.Metadata, retrievedCode.Metadata)
	})

	t.Run("successful retrieval with restrictions", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Restricted Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with restrictions
		promotionCode := helpers.NewValidPromotionCodeWithRestrictions("RESTRICTED", coupon.ID, []string{"customer_123"})
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by ID
		retrievedCode, err := repo.GetPromotionCodeByID(ctx, promotionCode.ID)

		require.NoError(t, err)
		assert.Equal(t, []string{"customer_123"}, retrievedCode.Restrictions.CurrencyOptions)
	})

	t.Run("successful retrieval with redemption limits", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Limited Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with limits
		promotionCode := helpers.NewValidPromotionCodeWithRedemptionLimits("LIMITED", coupon.ID, 25)
		promotionCode.TimesRedeemed = 5
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by ID
		retrievedCode, err := repo.GetPromotionCodeByID(ctx, promotionCode.ID)

		require.NoError(t, err)
		require.NotNil(t, retrievedCode.MaxRedemptions)
		assert.Equal(t, 25, *retrievedCode.MaxRedemptions)
		assert.Equal(t, 5, retrievedCode.TimesRedeemed)
	})

	t.Run("retrieval with non-existent ID should return not found error", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		nonExistentID := uuid.New()

		retrievedCode, err := repo.GetPromotionCodeByID(ctx, nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, retrievedCode)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("retrieval with nil UUID should return error", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		retrievedCode, err := repo.GetPromotionCodeByID(ctx, uuid.Nil)

		assert.Error(t, err)
		assert.Nil(t, retrievedCode)
	})

	t.Run("retrieval of expired promotion code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Expired Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert expired promotion code
		promotionCode := helpers.NewExpiredPromotionCode("EXPIRED", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by ID
		retrievedCode, err := repo.GetPromotionCodeByID(ctx, promotionCode.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCode.ExpiresAt)
		assert.Equal(t, promotionCode.ExpiresAt.Unix(), retrievedCode.ExpiresAt.Unix())
	})

	t.Run("retrieval of inactive promotion code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Inactive Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert inactive promotion code
		promotionCode := helpers.NewInactivePromotionCode("INACTIVE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by ID
		retrievedCode, err := repo.GetPromotionCodeByID(ctx, promotionCode.ID)

		require.NoError(t, err)
		assert.False(t, retrievedCode.Active)
	})

	t.Run("retrieval with no metadata", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("No Metadata Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code without metadata
		promotionCode := helpers.NewValidPromotionCode("NOMETA", coupon.ID)
		promotionCode.Metadata = nil
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by ID
		retrievedCode, err := repo.GetPromotionCodeByID(ctx, promotionCode.ID)

		require.NoError(t, err)
		assert.Equal(t, promotionCode.Code, retrievedCode.Code)
	})

	t.Run("retrieval with complex metadata", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Metadata Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with metadata
		promotionCode := helpers.NewValidPromotionCode("METADATA", coupon.ID)
		promotionCode.Metadata = map[string]string{
			"campaign":    "winter2024",
			"channel":     "social_media",
			"target_user": "premium",
		}
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by ID
		retrievedCode, err := repo.GetPromotionCodeByID(ctx, promotionCode.ID)

		require.NoError(t, err)
		assert.Equal(t, "winter2024", retrievedCode.Metadata["campaign"])
		assert.Equal(t, "social_media", retrievedCode.Metadata["channel"])
		assert.Equal(t, "premium", retrievedCode.Metadata["target_user"])
	})
}

