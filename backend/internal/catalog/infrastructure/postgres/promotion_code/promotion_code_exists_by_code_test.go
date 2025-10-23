package promotionCodeRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromotionCodeExistsByCode(t *testing.T) {
	ctx := context.Background()

	t.Run("returns true for existing code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Exists Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("EXISTS", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Check if code exists
		exists, err := repo.PromotionCodeExistsByCode(ctx, "EXISTS")

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("returns false for non-existent code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Check if non-existent code exists
		exists, err := repo.PromotionCodeExistsByCode(ctx, "NONEXISTENT")

		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("returns true for inactive code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and inactive promotion code
		coupon := helpers.NewValidPercentOffCoupon("Inactive Exists Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewInactivePromotionCode("INACTIVE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Check if inactive code exists
		exists, err := repo.PromotionCodeExistsByCode(ctx, "INACTIVE")

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("returns true for expired code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and expired promotion code
		coupon := helpers.NewValidPercentOffCoupon("Expired Exists Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewExpiredPromotionCode("EXPIRED", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Check if expired code exists
		exists, err := repo.PromotionCodeExistsByCode(ctx, "EXPIRED")

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("case sensitive check", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Case Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("CASE_SENSITIVE", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Check exact case
		exists, err := repo.PromotionCodeExistsByCode(ctx, "CASE_SENSITIVE")
		require.NoError(t, err)
		assert.True(t, exists)

		// Check different case
		exists, err = repo.PromotionCodeExistsByCode(ctx, "case_sensitive")
		require.NoError(t, err)
		assert.False(t, exists)

		// Check mixed case
		exists, err = repo.PromotionCodeExistsByCode(ctx, "Case_Sensitive")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("returns false for empty string", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		// Check if empty code exists
		exists, err := repo.PromotionCodeExistsByCode(ctx, "")

		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("multiple codes - specific existence check", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and multiple promotion codes
		coupon := helpers.NewValidPercentOffCoupon("Multiple Codes Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		code1 := helpers.NewValidPromotionCode("FIRST", coupon.ID)
		code2 := helpers.NewValidPromotionCode("SECOND", coupon.ID)
		code3 := helpers.NewValidPromotionCode("THIRD", coupon.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)
		helpers.InsertPromotionCode(t, ctx, testPool, code3)

		// Check existence of specific codes
		exists1, err1 := repo.PromotionCodeExistsByCode(ctx, "FIRST")
		exists2, err2 := repo.PromotionCodeExistsByCode(ctx, "SECOND")
		exists3, err3 := repo.PromotionCodeExistsByCode(ctx, "THIRD")
		existsNone, errNone := repo.PromotionCodeExistsByCode(ctx, "FOURTH")

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)
		require.NoError(t, errNone)

		assert.True(t, exists1)
		assert.True(t, exists2)
		assert.True(t, exists3)
		assert.False(t, existsNone)
	})

	t.Run("codes with special characters", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with special characters
		coupon := helpers.NewValidPercentOffCoupon("Special Chars Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		specialCode := helpers.NewValidPromotionCode("SPECIAL-CODE_123", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, specialCode)

		// Check if special character code exists
		exists, err := repo.PromotionCodeExistsByCode(ctx, "SPECIAL-CODE_123")

		require.NoError(t, err)
		assert.True(t, exists)

		// Check similar but different code
		existsSimilar, errSimilar := repo.PromotionCodeExistsByCode(ctx, "SPECIAL_CODE-123")
		require.NoError(t, errSimilar)
		assert.False(t, existsSimilar)
	})

	t.Run("codes with restrictions", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with restrictions
		coupon := helpers.NewValidPercentOffCoupon("Restricted Exists Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCodeWithRestrictions("VIP25", coupon.ID, []string{"premium_user"})
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Check if restricted code exists
		exists, err := repo.PromotionCodeExistsByCode(ctx, "VIP25")

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("codes with redemption limits", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with limits
		coupon := helpers.NewValidPercentOffCoupon("Limited Exists Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCodeWithRedemptionLimits("LIMITED50", coupon.ID, 50)
		promotionCode.TimesRedeemed = 25
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Check if limited code exists
		exists, err := repo.PromotionCodeExistsByCode(ctx, "LIMITED50")

		require.NoError(t, err)
		assert.True(t, exists)

		// Code at redemption limit should still exist
		promotionCode.TimesRedeemed = 50
		helpers.UpdatePromotionCodeTimesRedeemed(t, ctx, promotionCode.ID, 50, testPool)

		existsAtLimit, errAtLimit := repo.PromotionCodeExistsByCode(ctx, "LIMITED50")
		require.NoError(t, errAtLimit)
		assert.True(t, existsAtLimit)
	})

	t.Run("codes from different coupons", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create multiple coupons
		coupon1 := helpers.NewValidPercentOffCoupon("Coupon 1")
		coupon2 := helpers.NewValidPercentOffCoupon("Coupon 2")
		helpers.InsertCoupon(t, ctx, testPool, coupon1)
		helpers.InsertCoupon(t, ctx, testPool, coupon2)

		// Create codes for different coupons
		code1 := helpers.NewValidPromotionCode("COUPON1CODE", coupon1.ID)
		code2 := helpers.NewValidPromotionCode("COUPON2CODE", coupon2.ID)

		helpers.InsertPromotionCode(t, ctx, testPool, code1)
		helpers.InsertPromotionCode(t, ctx, testPool, code2)

		// Check existence of codes from different coupons
		exists1, err1 := repo.PromotionCodeExistsByCode(ctx, "COUPON1CODE")
		exists2, err2 := repo.PromotionCodeExistsByCode(ctx, "COUPON2CODE")

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.True(t, exists1)
		assert.True(t, exists2)
	})

	t.Run("whitespace handling", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Whitespace Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("TRIM", coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Check exact match
		exists, err := repo.PromotionCodeExistsByCode(ctx, "TRIM")
		require.NoError(t, err)
		assert.True(t, exists)

		// Check with leading/trailing spaces (should not exist)
		existsSpaced, errSpaced := repo.PromotionCodeExistsByCode(ctx, " TRIM ")
		require.NoError(t, errSpaced)
		assert.False(t, existsSpaced)

		// Check with just spaces (should not exist)
		existsSpaces, errSpaces := repo.PromotionCodeExistsByCode(ctx, "   ")
		require.NoError(t, errSpaces)
		assert.False(t, existsSpaces)
	})

	t.Run("long code names", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with long name
		coupon := helpers.NewValidPercentOffCoupon("Long Code Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		longCodeName := "LONG_PROMOTIONAL_CODE_NAME_FOR_TESTING_MAX50"
		promotionCode := helpers.NewValidPromotionCode(longCodeName, coupon.ID)
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Check if long code exists
		exists, err := repo.PromotionCodeExistsByCode(ctx, longCodeName)

		require.NoError(t, err)
		assert.True(t, exists)

		// Check similar but different long code
		similarLongCode := "LONG_PROMOTIONAL_CODE_NAME_FOR_TESTING_MAX49" // different ending
		existsSimilar, errSimilar := repo.PromotionCodeExistsByCode(ctx, similarLongCode)
		require.NoError(t, errSimilar)
		assert.False(t, existsSimilar)
	})
}

