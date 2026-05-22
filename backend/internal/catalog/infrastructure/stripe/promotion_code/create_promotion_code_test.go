package promotionCodePayment_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePromotionCode(t *testing.T) {
	ctx := context.Background()

	// Helper function to create a mock coupon first (since promotion codes require a coupon)
	createMockCoupon := func() string {
		// This assumes we have a way to create a coupon via the Stripe mock
		// In a real scenario, you might need to use the coupon service or mock directly
		return "coupon_test123456789" // Mock coupon ID
	}

	t.Run("successful promotion code creation", func(t *testing.T) {
		couponID := createMockCoupon()
		maxRedemptions := 50
		expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days from now
		minimumAmount := 2000                            // $20.00
		minimumAmountCurrency := "USD"

		req := &domain.CreatePromotionCodeRequest{
			CouponID:              couponID,
			Code:                  "SUMMER25",
			MaxRedemptions:        &maxRedemptions,
			ExpiresAt:             &expiresAt,
			FirstTimeTransaction:  true,
			MinimumAmount:         &minimumAmount,
			MinimumAmountCurrency: &minimumAmountCurrency,
			Restrictions: &domain.PromotionCodeRestrictionsRequest{
				CurrencyOptions: []string{"USD", "EUR"},
			},
			Metadata: map[string]string{
				"campaign": "summer_sale",
				"type":     "limited",
			},
		}

		promotionCode, err := stripeService.CreatePromotionCode(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, promotionCode)
		assert.NotEmpty(t, promotionCode.StripePromotionID)
		assert.Equal(t, req.Code, promotionCode.Code)
		assert.True(t, promotionCode.Active)
		assert.Equal(t, req.FirstTimeTransaction, promotionCode.FirstTimeTransaction)
		assert.Equal(t, &maxRedemptions, promotionCode.MaxRedemptions)
		assert.NotNil(t, promotionCode.ExpiresAt)
		assert.Equal(t, &minimumAmount, promotionCode.MinimumAmount)
		assert.Equal(t, &minimumAmountCurrency, promotionCode.MinimumAmountCurrency)
		assert.Equal(t, req.Restrictions, promotionCode.Restrictions)
		assert.Equal(t, req.Metadata, promotionCode.Metadata)
		assert.Equal(t, 0, promotionCode.TimesRedeemed)
		assert.False(t, promotionCode.CreatedAt.IsZero())
	})

	t.Run("create promotion code with minimal fields", func(t *testing.T) {
		couponID := createMockCoupon()

		req := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "SIMPLE10",
			FirstTimeTransaction: false,
		}

		promotionCode, err := stripeService.CreatePromotionCode(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, promotionCode)
		assert.Equal(t, req.Code, promotionCode.Code)
		assert.True(t, promotionCode.Active)
		assert.Equal(t, req.FirstTimeTransaction, promotionCode.FirstTimeTransaction)
		assert.Nil(t, promotionCode.MaxRedemptions)
		assert.Nil(t, promotionCode.ExpiresAt)
		assert.Nil(t, promotionCode.MinimumAmount)
		assert.Nil(t, promotionCode.Restrictions)
	})

	t.Run("create promotion code for first-time customers only", func(t *testing.T) {
		couponID := createMockCoupon()
		maxRedemptions := 100

		req := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "WELCOME20",
			MaxRedemptions:       &maxRedemptions,
			FirstTimeTransaction: true,
			Metadata: map[string]string{
				"type":   "welcome",
				"target": "new_customers",
			},
		}

		promotionCode, err := stripeService.CreatePromotionCode(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, promotionCode)
		assert.Equal(t, "WELCOME20", promotionCode.Code)
		assert.True(t, promotionCode.FirstTimeTransaction)
		assert.Equal(t, &maxRedemptions, promotionCode.MaxRedemptions)
		assert.Equal(t, req.Metadata, promotionCode.Metadata)
	})

	t.Run("create promotion code with expiry date", func(t *testing.T) {
		couponID := createMockCoupon()
		expiresAt := time.Now().Add(7 * 24 * time.Hour) // 1 week from now

		req := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "FLASH7DAYS",
			ExpiresAt:            &expiresAt,
			FirstTimeTransaction: false,
			Metadata: map[string]string{
				"promotion": "flash_sale",
				"duration":  "7_days",
			},
		}

		promotionCode, err := stripeService.CreatePromotionCode(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, promotionCode)
		assert.NotNil(t, promotionCode.ExpiresAt)
		// Allow some tolerance for time comparison (within 1 minute)
		assert.WithinDuration(t, expiresAt, *promotionCode.ExpiresAt, time.Minute)
	})

	t.Run("create promotion code with minimum amount", func(t *testing.T) {
		couponID := createMockCoupon()
		minimumAmount := 5000 // $50.00
		currency := "USD"

		req := &domain.CreatePromotionCodeRequest{
			CouponID:              couponID,
			Code:                  "MIN50USD",
			FirstTimeTransaction:  false,
			MinimumAmount:         &minimumAmount,
			MinimumAmountCurrency: &currency,
		}

		promotionCode, err := stripeService.CreatePromotionCode(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, promotionCode)
		assert.Equal(t, &minimumAmount, promotionCode.MinimumAmount)
		assert.Equal(t, &currency, promotionCode.MinimumAmountCurrency)
	})

	t.Run("create promotion code with currency restrictions", func(t *testing.T) {
		couponID := createMockCoupon()
		currencyOptions := []string{"EUR", "GBP", "cad"}

		req := &domain.CreatePromotionCodeRequest{
			CouponID:             couponID,
			Code:                 "INTL15",
			FirstTimeTransaction: false,
			Restrictions: &domain.PromotionCodeRestrictionsRequest{
				CurrencyOptions: currencyOptions,
			},
			Metadata: map[string]string{
				"type":   "international",
				"region": "europe_uk_canada",
			},
		}

		promotionCode, err := stripeService.CreatePromotionCode(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, promotionCode)
		assert.NotNil(t, promotionCode.Restrictions)
		assert.Equal(t, currencyOptions, promotionCode.Restrictions.CurrencyOptions)
		assert.Equal(t, req.Metadata, promotionCode.Metadata)
	})

	t.Run("create promotion code with different currencies for minimum amount", func(t *testing.T) {
		testCases := []struct {
			currency string
			amount   int
		}{
			{"EUR", 4000}, // €40.00
			{"GBP", 3500}, // £35.00
			{"cad", 6000}, // CA$60.00
		}

		for _, tc := range testCases {
			t.Run("currency_"+tc.currency, func(t *testing.T) {
				couponID := createMockCoupon()

				req := &domain.CreatePromotionCodeRequest{
					CouponID:              couponID,
					Code:                  "MIN" + tc.currency,
					FirstTimeTransaction:  false,
					MinimumAmount:         &tc.amount,
					MinimumAmountCurrency: &tc.currency,
				}

				promotionCode, err := stripeService.CreatePromotionCode(ctx, req)

				require.NoError(t, err)
				assert.Equal(t, &tc.currency, promotionCode.MinimumAmountCurrency)
				assert.Equal(t, &tc.amount, promotionCode.MinimumAmount)
			})
		}
	})

	t.Run("create promotion code with limited redemptions", func(t *testing.T) {
		testCases := []int{1, 10, 50, 100, 1000}

		for _, maxRedemptions := range testCases {
			t.Run("max_redemptions_"+string(rune(maxRedemptions)), func(t *testing.T) {
				couponID := createMockCoupon()

				req := &domain.CreatePromotionCodeRequest{
					CouponID:             couponID,
					Code:                 "LIMITED" + string(rune(maxRedemptions)),
					FirstTimeTransaction: false,
					MaxRedemptions:       &maxRedemptions,
				}

				promotionCode, err := stripeService.CreatePromotionCode(ctx, req)

				require.NoError(t, err)
				assert.Equal(t, &maxRedemptions, promotionCode.MaxRedemptions)
				assert.Equal(t, 0, promotionCode.TimesRedeemed)
			})
		}
	})
}
