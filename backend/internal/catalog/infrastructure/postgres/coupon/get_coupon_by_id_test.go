package couponRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCouponByID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful retrieval of percent-off coupon", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidPercentOffCoupon("Test Percent Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Retrieve coupon
		retrievedCoupon, err := repo.GetCouponByID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCoupon)
		assert.Equal(t, coupon.ID, retrievedCoupon.ID)
		assert.Equal(t, coupon.Name, retrievedCoupon.Name)
		assert.Equal(t, coupon.StripeCouponID, retrievedCoupon.StripeCouponID)
		assert.Equal(t, coupon.PercentOff, retrievedCoupon.PercentOff)
		assert.Equal(t, coupon.Duration, retrievedCoupon.Duration)
		assert.Equal(t, coupon.IsValid, retrievedCoupon.IsValid)
		assert.Equal(t, coupon.Metadata, retrievedCoupon.Metadata)
	})

	t.Run("successful retrieval of amount-off coupon", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidAmountOffCoupon("Test Amount Coupon", "EUR")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Retrieve coupon
		retrievedCoupon, err := repo.GetCouponByID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCoupon)
		assert.Equal(t, coupon.AmountOff, retrievedCoupon.AmountOff)
		assert.Equal(t, coupon.Currency, retrievedCoupon.Currency)
		assert.Nil(t, retrievedCoupon.PercentOff)
	})

	t.Run("successful retrieval of repeating coupon", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidRepeatingCoupon("Test Repeating Coupon", 6)
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Retrieve coupon
		retrievedCoupon, err := repo.GetCouponByID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCoupon)
		assert.Equal(t, coupon.Duration, retrievedCoupon.Duration)
		assert.Equal(t, coupon.DurationInMonths, retrievedCoupon.DurationInMonths)
	})

	t.Run("successful retrieval of forever coupon", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewValidForeverCoupon("Test Forever Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Retrieve coupon
		retrievedCoupon, err := repo.GetCouponByID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCoupon)
		assert.Equal(t, coupon.Duration, retrievedCoupon.Duration)
		assert.Nil(t, retrievedCoupon.DurationInMonths)
	})

	t.Run("successful retrieval of expired coupon", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon
		coupon := helpers.NewExpiredCoupon("Test Expired Coupon")
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Retrieve coupon
		retrievedCoupon, err := repo.GetCouponByID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCoupon)
		assert.NotNil(t, retrievedCoupon.RedeemBy)
		assert.Equal(t, coupon.RedeemBy.Unix(), retrievedCoupon.RedeemBy.Unix())
	})

	t.Run("retrieval with non-existent ID should return not found error", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		nonExistentID := uuid.New()

		retrievedCoupon, err := repo.GetCouponByID(ctx, nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, retrievedCoupon)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("retrieval with nil UUID should return error", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		retrievedCoupon, err := repo.GetCouponByID(ctx, uuid.Nil)

		assert.Error(t, err)
		assert.Nil(t, retrievedCoupon)
	})

	t.Run("retrieval of coupon with no metadata", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon without metadata
		coupon := helpers.NewValidPercentOffCoupon("Test No Metadata")
		coupon.Metadata = nil
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Retrieve coupon
		retrievedCoupon, err := repo.GetCouponByID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCoupon)
		assert.Equal(t, coupon.Name, retrievedCoupon.Name)
	})

	t.Run("retrieval of coupon with redemption data", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		// Insert test coupon with redemption data
		coupon := helpers.NewValidPercentOffCoupon("Test Redemption Data")
		coupon.TimesRedeemed = 5
		helpers.InsertCoupon(t, ctx, testPool, coupon)

		// Retrieve coupon
		retrievedCoupon, err := repo.GetCouponByID(ctx, coupon.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedCoupon)
		assert.Equal(t, coupon.TimesRedeemed, retrievedCoupon.TimesRedeemed)
		assert.Equal(t, coupon.MaxRedemptions, retrievedCoupon.MaxRedemptions)
	})
}

