package promotionCodePayment_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPromotionCode(t *testing.T) {
	ctx := context.Background()

	createMockCoupon := func() string {
		return "coupon_test123456789"
	}

	t.Run("successful promotion code retrieval", func(t *testing.T) {
		// First create a promotion code
		couponID := createMockCoupon()
		maxRedemptions := 75
		minimumAmount := 1500 // $15.00
		minimumAmountCurrency := "USD"

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:              couponID,
			Code:                  "TESTGET25",
			MaxRedemptions:        &maxRedemptions,
			FirstTimeTransaction:  true,
			MinimumAmount:         &minimumAmount,
			MinimumAmountCurrency: &minimumAmountCurrency,
			Metadata: map[string]string{
				"test": "get_promotion_code",
			},
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)
		require.NotNil(t, createdPromotionCode)

		// Now retrieve the promotion code
		retrievedPromotionCode, err := stripeService.GetPromotionCode(ctx, createdPromotionCode.StripePromotionID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedPromotionCode)
		assert.Equal(t, createdPromotionCode.StripePromotionID, retrievedPromotionCode.StripePromotionID)
		assert.Equal(t, createdPromotionCode.Code, retrievedPromotionCode.Code)
		assert.Equal(t, createdPromotionCode.Active, retrievedPromotionCode.Active)
		assert.Equal(t, createdPromotionCode.FirstTimeTransaction, retrievedPromotionCode.FirstTimeTransaction)
		assert.Equal(t, createdPromotionCode.MaxRedemptions, retrievedPromotionCode.MaxRedemptions)
		assert.Equal(t, createdPromotionCode.MinimumAmount, retrievedPromotionCode.MinimumAmount)
		assert.Equal(t, createdPromotionCode.MinimumAmountCurrency, retrievedPromotionCode.MinimumAmountCurrency)
		assert.Equal(t, createdPromotionCode.Metadata, retrievedPromotionCode.Metadata)
		assert.Equal(t, createdPromotionCode.TimesRedeemed, retrievedPromotionCode.TimesRedeemed)
	})

	t.Run("retrieve promotion code with minimal data", func(t *testing.T) {
		// Create a minimal promotion code
		couponID := createMockCoupon()

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "MINIMAL",
			FirstTimeTransaction: false,
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)

		// Retrieve it
		retrievedPromotionCode, err := stripeService.GetPromotionCode(ctx, createdPromotionCode.StripePromotionID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedPromotionCode)
		assert.Equal(t, createdPromotionCode.StripePromotionID, retrievedPromotionCode.StripePromotionID)
		assert.Equal(t, "MINIMAL", retrievedPromotionCode.Code)
		assert.True(t, retrievedPromotionCode.Active)
		assert.False(t, retrievedPromotionCode.FirstTimeTransaction)
		assert.Nil(t, retrievedPromotionCode.MaxRedemptions)
		assert.Nil(t, retrievedPromotionCode.ExpiresAt)
		assert.Nil(t, retrievedPromotionCode.MinimumAmount)
		assert.Equal(t, 0, retrievedPromotionCode.TimesRedeemed)
	})

	t.Run("retrieve promotion code with expiry", func(t *testing.T) {
		// Create a promotion code with expiry
		couponID := createMockCoupon()
		expiresAt := time.Now().Add(14 * 24 * time.Hour) // 2 weeks from now

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "EXPIRES2WK",
			FirstTimeTransaction: false,
			ExpiresAt:            &expiresAt,
			Metadata: map[string]string{
				"type": "limited_time",
			},
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)

		// Retrieve it
		retrievedPromotionCode, err := stripeService.GetPromotionCode(ctx, createdPromotionCode.StripePromotionID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedPromotionCode)
		assert.NotNil(t, retrievedPromotionCode.ExpiresAt)
		assert.WithinDuration(t, expiresAt, *retrievedPromotionCode.ExpiresAt, time.Minute)
		assert.Equal(t, createdPromotionCode.Metadata, retrievedPromotionCode.Metadata)
	})

	t.Run("retrieve promotion code with currency restrictions", func(t *testing.T) {
		// Create a promotion code with currency restrictions
		couponID := createMockCoupon()
		currencyOptions := []string{"EUR", "GBP"}

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "EUROPE20",
			FirstTimeTransaction: false,
			Restrictions: &domain.PromotionCodeRestrictions{
				CurrencyOptions: currencyOptions,
			},
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)

		// Retrieve it
		retrievedPromotionCode, err := stripeService.GetPromotionCode(ctx, createdPromotionCode.StripePromotionID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedPromotionCode)
		assert.NotNil(t, retrievedPromotionCode.Restrictions)
		assert.Equal(t, currencyOptions, retrievedPromotionCode.Restrictions.CurrencyOptions)
	})

	t.Run("retrieve first-time customer promotion code", func(t *testing.T) {
		// Create a first-time customer promotion code
		couponID := createMockCoupon()
		maxRedemptions := 200

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "FIRSTTIME30",
			FirstTimeTransaction: true,
			MaxRedemptions:       &maxRedemptions,
			Metadata: map[string]string{
				"target":   "new_customers",
				"discount": "30_percent",
			},
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)

		// Retrieve it
		retrievedPromotionCode, err := stripeService.GetPromotionCode(ctx, createdPromotionCode.StripePromotionID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedPromotionCode)
		assert.True(t, retrievedPromotionCode.FirstTimeTransaction)
		assert.Equal(t, &maxRedemptions, retrievedPromotionCode.MaxRedemptions)
		assert.Equal(t, createdPromotionCode.Metadata, retrievedPromotionCode.Metadata)
	})

	t.Run("retrieve promotion code with minimum amount in different currency", func(t *testing.T) {
		// Create a promotion code with EUR minimum amount
		couponID := createMockCoupon()
		minimumAmount := 2500 // €25.00
		currency := "EUR"

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:              couponID,
			Code:                  "EUR25MIN",
			FirstTimeTransaction:  false,
			MinimumAmount:         &minimumAmount,
			MinimumAmountCurrency: &currency,
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)

		// Retrieve it
		retrievedPromotionCode, err := stripeService.GetPromotionCode(ctx, createdPromotionCode.StripePromotionID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedPromotionCode)
		assert.Equal(t, &minimumAmount, retrievedPromotionCode.MinimumAmount)
		assert.Equal(t, &currency, retrievedPromotionCode.MinimumAmountCurrency)
	})

	t.Run("retrieval fails with non-existent promotion code ID", func(t *testing.T) {
		nonExistentID := "promo_nonexistent123456789"

		promotionCode, err := stripeService.GetPromotionCode(ctx, nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, promotionCode)
	})

	t.Run("retrieve promotion code with comprehensive metadata", func(t *testing.T) {
		// Create a promotion code with extensive metadata
		couponID := createMockCoupon()

		createReq := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "METADATA",
			FirstTimeTransaction: false,
			Metadata: map[string]string{
				"campaign":     "holiday_season",
				"type":         "promotional",
				"target_group": "vip_customers",
				"region":       "north_america",
				"created_by":   "marketing_team",
			},
		}

		createdPromotionCode, err := stripeService.CreatePromotionCode(ctx, createReq)
		require.NoError(t, err)

		// Retrieve it
		retrievedPromotionCode, err := stripeService.GetPromotionCode(ctx, createdPromotionCode.StripePromotionID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedPromotionCode)
		assert.Equal(t, createdPromotionCode.Metadata, retrievedPromotionCode.Metadata)
		assert.Len(t, retrievedPromotionCode.Metadata, 5)
	})
}
