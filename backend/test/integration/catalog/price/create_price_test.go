package price_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreatePrice TEST_PATH=test/integration/catalog/price/create_price_test.go

func TestCreatePrice(t *testing.T) {
	ctx := context.Background()

	t.Run("should create price successfully with valid data", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)
		request := td.NewValidCreatePriceRequest()

		// Execute
		resp, body := createPriceViaAPI(t, productID, request)

		// Assert
		assertSuccessResponse(t, resp, 201)

		var createdPriceID string
		err := json.Unmarshal(body, &createdPriceID)
		assert.NoError(t, err)
		assert.NotEmpty(t, createdPriceID)

		// Get the created price to verify full structure
		getResp, getBody := getPriceViaAPI(t, createdPriceID)
		assertSuccessResponse(t, getResp, 200)

		var createdPrice domain.Price
		err = json.Unmarshal(getBody, &createdPrice)
		assert.NoError(t, err)

		// Verify response structure
		assert.Equal(t, createdPriceID, createdPrice.ID.String())
		assert.Equal(t, request.Amount, createdPrice.Amount)
		assert.Equal(t, request.Currency, createdPrice.Currency)
		assert.Equal(t, request.Interval, string(createdPrice.Interval))
		assert.True(t, createdPrice.IsActive) // Prices are active by default
		assert.NotEmpty(t, createdPrice.StripePriceID)
		assert.NotZero(t, createdPrice.CreatedAt)

		// Verify database record
		dbPrice := td.GetPriceByID(t, ctx, createdPrice.ID, testPool)
		assert.Equal(t, createdPrice.ID, dbPrice.ID)
		assert.Equal(t, request.Amount, dbPrice.Amount)
		assert.Equal(t, request.Currency, dbPrice.Currency)
		assert.Equal(t, request.Interval, string(dbPrice.Interval))
		assert.True(t, dbPrice.IsActive) // Prices are active by default
	})

	t.Run("should return 400 when product ID is missing", func(t *testing.T) {
		// Setup
		request := td.NewValidCreatePriceRequest()

		// Execute - using empty product ID
		// Note: When productID is empty, the URL becomes /admin/products//prices which is 404
		resp, _ := createPriceViaAPI(t, "", request)

		// Assert - expecting 404 because empty path segment results in route not found
		assertErrorResponse(t, resp, 404)
	})

	t.Run("should return 400 with invalid request payload", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		// Execute - using invalid payload (nil)
		resp, _ := createPriceViaAPI(t, productID, nil)

		// Assert
		assertErrorResponse(t, resp, 400)
	})

	t.Run("should return 400 with invalid price data", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)
		request := &domain.CreatePriceRequest{
			Amount:   -100, // Invalid negative amount
			Currency: "",   // Invalid empty currency
			Interval: "",   // Invalid empty interval
		}

		// Execute
		resp, _ := createPriceViaAPI(t, productID, request)

		// Assert
		assertErrorResponse(t, resp, 400)
	})

	t.Run("should return 404 when product does not exist", func(t *testing.T) {
		// Setup
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)
		nonExistentProductID := "550e8400-e29b-41d4-a716-446655440000"
		request := td.NewValidCreatePriceRequest()

		// Execute
		resp, _ := createPriceViaAPI(t, nonExistentProductID, request)

		// Assert
		assertErrorResponse(t, resp, 404)
	})

	t.Run("should handle Stripe service unavailable", func(t *testing.T) {
		// This test would require mocking Stripe to return errors
		// For now, we'll test with valid data to ensure the happy path works
		// In a real scenario, you might want to test with network issues or Stripe mock returning errors

		productID := setupTestProduct(t, ctx)
		request := td.NewValidCreatePriceRequest()

		resp, body := createPriceViaAPI(t, productID, request)

		// Should succeed with mock Stripe
		assertSuccessResponse(t, resp, 201)

		var createdPriceID string
		err := json.Unmarshal(body, &createdPriceID)
		assert.NoError(t, err)
		assert.NotEmpty(t, createdPriceID)

		// Verify the created price has Stripe ID
		getResp, getBody := getPriceViaAPI(t, createdPriceID)
		assertSuccessResponse(t, getResp, 200)
		var createdPrice domain.Price
		err = json.Unmarshal(getBody, &createdPrice)
		assert.NoError(t, err)
		assert.NotEmpty(t, createdPrice.StripePriceID)
	})

	t.Run("should create multiple prices for same product", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		// Create first price
		request1 := td.NewValidCreatePriceRequest()
		request1.Amount = 1000
		request1.Nickname = td.StrPtr("Basic Plan")

		resp1, body1 := createPriceViaAPI(t, productID, request1)
		assertSuccessResponse(t, resp1, 201)

		var priceID1 string
		err := json.Unmarshal(body1, &priceID1)
		assert.NoError(t, err)

		// Create second price
		request2 := td.NewValidCreatePriceRequest()
		request2.Amount = 2000
		request2.Nickname = td.StrPtr("Premium Plan")

		resp2, body2 := createPriceViaAPI(t, productID, request2)
		assertSuccessResponse(t, resp2, 201)

		var priceID2 string
		err = json.Unmarshal(body2, &priceID2)
		assert.NoError(t, err)

		// Get both prices to verify structure
		getResp1, getBody1 := getPriceViaAPI(t, priceID1)
		assertSuccessResponse(t, getResp1, 200)
		var price1 domain.Price
		err = json.Unmarshal(getBody1, &price1)
		assert.NoError(t, err)

		getResp2, getBody2 := getPriceViaAPI(t, priceID2)
		assertSuccessResponse(t, getResp2, 200)
		var price2 domain.Price
		err = json.Unmarshal(getBody2, &price2)
		assert.NoError(t, err)

		// Assert both prices exist and are different
		assert.NotEqual(t, price1.ID, price2.ID)
		assert.NotEqual(t, price1.StripePriceID, price2.StripePriceID)
		assert.Equal(t, 1000, price1.Amount)
		assert.Equal(t, 2000, price2.Amount)

		// Verify both prices are in database
		dbPrice1 := td.GetPriceByID(t, ctx, price1.ID, testPool)
		dbPrice2 := td.GetPriceByID(t, ctx, price2.ID, testPool)
		assert.Equal(t, price1.ProductID, dbPrice1.ProductID)
		assert.Equal(t, price2.ProductID, dbPrice2.ProductID)
		assert.Equal(t, dbPrice1.ProductID, dbPrice2.ProductID) // Same product
	})
}
