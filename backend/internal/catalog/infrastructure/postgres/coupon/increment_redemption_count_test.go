package couponRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncrementRedemptionCount(t *testing.T) {
	ctx := context.Background()

	t.Run("successful increment from zero", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon with zero redemptions
		coupon := helpers.NewValidPercentOffCoupon("Increment Test")
		coupon.TimesRedeemed = 0
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Increment redemption count
		err := repo.IncrementRedemptionCount(ctx, coupon.ID)

		require.NoError(t, err)

		// Verify increment
		timesRedeemed := helpers.GetCouponTimesRedeemed(t, ctx, coupon.ID, testPool)
		assert.Equal(t, 1, timesRedeemed)
	})

	t.Run("successful increment from existing count", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon with existing redemptions
		coupon := helpers.NewValidPercentOffCoupon("Increment Existing")
		coupon.TimesRedeemed = 5
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Increment redemption count
		err := repo.IncrementRedemptionCount(ctx, coupon.ID)

		require.NoError(t, err)

		// Verify increment
		timesRedeemed := helpers.GetCouponTimesRedeemed(t, ctx, coupon.ID, testPool)
		assert.Equal(t, 6, timesRedeemed)
	})

	t.Run("multiple increments", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidPercentOffCoupon("Multiple Increments")
		coupon.TimesRedeemed = 0
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Perform multiple increments
		for i := 0; i < 3; i++ {
			err := repo.IncrementRedemptionCount(ctx, coupon.ID)
			require.NoError(t, err)
		}

		// Verify final count
		timesRedeemed := helpers.GetCouponTimesRedeemed(t, ctx, coupon.ID, testPool)
		assert.Equal(t, 3, timesRedeemed)
	})

	t.Run("increment for amount-off coupon", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert amount-off coupon
		coupon := helpers.NewValidAmountOffCoupon("Amount Off Increment", "USD")
		coupon.TimesRedeemed = 2
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Increment redemption count
		err := repo.IncrementRedemptionCount(ctx, coupon.ID)

		require.NoError(t, err)

		// Verify increment
		timesRedeemed := helpers.GetCouponTimesRedeemed(t, ctx, coupon.ID, testPool)
		assert.Equal(t, 3, timesRedeemed)
	})

	t.Run("increment for repeating coupon", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert repeating coupon
		coupon := helpers.NewValidRepeatingCoupon("Repeating Increment", 3)
		coupon.TimesRedeemed = 10
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Increment redemption count
		err := repo.IncrementRedemptionCount(ctx, coupon.ID)

		require.NoError(t, err)

		// Verify increment
		timesRedeemed := helpers.GetCouponTimesRedeemed(t, ctx, coupon.ID, testPool)
		assert.Equal(t, 11, timesRedeemed)
	})

	t.Run("increment for forever coupon", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert forever coupon
		coupon := helpers.NewValidForeverCoupon("Forever Increment")
		coupon.TimesRedeemed = 100
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Increment redemption count
		err := repo.IncrementRedemptionCount(ctx, coupon.ID)

		require.NoError(t, err)

		// Verify increment
		timesRedeemed := helpers.GetCouponTimesRedeemed(t, ctx, coupon.ID, testPool)
		assert.Equal(t, 101, timesRedeemed)
	})

	t.Run("increment non-existent coupon should return not found error", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		nonExistentID := uuid.New()

		err := repo.IncrementRedemptionCount(ctx, nonExistentID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("concurrent increments", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidPercentOffCoupon("Concurrent Test")
		coupon.TimesRedeemed = 0
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Perform concurrent increments
		numGoroutines := 5
		errChan := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				errChan <- repo.IncrementRedemptionCount(ctx, coupon.ID)
			}()
		}

		// Collect all errors
		for i := 0; i < numGoroutines; i++ {
			err := <-errChan
			require.NoError(t, err)
		}

		// Verify final count (should be exactly 5 due to database ACID properties)
		timesRedeemed := helpers.GetCouponTimesRedeemed(t, ctx, coupon.ID, testPool)
		assert.Equal(t, 5, timesRedeemed)
	})

	t.Run("increment large number", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon with large redemption count
		coupon := helpers.NewValidPercentOffCoupon("Large Number Test")
		coupon.MaxRedemptions = nil // Remove limit to allow large numbers
		coupon.TimesRedeemed = 999999
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Increment redemption count
		err := repo.IncrementRedemptionCount(ctx, coupon.ID)

		require.NoError(t, err)

		// Verify increment
		timesRedeemed := helpers.GetCouponTimesRedeemed(t, ctx, coupon.ID, testPool)
		assert.Equal(t, 1000000, timesRedeemed)
	})
}

