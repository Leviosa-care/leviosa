package couponPayment_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCoupon(t *testing.T) {
	ctx := context.Background()

	t.Run("successful percent-off coupon creation", func(t *testing.T) {
		percentOff := 25.0
		maxRedemptions := 100

		req := &domain.CreateCouponRequest{
			Name:           "25% Off Coupon",
			PercentOff:     &percentOff,
			Duration:       "once",
			MaxRedemptions: &maxRedemptions,
			Metadata: map[string]string{
				"campaign": "summer_sale",
				"type":     "percent",
			},
		}

		coupon, err := stripeService.CreateCoupon(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, coupon)
		assert.NotEmpty(t, coupon.StripeCouponID)
		assert.Equal(t, req.Name, coupon.Name)
		assert.Equal(t, &percentOff, coupon.PercentOff)
		assert.Nil(t, coupon.AmountOff)
		assert.Nil(t, coupon.Currency)
		assert.Equal(t, domain.CouponDuration(req.Duration), coupon.Duration)
		assert.Equal(t, &maxRedemptions, coupon.MaxRedemptions)
		assert.True(t, coupon.IsValid)
		assert.Equal(t, 0, coupon.TimesRedeemed)
		assert.Equal(t, req.Metadata, coupon.Metadata)
		assert.False(t, coupon.CreatedAt.IsZero())
	})

	t.Run("successful amount-off coupon creation", func(t *testing.T) {
		amountOff := 500 // $5.00
		currency := "USD"
		redeemBy := time.Now().Add(30 * 24 * time.Hour) // 30 days from now

		req := &domain.CreateCouponRequest{
			Name:      "5 Dollar Off Coupon",
			AmountOff: &amountOff,
			Currency:  &currency,
			Duration:  "once",
			RedeemBy:  &redeemBy,
			Metadata: map[string]string{
				"campaign": "new_user",
				"type":     "fixed",
			},
		}

		coupon, err := stripeService.CreateCoupon(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, coupon)
		assert.NotEmpty(t, coupon.StripeCouponID)
		assert.Equal(t, req.Name, coupon.Name)
		assert.Nil(t, coupon.PercentOff)
		assert.Equal(t, &amountOff, coupon.AmountOff)
		assert.Equal(t, &currency, coupon.Currency)
		assert.Equal(t, domain.CouponDuration(req.Duration), coupon.Duration)
		assert.NotNil(t, coupon.RedeemBy)
		assert.True(t, coupon.IsValid)
		assert.Equal(t, req.Metadata, coupon.Metadata)
	})

	t.Run("successful repeating coupon creation", func(t *testing.T) {
		percentOff := 15.0
		durationInMonths := 3

		req := &domain.CreateCouponRequest{
			Name:             "15% Off for 3 Months",
			PercentOff:       &percentOff,
			Duration:         "repeating",
			DurationInMonths: &durationInMonths,
			Metadata: map[string]string{
				"type": "repeating",
			},
		}

		coupon, err := stripeService.CreateCoupon(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, coupon)
		assert.Equal(t, req.Name, coupon.Name)
		assert.Equal(t, &percentOff, coupon.PercentOff)
		assert.Equal(t, domain.CouponDuration("repeating"), coupon.Duration)
		assert.Equal(t, &durationInMonths, coupon.DurationInMonths)
		assert.True(t, coupon.IsValid)
	})

	t.Run("successful forever coupon creation", func(t *testing.T) {
		percentOff := 10.0
		maxRedemptions := 1000

		req := &domain.CreateCouponRequest{
			Name:           "Forever 10% Off",
			PercentOff:     &percentOff,
			Duration:       "forever",
			MaxRedemptions: &maxRedemptions,
		}

		coupon, err := stripeService.CreateCoupon(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, coupon)
		assert.Equal(t, req.Name, coupon.Name)
		assert.Equal(t, &percentOff, coupon.PercentOff)
		assert.Equal(t, domain.CouponDuration("forever"), coupon.Duration)
		assert.Equal(t, &maxRedemptions, coupon.MaxRedemptions)
		assert.Nil(t, coupon.DurationInMonths)
		assert.True(t, coupon.IsValid)
	})

	t.Run("successful coupon creation with minimal fields", func(t *testing.T) {
		percentOff := 20.0

		req := &domain.CreateCouponRequest{
			Name:       "Minimal Coupon",
			PercentOff: &percentOff,
			Duration:   "once",
		}

		coupon, err := stripeService.CreateCoupon(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, coupon)
		assert.Equal(t, req.Name, coupon.Name)
		assert.Equal(t, &percentOff, coupon.PercentOff)
		assert.Equal(t, domain.CouponDuration("once"), coupon.Duration)
		assert.True(t, coupon.IsValid)
		assert.Equal(t, 0, coupon.TimesRedeemed)
	})

	t.Run("create coupon with different currencies", func(t *testing.T) {
		testCases := []struct {
			currency  string
			amountOff int
		}{
			{"EUR", 1000}, // €10.00
			{"GBP", 750},  // £7.50
			{"cad", 1250}, // CA$12.50
		}

		for _, tc := range testCases {
			t.Run("currency_"+tc.currency, func(t *testing.T) {
				req := &domain.CreateCouponRequest{
					Name:      "Amount Off " + tc.currency,
					AmountOff: &tc.amountOff,
					Currency:  &tc.currency,
					Duration:  "once",
				}

				coupon, err := stripeService.CreateCoupon(ctx, req)

				require.NoError(t, err)
				assert.Equal(t, &tc.currency, coupon.Currency)
				assert.Equal(t, &tc.amountOff, coupon.AmountOff)
			})
		}
	})

	t.Run("create coupon with expiry date", func(t *testing.T) {
		percentOff := 30.0
		redeemBy := time.Now().Add(7 * 24 * time.Hour) // 1 week from now

		req := &domain.CreateCouponRequest{
			Name:       "Limited Time 30% Off",
			PercentOff: &percentOff,
			Duration:   "once",
			RedeemBy:   &redeemBy,
			Metadata: map[string]string{
				"promotion": "flash_sale",
			},
		}

		coupon, err := stripeService.CreateCoupon(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, coupon)
		assert.NotNil(t, coupon.RedeemBy)
		// Allow some tolerance for time comparison (within 1 minute)
		assert.WithinDuration(t, redeemBy, *coupon.RedeemBy, time.Minute)
		assert.Equal(t, req.Metadata, coupon.Metadata)
	})
}