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

// make test-func TEST_NAME='^TestGetPrice$$' TEST_PATH=test/integration/catalog/price/get_price_test.go

func TestGetPrice(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get price with valid admin token", func(t *testing.T) {
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup
		productIDStr := setupTestProduct(t, ctx)
		productID, err := uuid.Parse(productIDStr)
		require.NoError(t, err)

		price := td.NewValidPrice()
		price.ProductID = productID
		td.InsertPrice(t, ctx, price, testPool)

		// Execute
		req := th.NewGetPriceRequest(t, ctx, testServerURL, price.ID.String(), accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var retrievedPrice domain.Price
		err = json.NewDecoder(resp.Body).Decode(&retrievedPrice)
		assert.NoError(t, err)

		// Verify response structure
		assert.Equal(t, price.ID.String(), retrievedPrice.ID.String())
		assert.Equal(t, price.Amount, retrievedPrice.Amount)
		assert.Equal(t, price.Currency, retrievedPrice.Currency)
		assert.Equal(t, price.Interval, retrievedPrice.Interval)
		assert.True(t, retrievedPrice.IsActive) // Should be active by default
		assert.NotZero(t, retrievedPrice.CreatedAt)
		assert.NotZero(t, retrievedPrice.UpdatedAt)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		td.ClearPricesTable(t, ctx, testPool)

		priceID := "550e8400-e29b-41d4-a716-446655440000"

		req := th.NewGetPriceRequest(t, ctx, testServerURL, priceID, "")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		td.ClearPricesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		priceID := "550e8400-e29b-41d4-a716-446655440000"

		req := th.NewGetPriceRequest(t, ctx, testServerURL, priceID, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		td.ClearPricesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		priceID := "550e8400-e29b-41d4-a716-446655440000"

		req := th.NewGetPriceRequest(t, ctx, testServerURL, priceID, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		td.ClearPricesTable(t, ctx, testPool)

		priceID := "550e8400-e29b-41d4-a716-446655440000"

		req := th.NewGetPriceRequest(t, ctx, testServerURL, priceID, "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
