package pricePayment_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePrice(t *testing.T) {
	ctx := context.Background()

	// Helper function to create a mock product first
	createMockProduct := func() string {
		// This assumes we have a way to create a product via the Stripe mock
		// In a real scenario, you might need to use the product service or mock directly
		return "prod_test123456789" // Mock product ID
	}

	t.Run("successful price creation", func(t *testing.T) {
		productID := createMockProduct()

		req := domain.CreateStripePriceRequest{
			ProductID: productID,
			Amount:    2999, // $29.99
			Currency:  "USD",
			Interval:  "month",
			Active:    true,
			Nickname:  "Monthly Plan",
			Metadata: map[string]string{
				"plan_type": "monthly",
				"tier":      "premium",
			},
		}

		price, err := stripeService.CreatePrice(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, price)
		assert.NotEmpty(t, price.ID)
		assert.Equal(t, productID, price.Product)
		assert.Equal(t, req.Amount, price.Amount)
		assert.Equal(t, req.Currency, price.Currency)
		assert.Equal(t, req.Interval, price.Interval)
		assert.Equal(t, req.Active, price.Active)
		assert.Equal(t, req.Nickname, price.Nickname)
		assert.Equal(t, req.Metadata, price.Metadata)
	})

	t.Run("create price with minimal fields", func(t *testing.T) {
		productID := createMockProduct()

		req := domain.CreateStripePriceRequest{
			ProductID: productID,
			Amount:    1000, // $10.00
			Currency:  "USD",
			Interval:  "month",
			Active:    true,
		}

		price, err := stripeService.CreatePrice(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, price)
		assert.Equal(t, req.Amount, price.Amount)
		assert.Equal(t, req.Currency, price.Currency)
		assert.Equal(t, req.Interval, price.Interval)
		assert.True(t, price.Active)
	})

	t.Run("create inactive price", func(t *testing.T) {
		productID := createMockProduct()

		req := domain.CreateStripePriceRequest{
			ProductID: productID,
			Amount:    5000, // $50.00
			Currency:  "USD",
			Interval:  "year",
			Active:    false,
			Nickname:  "Yearly Plan (Inactive)",
		}

		price, err := stripeService.CreatePrice(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, price)
		assert.False(t, price.Active)
		assert.Equal(t, "Yearly Plan (Inactive)", price.Nickname)
	})

	t.Run("create price with different currencies", func(t *testing.T) {
		testCases := []struct {
			currency string
			amount   int
		}{
			{"EUR", 2500},
			{"GBP", 2000},
			{"jpy", 300000}, // Yen doesn't use decimal places
		}

		for _, tc := range testCases {
			t.Run("currency_"+tc.currency, func(t *testing.T) {
				productID := createMockProduct()

				req := domain.CreateStripePriceRequest{
					ProductID: productID,
					Amount:    tc.amount,
					Currency:  tc.currency,
					Interval:  "month",
					Active:    true,
				}

				price, err := stripeService.CreatePrice(ctx, req)

				require.NoError(t, err)
				assert.Equal(t, tc.currency, price.Currency)
				assert.Equal(t, tc.amount, price.Amount)
			})
		}
	})

	t.Run("create price with different intervals", func(t *testing.T) {
		intervals := []string{"day", "week", "month", "year"}

		for _, interval := range intervals {
			t.Run("interval_"+interval, func(t *testing.T) {
				productID := createMockProduct()

				req := domain.CreateStripePriceRequest{
					ProductID: productID,
					Amount:    1000,
					Currency:  "USD",
					Interval:  interval,
					Active:    true,
				}

				price, err := stripeService.CreatePrice(ctx, req)

				require.NoError(t, err)
				assert.Equal(t, interval, price.Interval)
			})
		}
	})
}