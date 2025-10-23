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

func TestIncrementPromotionCodeRedemptionCount(t *testing.T) {
	ctx := context.Background()

	t.Run("successful increment from zero", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Increment Test Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("INCREMENT", coupon.ID)
		promotionCode.TimesRedeemed = 0
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Increment redemption count
		err := repo.IncrementRedemptionCount(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify increment
		timesRedeemed := helpers.GetPromotionCodeTimesRedeemed(t, ctx, promotionCode.ID, testPool)
		assert.Equal(t, 1, timesRedeemed)
	})

	t.Run("successful increment from existing count", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Existing Count Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("EXISTING", coupon.ID)
		promotionCode.TimesRedeemed = 15
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Increment redemption count
		err := repo.IncrementRedemptionCount(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify increment
		timesRedeemed := helpers.GetPromotionCodeTimesRedeemed(t, ctx, promotionCode.ID, testPool)
		assert.Equal(t, 16, timesRedeemed)
	})

	t.Run("multiple increments", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Multiple Increments")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("MULTIPLE", coupon.ID)
		promotionCode.TimesRedeemed = 0
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Perform multiple increments
		for i := 0; i < 5; i++ {
			err := repo.IncrementRedemptionCount(ctx, promotionCode.ID)
			require.NoError(t, err)
		}

		// Verify final count
		timesRedeemed := helpers.GetPromotionCodeTimesRedeemed(t, ctx, promotionCode.ID, testPool)
		assert.Equal(t, 5, timesRedeemed)
	})

	t.Run("increment with redemption limits", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with limits
		coupon := helpers.NewValidPercentOffCoupon("Limited Increments")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCodeWithRedemptionLimits("LIMITED", coupon.ID, 10)
		promotionCode.TimesRedeemed = 7
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Increment redemption count
		err := repo.IncrementRedemptionCount(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify increment
		timesRedeemed := helpers.GetPromotionCodeTimesRedeemed(t, ctx, promotionCode.ID, testPool)
		assert.Equal(t, 8, timesRedeemed)
	})

	t.Run("increment expired promotion code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and expired promotion code
		coupon := helpers.NewValidPercentOffCoupon("Expired Increment Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewExpiredPromotionCode("EXPIRED", coupon.ID)
		promotionCode.TimesRedeemed = 2
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Increment redemption count (should still work)
		err := repo.IncrementRedemptionCount(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify increment
		timesRedeemed := helpers.GetPromotionCodeTimesRedeemed(t, ctx, promotionCode.ID, testPool)
		assert.Equal(t, 3, timesRedeemed)
	})

	t.Run("increment inactive promotion code", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and inactive promotion code
		coupon := helpers.NewValidPercentOffCoupon("Inactive Increment Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewInactivePromotionCode("INACTIVE", coupon.ID)
		promotionCode.TimesRedeemed = 1
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Increment redemption count (should still work)
		err := repo.IncrementRedemptionCount(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify increment
		timesRedeemed := helpers.GetPromotionCodeTimesRedeemed(t, ctx, promotionCode.ID, testPool)
		assert.Equal(t, 2, timesRedeemed)
	})

	t.Run("increment non-existent promotion code should return not found error", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)

		nonExistentID := uuid.New()

		err := repo.IncrementRedemptionCount(ctx, nonExistentID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("concurrent increments", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code
		coupon := helpers.NewValidPercentOffCoupon("Concurrent Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCode("CONCURRENT", coupon.ID)
		promotionCode.TimesRedeemed = 0
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Perform concurrent increments
		numGoroutines := 3
		errChan := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				errChan <- repo.IncrementRedemptionCount(ctx, promotionCode.ID)
			}()
		}

		// Collect all errors
		for i := 0; i < numGoroutines; i++ {
			err := <-errChan
			require.NoError(t, err)
		}

		// Verify final count (should be exactly 3 due to database ACID properties)
		timesRedeemed := helpers.GetPromotionCodeTimesRedeemed(t, ctx, promotionCode.ID, testPool)
		assert.Equal(t, 3, timesRedeemed)
	})

	t.Run("increment large number", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with large count
		coupon := helpers.NewValidPercentOffCoupon("Large Number Test")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCodeWithRedemptionLimits("LARGE", coupon.ID, 2000000)
		promotionCode.TimesRedeemed = 999999
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Increment redemption count
		err := repo.IncrementRedemptionCount(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify increment
		timesRedeemed := helpers.GetPromotionCodeTimesRedeemed(t, ctx, promotionCode.ID, testPool)
		assert.Equal(t, 1000000, timesRedeemed)
	})

	t.Run("increment with restrictions", func(t *testing.T) {
		helpers.ClearPromotionCodesTable(t, ctx, testPool)
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Create a coupon and promotion code with restrictions
		coupon := helpers.NewValidPercentOffCoupon("Restricted Increment")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		promotionCode := helpers.NewValidPromotionCodeWithRestrictions("RESTRICTED", coupon.ID, []string{"vip_user"})
		promotionCode.TimesRedeemed = 3
		helpers.InsertPromotionCode(t, ctx, testPool, promotionCode)

		// Increment redemption count
		err := repo.IncrementRedemptionCount(ctx, promotionCode.ID)

		require.NoError(t, err)

		// Verify increment and that restrictions are preserved
		timesRedeemed := helpers.GetPromotionCodeTimesRedeemed(t, ctx, promotionCode.ID, testPool)
		assert.Equal(t, 4, timesRedeemed)

		// Verify restrictions still intact
		updatedCode, err := helpers.GetPromotionCodeByID(t, ctx, promotionCode.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, []string{"vip_user"}, updatedCode.Restrictions.CurrencyOptions)
	})
}

