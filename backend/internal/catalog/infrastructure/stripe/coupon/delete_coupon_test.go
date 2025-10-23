package couponPayment_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteCoupon(t *testing.T) {
	ctx := context.Background()

	t.Run("successful coupon deletion", func(t *testing.T) {
		// Create a coupon to delete
		percentOff := 25.0

		createReq := &domain.CreateCouponRequest{
			Name:       "Coupon to Delete",
			PercentOff: &percentOff,
			Duration:   "once",
			Metadata: map[string]string{
				"test": "deletion",
			},
		}

		createdCoupon, err := stripeService.CreateCoupon(ctx, createReq)
		require.NoError(t, err)
		require.NotNil(t, createdCoupon)

		// Verify the coupon exists before deletion
		retrievedCoupon, err := stripeService.GetCoupon(ctx, createdCoupon.StripeCouponID)
		require.NoError(t, err)
		assert.NotNil(t, retrievedCoupon)

		// Delete the coupon
		err = stripeService.DeleteCoupon(ctx, createdCoupon.StripeCouponID)

		require.NoError(t, err)

		// Verify the coupon no longer exists
		deletedCoupon, err := stripeService.GetCoupon(ctx, createdCoupon.StripeCouponID)
		assert.Error(t, err)
		assert.Nil(t, deletedCoupon)
	})

	t.Run("delete amount-off coupon", func(t *testing.T) {
		// Create an amount-off coupon to delete
		amountOff := 750 // $7.50
		currency := "USD"

		createReq := &domain.CreateCouponRequest{
			Name:      "Amount Off Coupon to Delete",
			AmountOff: &amountOff,
			Currency:  &currency,
			Duration:  "once",
		}

		createdCoupon, err := stripeService.CreateCoupon(ctx, createReq)
		require.NoError(t, err)

		// Delete the coupon
		err = stripeService.DeleteCoupon(ctx, createdCoupon.StripeCouponID)

		require.NoError(t, err)

		// Verify deletion
		deletedCoupon, err := stripeService.GetCoupon(ctx, createdCoupon.StripeCouponID)
		assert.Error(t, err)
		assert.Nil(t, deletedCoupon)
	})

	t.Run("delete repeating coupon", func(t *testing.T) {
		// Create a repeating coupon to delete
		percentOff := 15.0
		durationInMonths := 6

		createReq := &domain.CreateCouponRequest{
			Name:             "Repeating Coupon to Delete",
			PercentOff:       &percentOff,
			Duration:         "repeating",
			DurationInMonths: &durationInMonths,
		}

		createdCoupon, err := stripeService.CreateCoupon(ctx, createReq)
		require.NoError(t, err)

		// Delete the coupon
		err = stripeService.DeleteCoupon(ctx, createdCoupon.StripeCouponID)

		require.NoError(t, err)

		// Verify deletion
		deletedCoupon, err := stripeService.GetCoupon(ctx, createdCoupon.StripeCouponID)
		assert.Error(t, err)
		assert.Nil(t, deletedCoupon)
	})

	t.Run("delete forever coupon", func(t *testing.T) {
		// Create a forever coupon to delete
		percentOff := 5.0
		maxRedemptions := 1000

		createReq := &domain.CreateCouponRequest{
			Name:           "Forever Coupon to Delete",
			PercentOff:     &percentOff,
			Duration:       "forever",
			MaxRedemptions: &maxRedemptions,
		}

		createdCoupon, err := stripeService.CreateCoupon(ctx, createReq)
		require.NoError(t, err)

		// Delete the coupon
		err = stripeService.DeleteCoupon(ctx, createdCoupon.StripeCouponID)

		require.NoError(t, err)

		// Verify deletion
		deletedCoupon, err := stripeService.GetCoupon(ctx, createdCoupon.StripeCouponID)
		assert.Error(t, err)
		assert.Nil(t, deletedCoupon)
	})

	t.Run("deletion fails with non-existent coupon ID", func(t *testing.T) {
		nonExistentID := "coupon_nonexistent123456789"

		err := stripeService.DeleteCoupon(ctx, nonExistentID)

		assert.Error(t, err)
	})

	t.Run("delete coupon with metadata", func(t *testing.T) {
		// Create a coupon with extensive metadata
		percentOff := 35.0

		createReq := &domain.CreateCouponRequest{
			Name:       "Coupon with Metadata to Delete",
			PercentOff: &percentOff,
			Duration:   "once",
			Metadata: map[string]string{
				"campaign":     "holiday_sale",
				"type":         "percent_off",
				"target_group": "premium_users",
				"region":       "north_america",
			},
		}

		createdCoupon, err := stripeService.CreateCoupon(ctx, createReq)
		require.NoError(t, err)
		require.NotEmpty(t, createdCoupon.Metadata)

		// Delete the coupon
		err = stripeService.DeleteCoupon(ctx, createdCoupon.StripeCouponID)

		require.NoError(t, err)

		// Verify deletion
		deletedCoupon, err := stripeService.GetCoupon(ctx, createdCoupon.StripeCouponID)
		assert.Error(t, err)
		assert.Nil(t, deletedCoupon)
	})
}

func TestCouponLifecycle(t *testing.T) {
	ctx := context.Background()

	t.Run("full coupon lifecycle - create, get, update, delete", func(t *testing.T) {
		// Create
		percentOff := 20.0

		createReq := &domain.CreateCouponRequest{
			Name:       "Lifecycle Test Coupon",
			PercentOff: &percentOff,
			Duration:   "once",
			Metadata: map[string]string{
				"test": "lifecycle",
			},
		}

		createdCoupon, err := stripeService.CreateCoupon(ctx, createReq)
		require.NoError(t, err)
		assert.Equal(t, "Lifecycle Test Coupon", createdCoupon.Name)

		// Get
		retrievedCoupon, err := stripeService.GetCoupon(ctx, createdCoupon.StripeCouponID)
		require.NoError(t, err)
		assert.Equal(t, createdCoupon.StripeCouponID, retrievedCoupon.StripeCouponID)

		// Update
		newName := "Updated Lifecycle Coupon"
		updateReq := &domain.UpdateCouponRequest{
			Name: &newName,
		}

		updatedCoupon, err := stripeService.UpdateCoupon(ctx, createdCoupon.StripeCouponID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, newName, updatedCoupon.Name)

		// Delete
		err = stripeService.DeleteCoupon(ctx, createdCoupon.StripeCouponID)
		require.NoError(t, err)

		// Verify deletion
		deletedCoupon, err := stripeService.GetCoupon(ctx, createdCoupon.StripeCouponID)
		assert.Error(t, err)
		assert.Nil(t, deletedCoupon)
	})
}