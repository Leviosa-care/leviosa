package couponPayment_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateCoupon(t *testing.T) {
	ctx := context.Background()

	t.Run("successful coupon update with all fields", func(t *testing.T) {
		// First create a coupon
		percentOff := 10.0

		createReq := &domain.CreateCouponRequest{
			Name:       "Original Coupon",
			PercentOff: &percentOff,
			Duration:   "once",
			Metadata: map[string]string{
				"version": "1.0",
			},
		}

		createdCoupon, err := stripeService.CreateCoupon(ctx, createReq)
		require.NoError(t, err)

		// Update the coupon
		newName := "Updated Coupon"
		newMetadata := map[string]string{
			"version": "2.0",
			"updated": "true",
		}

		updateReq := &domain.UpdateCouponRequest{
			Name:     &newName,
			Metadata: newMetadata,
		}

		updatedCoupon, err := stripeService.UpdateCoupon(ctx, createdCoupon.StripeCouponID, updateReq)

		require.NoError(t, err)
		assert.NotNil(t, updatedCoupon)
		assert.Equal(t, createdCoupon.StripeCouponID, updatedCoupon.StripeCouponID)
		assert.Equal(t, newName, updatedCoupon.Name)
		assert.Equal(t, newMetadata, updatedCoupon.Metadata)
		// Other fields should remain unchanged
		assert.Equal(t, createdCoupon.PercentOff, updatedCoupon.PercentOff)
		assert.Equal(t, createdCoupon.Duration, updatedCoupon.Duration)
		assert.Equal(t, createdCoupon.IsValid, updatedCoupon.IsValid)
	})

	t.Run("successful partial update - name only", func(t *testing.T) {
		// Create a coupon
		amountOff := 500 // $5.00
		currency := "USD"

		createReq := &domain.CreateCouponRequest{
			Name:      "Original Name",
			AmountOff: &amountOff,
			Currency:  &currency,
			Duration:  "once",
			Metadata: map[string]string{
				"keep": "this",
			},
		}

		createdCoupon, err := stripeService.CreateCoupon(ctx, createReq)
		require.NoError(t, err)

		// Update only the name
		newName := "Updated Name Only"
		updateReq := &domain.UpdateCouponRequest{
			Name: &newName,
		}

		updatedCoupon, err := stripeService.UpdateCoupon(ctx, createdCoupon.StripeCouponID, updateReq)

		require.NoError(t, err)
		assert.Equal(t, newName, updatedCoupon.Name)
		assert.Equal(t, createdCoupon.AmountOff, updatedCoupon.AmountOff)
		assert.Equal(t, createdCoupon.Currency, updatedCoupon.Currency)
		assert.Equal(t, createdCoupon.Duration, updatedCoupon.Duration)
		assert.Equal(t, createdCoupon.Metadata, updatedCoupon.Metadata) // Should remain unchanged
	})

	t.Run("successful partial update - metadata only", func(t *testing.T) {
		// Create a coupon
		percentOff := 15.0

		createReq := &domain.CreateCouponRequest{
			Name:       "Coupon Name",
			PercentOff: &percentOff,
			Duration:   "once",
			Metadata: map[string]string{
				"original": "value",
			},
		}

		createdCoupon, err := stripeService.CreateCoupon(ctx, createReq)
		require.NoError(t, err)

		// Update only the metadata
		newMetadata := map[string]string{
			"updated":  "metadata",
			"category": "premium",
			"source":   "api",
		}
		updateReq := &domain.UpdateCouponRequest{
			Metadata: newMetadata,
		}

		updatedCoupon, err := stripeService.UpdateCoupon(ctx, createdCoupon.StripeCouponID, updateReq)

		require.NoError(t, err)
		assert.Equal(t, newMetadata, updatedCoupon.Metadata)
		assert.Equal(t, createdCoupon.Name, updatedCoupon.Name)
		assert.Equal(t, createdCoupon.PercentOff, updatedCoupon.PercentOff)
		assert.Equal(t, createdCoupon.Duration, updatedCoupon.Duration)
	})

	t.Run("update fails with non-existent coupon ID", func(t *testing.T) {
		nonExistentID := "coupon_nonexistent123456789"
		newName := "Updated Name"

		updateReq := &domain.UpdateCouponRequest{
			Name: &newName,
		}

		updatedCoupon, err := stripeService.UpdateCoupon(ctx, nonExistentID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, updatedCoupon)
	})

	t.Run("update with empty request does nothing", func(t *testing.T) {
		// Create a coupon
		percentOff := 30.0

		createReq := &domain.CreateCouponRequest{
			Name:       "Original Coupon",
			PercentOff: &percentOff,
			Duration:   "once",
			Metadata: map[string]string{
				"original": "data",
			},
		}

		createdCoupon, err := stripeService.CreateCoupon(ctx, createReq)
		require.NoError(t, err)

		// Update with empty request
		updateReq := &domain.UpdateCouponRequest{}

		updatedCoupon, err := stripeService.UpdateCoupon(ctx, createdCoupon.StripeCouponID, updateReq)

		require.NoError(t, err)
		assert.Equal(t, createdCoupon.Name, updatedCoupon.Name)
		assert.Equal(t, createdCoupon.PercentOff, updatedCoupon.PercentOff)
		assert.Equal(t, createdCoupon.Metadata, updatedCoupon.Metadata)
		assert.Equal(t, createdCoupon.Duration, updatedCoupon.Duration)
	})

	t.Run("update repeating coupon", func(t *testing.T) {
		// Create a repeating coupon
		percentOff := 20.0
		durationInMonths := 3

		createReq := &domain.CreateCouponRequest{
			Name:             "Repeating Coupon",
			PercentOff:       &percentOff,
			Duration:         "repeating",
			DurationInMonths: &durationInMonths,
		}

		createdCoupon, err := stripeService.CreateCoupon(ctx, createReq)
		require.NoError(t, err)

		// Update the name and metadata
		newName := "Updated Repeating Coupon"
		newMetadata := map[string]string{
			"type":    "repeating",
			"updated": "true",
		}

		updateReq := &domain.UpdateCouponRequest{
			Name:     &newName,
			Metadata: newMetadata,
		}

		updatedCoupon, err := stripeService.UpdateCoupon(ctx, createdCoupon.StripeCouponID, updateReq)

		require.NoError(t, err)
		assert.Equal(t, newName, updatedCoupon.Name)
		assert.Equal(t, newMetadata, updatedCoupon.Metadata)
		assert.Equal(t, createdCoupon.PercentOff, updatedCoupon.PercentOff)
		assert.Equal(t, createdCoupon.Duration, updatedCoupon.Duration)
		assert.Equal(t, createdCoupon.DurationInMonths, updatedCoupon.DurationInMonths)
	})

	t.Run("clear metadata with empty map", func(t *testing.T) {
		// Create a coupon with metadata
		percentOff := 12.0

		createReq := &domain.CreateCouponRequest{
			Name:       "Coupon with Metadata",
			PercentOff: &percentOff,
			Duration:   "once",
			Metadata: map[string]string{
				"original": "metadata",
				"to":       "clear",
			},
		}

		createdCoupon, err := stripeService.CreateCoupon(ctx, createReq)
		require.NoError(t, err)
		require.NotEmpty(t, createdCoupon.Metadata)

		// Clear metadata with empty map
		updateReq := &domain.UpdateCouponRequest{
			Metadata: map[string]string{},
		}

		updatedCoupon, err := stripeService.UpdateCoupon(ctx, createdCoupon.StripeCouponID, updateReq)

		require.NoError(t, err)
		assert.Empty(t, updatedCoupon.Metadata)
		assert.Equal(t, createdCoupon.Name, updatedCoupon.Name)
		assert.Equal(t, createdCoupon.PercentOff, updatedCoupon.PercentOff)
	})
}