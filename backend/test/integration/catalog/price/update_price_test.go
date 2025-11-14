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

// make test-func TEST_NAME=TestUpdatePrice TEST_PATH=test/integration/catalog/price/update_price_test.go

func TestUpdatePrice(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully update price with valid admin token", func(t *testing.T) {
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup
		productIDStr := setupTestProduct(t, ctx)
		productID, err := uuid.Parse(productIDStr)
		require.NoError(t, err)

		// Create a price first
		price := td.NewValidPrice()
		price.ProductID = productID
		td.InsertPrice(t, ctx, price, testPool)

		// Execute - Update the price
		active := false
		nickname := "Updated Plan"
		updateRequest := &domain.UpdatePriceRequest{
			Active:   &active,
			Nickname: &nickname,
			Metadata: map[string]string{
				"updated_key": "updated_value",
				"plan_type":   "premium",
			},
		}

		req := th.NewUpdatePriceRequest(t, ctx, testServerURL, price.ID.String(), updateRequest, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updatedPrice domain.Price
		err = json.NewDecoder(resp.Body).Decode(&updatedPrice)
		assert.NoError(t, err)

		// Verify updated fields
		assert.Equal(t, price.ID, updatedPrice.ID)
		assert.False(t, updatedPrice.IsActive)

		// Verify database record is updated
		dbPrice := td.GetPriceByID(t, ctx, updatedPrice.ID, testPool)
		assert.False(t, dbPrice.IsActive)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)

		productIDStr := setupTestProduct(t, ctx)
		productID, err := uuid.Parse(productIDStr)
		require.NoError(t, err)

		price := td.NewValidPrice()
		price.ProductID = productID
		td.InsertPrice(t, ctx, price, testPool)

		active := false
		updateRequest := &domain.UpdatePriceRequest{
			Active: &active,
		}

		req := th.NewUpdatePriceRequest(t, ctx, testServerURL, price.ID.String(), updateRequest, "")

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

		// Create a price first
		productIDStr := setupTestProduct(t, ctx)
		productID, err := uuid.Parse(productIDStr)
		require.NoError(t, err)

		price := td.NewValidPrice()
		price.ProductID = productID
		td.InsertPrice(t, ctx, price, testPool)

		active := false
		updateRequest := &domain.UpdatePriceRequest{
			Active: &active,
		}

		req := th.NewUpdatePriceRequest(t, ctx, testServerURL, price.ID.String(), updateRequest, accessToken)

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

		// Create a price first
		productIDStr := setupTestProduct(t, ctx)
		productID, err := uuid.Parse(productIDStr)
		require.NoError(t, err)

		price := td.NewValidPrice()
		price.ProductID = productID
		td.InsertPrice(t, ctx, price, testPool)

		active := false
		updateRequest := &domain.UpdatePriceRequest{
			Active: &active,
		}

		req := th.NewUpdatePriceRequest(t, ctx, testServerURL, price.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)

		// Create a price first
		productIDStr := setupTestProduct(t, ctx)
		productID, err := uuid.Parse(productIDStr)
		require.NoError(t, err)

		price := td.NewValidPrice()
		price.ProductID = productID
		td.InsertPrice(t, ctx, price, testPool)

		active := false
		updateRequest := &domain.UpdatePriceRequest{
			Active: &active,
		}

		req := th.NewUpdatePriceRequest(t, ctx, testServerURL, price.ID.String(), updateRequest, "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
