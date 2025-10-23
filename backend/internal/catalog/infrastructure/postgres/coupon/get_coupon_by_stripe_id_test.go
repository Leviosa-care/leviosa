package couponRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCouponByStripeID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful retrieval by stripe ID", func(t *testing.T) {
		testdata.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidPercentOffCoupon("Test Stripe ID Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Retrieve coupon by Stripe ID
		retrievedCoupon, err := repo.GetCouponByStripeID(ctx, coupon.StripeCouponID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCoupon)
		assert.Equal(t, coupon.ID, retrievedCoupon.ID)
		assert.Equal(t, coupon.StripeCouponID, retrievedCoupon.StripeCouponID)
		assert.Equal(t, coupon.Name, retrievedCoupon.Name)
		assert.Equal(t, coupon.PercentOff, retrievedCoupon.PercentOff)
		assert.Equal(t, coupon.Duration, retrievedCoupon.Duration)
	})

	t.Run("successful retrieval of amount-off coupon by stripe ID", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidAmountOffCoupon("Test Amount Stripe", "GBP")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Retrieve coupon by Stripe ID
		retrievedCoupon, err := repo.GetCouponByStripeID(ctx, coupon.StripeCouponID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCoupon)
		assert.Equal(t, coupon.AmountOff, retrievedCoupon.AmountOff)
		assert.Equal(t, coupon.Currency, retrievedCoupon.Currency)
	})

	t.Run("retrieval with non-existent stripe ID should return not found error", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		nonExistentStripeID := "coupon_nonexistent123"

		retrievedCoupon, err := repo.GetCouponByStripeID(ctx, nonExistentStripeID)

		assert.Error(t, err)
		assert.Nil(t, retrievedCoupon)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("retrieval with empty stripe ID should return error", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		retrievedCoupon, err := repo.GetCouponByStripeID(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, retrievedCoupon)
	})

	t.Run("multiple coupons, retrieve specific one by stripe ID", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert multiple coupons
		coupon1 := helpers.NewValidPercentOffCoupon("First Coupon")
		coupon2 := helpers.NewValidAmountOffCoupon("Second Coupon", "USD")
		coupon3 := helpers.NewValidForeverCoupon("Third Coupon")

		helpers.InsertCoupon(t, ctx, testPool, coupon1)
		helpers.InsertCoupon(t, ctx, testPool, coupon2)
		helpers.InsertCoupon(t, ctx, testPool, coupon3)

		// Retrieve specific coupon by Stripe ID
		retrievedCoupon, err := repo.GetCouponByStripeID(ctx, coupon2.StripeCouponID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCoupon)
		assert.Equal(t, coupon2.ID, retrievedCoupon.ID)
		assert.Equal(t, coupon2.Name, retrievedCoupon.Name)
		assert.Equal(t, coupon2.AmountOff, retrievedCoupon.AmountOff)
	})

	t.Run("retrieval of expired coupon by stripe ID", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert expired coupon
		coupon := helpers.NewExpiredCoupon("Expired Stripe Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Retrieve coupon by Stripe ID
		retrievedCoupon, err := repo.GetCouponByStripeID(ctx, coupon.StripeCouponID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCoupon)
		assert.NotNil(t, retrievedCoupon.RedeemBy)
		assert.Equal(t, coupon.StripeCouponID, retrievedCoupon.StripeCouponID)
	})

	t.Run("retrieval with special characters in stripe ID", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert coupon with special Stripe ID
		coupon := helpers.NewValidPercentOffCoupon("Special Char Coupon")
		coupon.StripeCouponID = "coupon_test-123_special.id"
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Retrieve coupon by Stripe ID
		retrievedCoupon, err := repo.GetCouponByStripeID(ctx, coupon.StripeCouponID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCoupon)
		assert.Equal(t, coupon.StripeCouponID, retrievedCoupon.StripeCouponID)
	})
}

