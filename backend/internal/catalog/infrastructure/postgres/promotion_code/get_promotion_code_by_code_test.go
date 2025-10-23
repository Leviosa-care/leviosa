package promotionCodeRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPromotionCodeByCode(t *testing.T) {
	ctx := context.Background()

	t.Run("successful retrieval by code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Code Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("GETBYCODE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve promotion code by code
		retrievedCode, err := repo.GetPromotionCodeByCode(ctx, "GETBYCODE")

		require.NoError(t, err)
		assert.NotNil(t, retrievedCode)
		assert.Equal(t, promotionCode.ID, retrievedCode.ID)
		assert.Equal(t, promotionCode.Code, retrievedCode.Code)
		assert.Equal(t, promotionCode.CouponID, retrievedCode.CouponID)
		assert.Equal(t, promotionCode.StripePromotionID, retrievedCode.StripePromotionID)
		assert.Equal(t, promotionCode.Active, retrievedCode.Active)
	})

	t.Run("successful retrieval with case sensitivity", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Case Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with mixed case
		promotionCode := helpers.NewValidPromotionCode("MiXeDcAsE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve with exact case
		retrievedCode, err := repo.GetPromotionCodeByCode(ctx, "MiXeDcAsE")

		require.NoError(t, err)
		assert.Equal(t, "MiXeDcAsE", retrievedCode.Code)
	})

	t.Run("retrieval with non-existent code should return not found error", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		nonExistentCode := "NONEXISTENT"

		retrievedCode, err := repo.GetPromotionCodeByCode(ctx, nonExistentCode)

		assert.Error(t, err)
		assert.Nil(t, retrievedCode)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("retrieval with empty code should return error", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		retrievedCode, err := repo.GetPromotionCodeByCode(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, retrievedCode)
	})

	t.Run("multiple codes, retrieve specific one", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Multiple Codes Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert multiple promotion codes
		code1 := helpers.NewValidPromotionCode("FIRST", coupon.ID)
		code2 := helpers.NewValidPromotionCode("SECOND", coupon.ID)
		code3 := helpers.NewValidPromotionCode("THIRD", coupon.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)
		helpers.InsertPromotionCode(t, ctx, testPool, code3)

		// Retrieve specific code
		retrievedCode, err := repo.GetPromotionCodeByCode(ctx, "SECOND")

		require.NoError(t, err)
		assert.Equal(t, code2.ID, retrievedCode.ID)
		assert.Equal(t, "SECOND", retrievedCode.Code)
	})

	t.Run("retrieval of expired promotion code by code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Expired Code Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert expired promotion code
		promotionCode := helpers.NewExpiredPromotionCode("EXPIREDCODE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by code
		retrievedCode, err := repo.GetPromotionCodeByCode(ctx, "EXPIREDCODE")

		require.NoError(t, err)
		assert.NotNil(t, retrievedCode.ExpiresAt)
		assert.Equal(t, promotionCode.ExpiresAt.Unix(), retrievedCode.ExpiresAt.Unix())
	})

	t.Run("retrieval of inactive promotion code by code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Inactive Code Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert inactive promotion code
		promotionCode := helpers.NewInactivePromotionCode("INACTIVECODE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by code
		retrievedCode, err := repo.GetPromotionCodeByCode(ctx, "INACTIVECODE")

		require.NoError(t, err)
		assert.False(t, retrievedCode.Active)
	})

	t.Run("retrieval with special characters in code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Special Char Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with special characters
		promotionCode := helpers.NewValidPromotionCode("SAVE-25_NOW", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by code
		retrievedCode, err := repo.GetPromotionCodeByCode(ctx, "SAVE-25_NOW")

		require.NoError(t, err)
		assert.Equal(t, "SAVE-25_NOW", retrievedCode.Code)
	})

	t.Run("retrieval with restrictions", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Restricted Code Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with restrictions
		promotionCode := helpers.NewValidPromotionCodeWithRestrictions("VIP25", coupon.ID, []string{"vip_customer_1", "vip_customer_2"})
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by code
		retrievedCode, err := repo.GetPromotionCodeByCode(ctx, "VIP25")

		require.NoError(t, err)
		assert.Equal(t, []string{"vip_customer_1", "vip_customer_2"}, retrievedCode.Restrictions.CurrencyOptions)
	})

	t.Run("retrieval with redemption data", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon first
		coupon := helpers.NewValidPercentOffCoupon("Redemption Data Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Insert promotion code with redemption limits
		promotionCode := helpers.NewValidPromotionCodeWithRedemptionLimits("REDEMPTION", coupon.ID, 100)
		promotionCode.TimesRedeemed = 25
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Retrieve by code
		retrievedCode, err := repo.GetPromotionCodeByCode(ctx, "REDEMPTION")

		require.NoError(t, err)
		require.NotNil(t, retrievedCode.MaxRedemptions)
		assert.Equal(t, 100, *retrievedCode.MaxRedemptions)
		assert.Equal(t, 25, retrievedCode.TimesRedeemed)
	})
}
