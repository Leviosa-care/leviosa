package price_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME='^TestGetPricesByProductID$$' TEST_PATH=test/integration/catalog/price/get_prices_by_product_id_test.go

func TestGetPricesByProductID(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get prices by product ID with valid admin token", func(t *testing.T) {
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup
		productIDStr := setupTestProduct(t, ctx)
		productID, err := uuid.Parse(productIDStr)
		require.NoError(t, err)

		const monthAmount = 100
		const annualAmount = 1200
		const currency = "EUR"

		// Create multiple prices for the same product
		createdPrice1 := td.NewValidPrice()
		createdPrice1.ProductID = productID
		createdPrice1.Amount = monthAmount
		createdPrice1.Interval = domain.Month

		createdPrice2 := td.NewValidPrice()
		createdPrice2.ProductID = productID
		createdPrice2.Amount = annualAmount
		createdPrice2.Interval = domain.Year

		td.InsertPrice(t, ctx, createdPrice1, testPool)
		td.InsertPrice(t, ctx, createdPrice2, testPool)

		// Execute
		req := th.NewGetPricesByProductIDRequest(t, ctx, testServerURL, productID.String(), accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var prices []*domain.Price
		err = json.NewDecoder(resp.Body).Decode(&prices)
		assert.NoError(t, err)

		// Verify we got both prices
		assert.Len(t, prices, 2)

		// Sort prices by amount for consistent comparison
		if prices[0].Amount > prices[1].Amount {
			prices[0], prices[1] = prices[1], prices[0]
		}

		// Verify first price (Basic Plan)
		assert.Equal(t, monthAmount, prices[0].Amount)
		assert.Equal(t, domain.Month, prices[0].Interval)
		assert.Equal(t, currency, prices[0].Currency)
		assert.True(t, prices[0].IsActive)

		// Verify second price (Premium Plan)
		assert.Equal(t, annualAmount, prices[1].Amount)
		assert.Equal(t, domain.Year, prices[1].Interval)
		assert.Equal(t, currency, prices[1].Currency)
		assert.True(t, prices[1].IsActive)

		// Verify both prices belong to the same product
		assert.Equal(t, prices[0].ProductID, prices[1].ProductID)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)

		productID := setupTestProduct(t, ctx)

		req := th.NewGetPricesByProductIDRequest(t, ctx, testServerURL, productID, "")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		productID := setupTestProduct(t, ctx)

		req := th.NewGetPricesByProductIDRequest(t, ctx, testServerURL, productID, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		productID := setupTestProduct(t, ctx)

		req := th.NewGetPricesByProductIDRequest(t, ctx, testServerURL, productID, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)

		productID := setupTestProduct(t, ctx)

		req := th.NewGetPricesByProductIDRequest(t, ctx, testServerURL, productID, "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
