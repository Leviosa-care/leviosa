package couponRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetValidCoupons(t *testing.T) {
	ctx := context.Background()

	t.Run("retrieves only valid coupons", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert valid coupon
		validCoupon := helpers.NewValidPercentOffCoupon("Valid Coupon")
		helpers.InsertCoupon(t, ctx, testPool, validCoupon)

		// Insert invalid coupon
		invalidCoupon := helpers.NewValidPercentOffCoupon("Invalid Coupon")
		invalidCoupon.IsValid = false
		helpers.InsertCoupon(t, ctx, testPool, invalidCoupon)

		// Get valid coupons
		coupons, err := repo.GetValidCoupons(ctx)

		require.NoError(t, err)
		assert.Len(t, coupons, 1)
		assert.Equal(t, validCoupon.ID, coupons[0].ID)
		assert.True(t, coupons[0].IsValid)
	})

	t.Run("excludes expired coupons", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert valid coupon
		validCoupon := helpers.NewValidPercentOffCoupon("Valid Coupon")
		helpers.InsertCoupon(t, ctx, testPool, validCoupon)

		// Insert expired coupon
		expiredCoupon := helpers.NewExpiredCoupon("Expired Coupon")
		helpers.InsertCoupon(t, ctx, testPool, expiredCoupon)

		// Get valid coupons
		coupons, err := repo.GetValidCoupons(ctx)

		require.NoError(t, err)
		assert.Len(t, coupons, 1)
		assert.Equal(t, validCoupon.ID, coupons[0].ID)
	})

	t.Run("excludes coupons at redemption limit", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert valid coupon
		validCoupon := helpers.NewValidPercentOffCoupon("Valid Coupon")
		helpers.InsertCoupon(t, ctx, testPool, validCoupon)

		// Insert coupon at redemption limit
		limitCoupon := helpers.NewValidPercentOffCoupon("Limit Coupon")
		maxRedemptions := 10
		limitCoupon.MaxRedemptions = &maxRedemptions
		limitCoupon.TimesRedeemed = 10
		helpers.InsertCoupon(t, ctx, testPool, limitCoupon)

		// Get valid coupons
		coupons, err := repo.GetValidCoupons(ctx)

		require.NoError(t, err)
		assert.Len(t, coupons, 1)
		assert.Equal(t, validCoupon.ID, coupons[0].ID)
	})

	t.Run("includes coupons under redemption limit", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert coupon under limit
		underLimitCoupon := helpers.NewValidPercentOffCoupon("Under Limit")
		maxRedemptions := 10
		underLimitCoupon.MaxRedemptions = &maxRedemptions
		underLimitCoupon.TimesRedeemed = 5
		helpers.InsertCoupon(t, ctx, testPool, underLimitCoupon)

		// Get valid coupons
		coupons, err := repo.GetValidCoupons(ctx)

		require.NoError(t, err)
		assert.Len(t, coupons, 1)
		assert.Equal(t, underLimitCoupon.ID, coupons[0].ID)
		assert.Equal(t, 5, coupons[0].TimesRedeemed)
	})

	t.Run("includes coupons with no redemption limit", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert coupon without redemption limit
		unlimitedCoupon := helpers.NewValidPercentOffCoupon("Unlimited")
		unlimitedCoupon.MaxRedemptions = nil
		unlimitedCoupon.TimesRedeemed = 1000
		helpers.InsertCoupon(t, ctx, testPool, unlimitedCoupon)

		// Get valid coupons
		coupons, err := repo.GetValidCoupons(ctx)

		require.NoError(t, err)
		assert.Len(t, coupons, 1)
		assert.Equal(t, unlimitedCoupon.ID, coupons[0].ID)
	})

	t.Run("includes coupons with future expiry", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert coupon with future expiry
		futureCoupon := helpers.NewValidPercentOffCoupon("Future Expiry")
		futureDate := time.Now().Add(24 * time.Hour)
		futureCoupon.RedeemBy = &futureDate
		helpers.InsertCoupon(t, ctx, testPool, futureCoupon)

		// Get valid coupons
		coupons, err := repo.GetValidCoupons(ctx)

		require.NoError(t, err)
		assert.Len(t, coupons, 1)
		assert.Equal(t, futureCoupon.ID, coupons[0].ID)
	})

	t.Run("includes coupons with no expiry date", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert coupon without expiry
		noExpiryCoupon := helpers.NewValidPercentOffCoupon("No Expiry")
		noExpiryCoupon.RedeemBy = nil
		helpers.InsertCoupon(t, ctx, testPool, noExpiryCoupon)

		// Get valid coupons
		coupons, err := repo.GetValidCoupons(ctx)

		require.NoError(t, err)
		assert.Len(t, coupons, 1)
		assert.Equal(t, noExpiryCoupon.ID, coupons[0].ID)
	})

	t.Run("returns empty list when no valid coupons", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert only invalid coupons
		invalidCoupon := helpers.NewValidPercentOffCoupon("Invalid")
		invalidCoupon.IsValid = false
		helpers.InsertCoupon(t, ctx, testPool, invalidCoupon)

		expiredCoupon := helpers.NewExpiredCoupon("Expired")
		helpers.InsertCoupon(t, ctx, testPool, expiredCoupon)

		// Get valid coupons
		coupons, err := repo.GetValidCoupons(ctx)

		require.NoError(t, err)
		assert.Empty(t, coupons)
	})

	t.Run("returns coupons ordered by creation date descending", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert multiple valid coupons with different creation times
		coupon1 := helpers.NewValidPercentOffCoupon("First Coupon")
		coupon1.CreatedAt = time.Now().Add(-2 * time.Hour)
		helpers.InsertCoupon(t, ctx, testPool, coupon1)

		coupon2 := helpers.NewValidPercentOffCoupon("Second Coupon")
		coupon2.CreatedAt = time.Now().Add(-1 * time.Hour)
		helpers.InsertCoupon(t, ctx, testPool, coupon2)

		coupon3 := helpers.NewValidPercentOffCoupon("Third Coupon")
		coupon3.CreatedAt = time.Now()
		helpers.InsertCoupon(t, ctx, testPool, coupon3)

		// Get valid coupons
		coupons, err := repo.GetValidCoupons(ctx)

		require.NoError(t, err)
		assert.Len(t, coupons, 3)

		// Should be ordered by creation date descending (newest first)
		assert.Equal(t, coupon3.ID, coupons[0].ID)
		assert.Equal(t, coupon2.ID, coupons[1].ID)
		assert.Equal(t, coupon1.ID, coupons[2].ID)
	})

	t.Run("returns mixed coupon types", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert different types of valid coupons
		percentCoupon := helpers.NewValidPercentOffCoupon("Percent Coupon")
		helpers.InsertCoupon(t, ctx, testPool, percentCoupon)

		amountCoupon := helpers.NewValidAmountOffCoupon("Amount Coupon", "USD")
		helpers.InsertCoupon(t, ctx, testPool, amountCoupon)

		repeatingCoupon := helpers.NewValidRepeatingCoupon("Repeating Coupon", 6)
		helpers.InsertCoupon(t, ctx, testPool, repeatingCoupon)

		foreverCoupon := helpers.NewValidForeverCoupon("Forever Coupon")
		helpers.InsertCoupon(t, ctx, testPool, foreverCoupon)

		// Get valid coupons
		coupons, err := repo.GetValidCoupons(ctx)

		require.NoError(t, err)
		assert.Len(t, coupons, 4)

		// Verify all types are included
		couponIDs := make(map[string]bool)
		for _, c := range coupons {
			couponIDs[c.ID.String()] = true
		}

		assert.True(t, couponIDs[percentCoupon.ID.String()])
		assert.True(t, couponIDs[amountCoupon.ID.String()])
		assert.True(t, couponIDs[repeatingCoupon.ID.String()])
		assert.True(t, couponIDs[foreverCoupon.ID.String()])
	})
}

