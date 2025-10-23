package couponRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCoupon(t *testing.T) {
	ctx := context.Background()

	t.Run("successful percent-off coupon creation", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		coupon := helpers.NewValidPercentOffCoupon("25% Off Coupon")

		couponID, err := repo.CreateCoupon(ctx, coupon)

		require.NoError(t, err)
		assert.NotEmpty(t, couponID)

		// Verify coupon was created in database
		createdCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, coupon.ID.String(), couponID)
		assert.Equal(t, coupon.Name, createdCoupon.Name)
		assert.Equal(t, coupon.PercentOff, createdCoupon.PercentOff)
		assert.Equal(t, coupon.Duration, createdCoupon.Duration)
		assert.Equal(t, coupon.IsValid, createdCoupon.IsValid)
		assert.Equal(t, coupon.Metadata, createdCoupon.Metadata)
	})

	t.Run("successful amount-off coupon creation", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		coupon := helpers.NewValidAmountOffCoupon("$5 Off Coupon", "USD")

		couponID, err := repo.CreateCoupon(ctx, coupon)

		require.NoError(t, err)
		assert.NotEmpty(t, couponID)

		// Verify coupon was created in database
		createdCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, coupon.AmountOff, createdCoupon.AmountOff)
		assert.Equal(t, coupon.Currency, createdCoupon.Currency)
		assert.Nil(t, createdCoupon.PercentOff)
	})

	t.Run("successful repeating coupon creation", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		coupon := helpers.NewValidRepeatingCoupon("3 Month Coupon", 3)

		couponID, err := repo.CreateCoupon(ctx, coupon)

		require.NoError(t, err)
		assert.NotEmpty(t, couponID)

		// Verify coupon was created in database
		createdCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, domain.CouponDurationRepeating, createdCoupon.Duration)
		assert.Equal(t, coupon.DurationInMonths, createdCoupon.DurationInMonths)
	})

	t.Run("successful forever coupon creation", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		coupon := helpers.NewValidForeverCoupon("Forever Coupon")

		couponID, err := repo.CreateCoupon(ctx, coupon)

		require.NoError(t, err)
		assert.NotEmpty(t, couponID)

		// Verify coupon was created in database
		createdCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, domain.CouponDurationForever, createdCoupon.Duration)
		assert.Nil(t, createdCoupon.DurationInMonths)
	})

	t.Run("creation with auto-generated UUID", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		coupon := helpers.NewValidPercentOffCoupon("Auto UUID Coupon")
		coupon.ID = uuid.Nil // Clear the ID to test auto-generation

		couponID, err := repo.CreateCoupon(ctx, coupon)

		require.NoError(t, err)
		assert.NotEmpty(t, couponID)

		// Verify a UUID was generated
		parsedID, err := uuid.Parse(couponID)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, parsedID)
	})

	t.Run("creation with expired coupon", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		coupon := helpers.NewExpiredCoupon("Expired Coupon")

		couponID, err := repo.CreateCoupon(ctx, coupon)

		require.NoError(t, err)
		assert.NotEmpty(t, couponID)

		// Verify coupon was created with correct expiry
		createdCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.NotNil(t, createdCoupon.RedeemBy)
		assert.True(t, createdCoupon.RedeemBy.Before(time.Now()))
	})

	t.Run("creation with duplicate stripe coupon ID should fail", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		coupon1 := helpers.NewValidPercentOffCoupon("First Coupon")
		coupon2 := helpers.NewValidPercentOffCoupon("Second Coupon")
		coupon2.StripeCouponID = coupon1.StripeCouponID // Same Stripe ID

		// Create first coupon
		_, err := repo.CreateCoupon(ctx, coupon1)
		require.NoError(t, err)

		// Creating second coupon with same Stripe ID should fail
		_, err = repo.CreateCoupon(ctx, coupon2)
		assert.Error(t, err)
	})

	t.Run("creation with nil metadata", func(t *testing.T) {
		helpers.ClearCouponsTable(t, ctx, testPool)

		coupon := helpers.NewValidPercentOffCoupon("No Metadata Coupon")
		coupon.Metadata = nil

		couponID, err := repo.CreateCoupon(ctx, coupon)

		require.NoError(t, err)
		assert.NotEmpty(t, couponID)

		// Verify coupon was created
		createdCoupon, err := helpers.GetCouponByID(t, ctx, coupon.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, coupon.Name, createdCoupon.Name)
	})
}
