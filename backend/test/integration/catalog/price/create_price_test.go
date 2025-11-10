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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreatePrice TEST_PATH=test/integration/catalog/price/create_price_test.go

func TestCreatePrice(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully create price with valid admin token", func(t *testing.T) {
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup
		productID := setupTestProduct(t, ctx)
		request := td.NewValidCreatePriceRequest()

		// Execute
		req := th.NewCreatePriceRequest(t, ctx, testServerURL, productID, request, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var createdPriceID string
		err = json.NewDecoder(resp.Body).Decode(&createdPriceID)
		assert.NoError(t, err)
		assert.NotEmpty(t, createdPriceID)

		// Get the created price to verify full structure
		getReq := th.NewGetPriceRequest(t, ctx, testServerURL, createdPriceID, accessToken)
		getResp, err := client.Do(getReq)
		require.NoError(t, err)
		defer getResp.Body.Close()
		assert.Equal(t, http.StatusOK, getResp.StatusCode)

		var createdPrice domain.Price
		err = json.NewDecoder(getResp.Body).Decode(&createdPrice)
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

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)

		productID := setupTestProduct(t, ctx)
		request := td.NewValidCreatePriceRequest()

		req := th.NewCreatePriceRequest(t, ctx, testServerURL, productID, request, "")

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
		request := td.NewValidCreatePriceRequest()

		req := th.NewCreatePriceRequest(t, ctx, testServerURL, productID, request, accessToken)

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
		request := td.NewValidCreatePriceRequest()

		req := th.NewCreatePriceRequest(t, ctx, testServerURL, productID, request, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)

		productID := setupTestProduct(t, ctx)
		request := td.NewValidCreatePriceRequest()

		req := th.NewCreatePriceRequest(t, ctx, testServerURL, productID, request, "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
