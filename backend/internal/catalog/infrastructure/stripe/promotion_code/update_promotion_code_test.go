package promotionCodePayment_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdatePromotionCode(t *testing.T) {
	ctx := context.Background()

	createMockCoupon := func() string {
		return "coupon_test123456789"
	}

	t.Run("successful promotion code update with all fields", func(t *testing.T) {
		// First create a promotion code
		couponID := createMockCoupon()

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "ORIGINAL",
			FirstTimeTransaction: true,
			Metadata: map[string]string{
				"version": "1.0",
			},
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)
		require.True(t, createdPromotionCode.Active)

		// Update the promotion code
		newActive := false
		newMetadata := map[string]string{
			"version": "2.0",
			"updated": "true",
			"status":  "deactivated",
		}

		updateReq := &domain.UpdatePromotionCodeRequest{
			Active:   &newActive,
			Metadata: newMetadata,
		}

		updatedPromotionCode, err := stripeService.UpdatePromotionCode(ctx, createdPromotionCode.StripePromotionID, updateReq)

		require.NoError(t, err)
		assert.NotNil(t, updatedPromotionCode)
		assert.Equal(t, createdPromotionCode.StripePromotionID, updatedPromotionCode.StripePromotionID)
		assert.Equal(t, newActive, updatedPromotionCode.Active)
		assert.Equal(t, newMetadata, updatedPromotionCode.Metadata)
		// Other fields should remain unchanged
		assert.Equal(t, createdPromotionCode.Code, updatedPromotionCode.Code)
		assert.Equal(t, createdPromotionCode.FirstTimeTransaction, updatedPromotionCode.FirstTimeTransaction)
	})

	t.Run("successful partial update - active status only", func(t *testing.T) {
		// Create a promotion code
		couponID := createMockCoupon()

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "ACTIVETEST",
			FirstTimeTransaction: false,
			Metadata: map[string]string{
				"keep": "this_metadata",
			},
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)
		require.True(t, createdPromotionCode.Active)

		// Update only the active status
		newActive := false
		updateReq := &domain.UpdatePromotionCodeRequest{
			Active: &newActive,
		}

		updatedPromotionCode, err := stripeService.UpdatePromotionCode(ctx, createdPromotionCode.StripePromotionID, updateReq)

		require.NoError(t, err)
		assert.False(t, updatedPromotionCode.Active)
		assert.Equal(t, createdPromotionCode.Code, updatedPromotionCode.Code)
		assert.Equal(t, createdPromotionCode.FirstTimeTransaction, updatedPromotionCode.FirstTimeTransaction)
		assert.Equal(t, createdPromotionCode.Metadata, updatedPromotionCode.Metadata) // Should remain unchanged
	})

	t.Run("successful partial update - metadata only", func(t *testing.T) {
		// Create a promotion code
		couponID := createMockCoupon()

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "METAUPDATE",
			FirstTimeTransaction: true,
			Metadata: map[string]string{
				"original": "value",
			},
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)

		// Update only the metadata
		newMetadata := map[string]string{
			"updated":  "metadata",
			"category": "premium",
			"source":   "api_update",
		}
		updateReq := &domain.UpdatePromotionCodeRequest{
			Metadata: newMetadata,
		}

		updatedPromotionCode, err := stripeService.UpdatePromotionCode(ctx, createdPromotionCode.StripePromotionID, updateReq)

		require.NoError(t, err)
		assert.Equal(t, newMetadata, updatedPromotionCode.Metadata)
		assert.Equal(t, createdPromotionCode.Code, updatedPromotionCode.Code)
		assert.Equal(t, createdPromotionCode.Active, updatedPromotionCode.Active)
		assert.Equal(t, createdPromotionCode.FirstTimeTransaction, updatedPromotionCode.FirstTimeTransaction)
	})

	t.Run("deactivate active promotion code", func(t *testing.T) {
		// Create an active promotion code
		couponID := createMockCoupon()

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "DEACTIVATE",
			FirstTimeTransaction: false,
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)
		require.True(t, createdPromotionCode.Active)

		// Deactivate it
		newActive := false
		updateReq := &domain.UpdatePromotionCodeRequest{
			Active: &newActive,
		}

		updatedPromotionCode, err := stripeService.UpdatePromotionCode(ctx, createdPromotionCode.StripePromotionID, updateReq)

		require.NoError(t, err)
		assert.False(t, updatedPromotionCode.Active)
	})

	t.Run("reactivate deactivated promotion code", func(t *testing.T) {
		// Create a promotion code
		couponID := createMockCoupon()

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "REACTIVATE",
			FirstTimeTransaction: false,
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)

		// First deactivate it
		deactivate := false
		deactivateReq := &domain.UpdatePromotionCodeRequest{
			Active: &deactivate,
		}

		deactivatedPromotionCode, err := stripeService.UpdatePromotionCode(ctx, createdPromotionCode.StripePromotionID, deactivateReq)
		require.NoError(t, err)
		require.False(t, deactivatedPromotionCode.Active)

		// Now reactivate it
		reactivate := true
		reactivateReq := &domain.UpdatePromotionCodeRequest{
			Active: &reactivate,
		}

		reactivatedPromotionCode, err := stripeService.UpdatePromotionCode(ctx, createdPromotionCode.StripePromotionID, reactivateReq)

		require.NoError(t, err)
		assert.True(t, reactivatedPromotionCode.Active)
	})

	t.Run("update fails with non-existent promotion code ID", func(t *testing.T) {
		nonExistentID := "promo_nonexistent123456789"
		newActive := false

		updateReq := &domain.UpdatePromotionCodeRequest{
			Active: &newActive,
		}

		updatedPromotionCode, err := stripeService.UpdatePromotionCode(ctx, nonExistentID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, updatedPromotionCode)
	})

	t.Run("update with empty request does nothing", func(t *testing.T) {
		// Create a promotion code
		couponID := createMockCoupon()

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "EMPTYUPDATE",
			FirstTimeTransaction: true,
			Metadata: map[string]string{
				"original": "data",
			},
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)

		// Update with empty request
		updateReq := &domain.UpdatePromotionCodeRequest{}

		updatedPromotionCode, err := stripeService.UpdatePromotionCode(ctx, createdPromotionCode.StripePromotionID, updateReq)

		require.NoError(t, err)
		assert.Equal(t, createdPromotionCode.Code, updatedPromotionCode.Code)
		assert.Equal(t, createdPromotionCode.Active, updatedPromotionCode.Active)
		assert.Equal(t, createdPromotionCode.FirstTimeTransaction, updatedPromotionCode.FirstTimeTransaction)
		assert.Equal(t, createdPromotionCode.Metadata, updatedPromotionCode.Metadata)
	})

	t.Run("clear metadata with empty map", func(t *testing.T) {
		// Create a promotion code with metadata
		couponID := createMockCoupon()

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "CLEARMETA",
			FirstTimeTransaction: false,
			Metadata: map[string]string{
				"original": "metadata",
				"to":       "clear",
			},
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)
		require.NotEmpty(t, createdPromotionCode.Metadata)

		// Clear metadata with empty map
		updateReq := &domain.UpdatePromotionCodeRequest{
			Metadata: map[string]string{},
		}

		updatedPromotionCode, err := stripeService.UpdatePromotionCode(ctx, createdPromotionCode.StripePromotionID, updateReq)

		require.NoError(t, err)
		assert.Empty(t, updatedPromotionCode.Metadata)
		assert.Equal(t, createdPromotionCode.Code, updatedPromotionCode.Code)
		assert.Equal(t, createdPromotionCode.Active, updatedPromotionCode.Active)
	})

	t.Run("update promotion code with extensive metadata", func(t *testing.T) {
		// Create a promotion code
		couponID := createMockCoupon()

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "EXTENSIVEMETA",
			FirstTimeTransaction: false,
			Metadata: map[string]string{
				"version": "1.0",
			},
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)

		// Update with extensive metadata
		newMetadata := map[string]string{
			"version":      "2.0",
			"campaign":     "winter_sale",
			"target_group": "premium_users",
			"region":       "europe",
			"updated_by":   "marketing_api",
			"status":       "active",
		}

		updateReq := &domain.UpdatePromotionCodeRequest{
			Metadata: newMetadata,
		}

		updatedPromotionCode, err := stripeService.UpdatePromotionCode(ctx, createdPromotionCode.StripePromotionID, updateReq)

		require.NoError(t, err)
		assert.Equal(t, newMetadata, updatedPromotionCode.Metadata)
		assert.Len(t, updatedPromotionCode.Metadata, 6)
	})
}

func TestPromotionCodeActivationWorkflow(t *testing.T) {
	ctx := context.Background()

	createMockCoupon := func() string {
		return "coupon_test123456789"
	}

	t.Run("full activation workflow", func(t *testing.T) {
		// Create promotion code
		couponID := createMockCoupon()

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "WORKFLOW",
			FirstTimeTransaction: false,
			Metadata: map[string]string{
				"test": "workflow",
			},
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)
		assert.True(t, createdPromotionCode.Active)

		// Deactivate
		deactivate := false
		deactivateReq := &domain.UpdatePromotionCodeRequest{
			Active: &deactivate,
		}

		deactivatedPromotionCode, err := stripeService.UpdatePromotionCode(ctx, createdPromotionCode.StripePromotionID, deactivateReq)
		require.NoError(t, err)
		assert.False(t, deactivatedPromotionCode.Active)

		// Reactivate
		reactivate := true
		reactivateReq := &domain.UpdatePromotionCodeRequest{
			Active: &reactivate,
		}

		reactivatedPromotionCode, err := stripeService.UpdatePromotionCode(ctx, createdPromotionCode.StripePromotionID, reactivateReq)
		require.NoError(t, err)
		assert.True(t, reactivatedPromotionCode.Active)

		// Deactivate again with metadata update
		finalDeactivate := false
		finalMetadata := map[string]string{
			"test":   "workflow",
			"status": "permanently_disabled",
		}

		finalUpdateReq := &domain.UpdatePromotionCodeRequest{
			Active:   &finalDeactivate,
			Metadata: finalMetadata,
		}

		finalPromotionCode, err := stripeService.UpdatePromotionCode(ctx, createdPromotionCode.StripePromotionID, finalUpdateReq)
		require.NoError(t, err)
		assert.False(t, finalPromotionCode.Active)
		assert.Equal(t, finalMetadata, finalPromotionCode.Metadata)
	})
}