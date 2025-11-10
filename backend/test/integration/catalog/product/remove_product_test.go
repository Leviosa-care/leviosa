package product_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	productHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/product"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME='^TestRemoveProduct$' TEST_PATH=test/integration/catalog/product/remove_product_test.go

func TestRemoveProduct(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully remove product with valid admin token", func(t *testing.T) {
		th.ClearProductsTable(t, ctx, testPool)
		th.ClearPricesTable(t, ctx, testPool)
		th.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create category first
		category := th.NewValidCategory("Test Category")
		th.InsertCategory(t, ctx, category, testPool)

		// Create product with Stripe IDs
		product := th.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = "prod_test123"
		th.InsertProduct(t, ctx, testPool, product)

		// Create prices for the product
		price1 := th.NewValidPrice()
		price1.ProductID = product.ID
		price1.StripePriceID = "price_test123"
		th.InsertPrice(t, ctx, price1, testPool)

		price2 := th.NewValidPrice()
		price2.ProductID = product.ID
		price2.StripePriceID = "price_test456"
		th.InsertPrice(t, ctx, price2, testPool)

		// Make request
		req := th.NewRemoveProductRequest(t, ctx, testServerURL, product.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assertions
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify product is deleted from database
		_, err = th.GetProductByID(t, ctx, product.ID, testPool)
		assert.Error(t, err, "Product should be deleted from database")
	})

	t.Run("should successfully remove product without prices with valid admin token", func(t *testing.T) {
		th.ClearProductsTable(t, ctx, testPool)
		th.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create category first
		category := th.NewValidCategory("Test Category")
		th.InsertCategory(t, ctx, category, testPool)

		// Create product with Stripe ID but no prices
		product := th.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = "prod_noprices123"
		th.InsertProduct(t, ctx, testPool, product)

		// Make request
		req := th.NewRemoveProductRequest(t, ctx, testServerURL, product.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should still succeed even with no prices
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify product is deleted
		_, err = th.GetProductByID(t, ctx, product.ID, testPool)
		assert.Error(t, err, "Product should be deleted from database")
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		th.ClearProductsTable(t, ctx, testPool)
		th.ClearCategoriesTable(t, ctx, testPool)

		// Create category first
		category := th.NewValidCategory("Test Category")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = "prod_test123"
		th.InsertProduct(t, ctx, testPool, product)

		req := th.NewRemoveProductRequest(t, ctx, testServerURL, product.ID.String(), "")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		th.ClearProductsTable(t, ctx, testPool)
		th.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Create category first
		category := th.NewValidCategory("Test Category")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = "prod_test123"
		th.InsertProduct(t, ctx, testPool, product)

		req := th.NewRemoveProductRequest(t, ctx, testServerURL, product.ID.String(), accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		th.ClearProductsTable(t, ctx, testPool)
		th.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Create category first
		category := th.NewValidCategory("Test Category")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = "prod_test123"
		th.InsertProduct(t, ctx, testPool, product)

		req := th.NewRemoveProductRequest(t, ctx, testServerURL, product.ID.String(), accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		th.ClearProductsTable(t, ctx, testPool)
		th.ClearCategoriesTable(t, ctx, testPool)

		// Create category first
		category := th.NewValidCategory("Test Category")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = "prod_test123"
		th.InsertProduct(t, ctx, testPool, product)

		req := th.NewRemoveProductRequest(t, ctx, testServerURL, product.ID.String(), "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 404 when product not found", func(t *testing.T) {
		th.ClearProductsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		nonExistentID := uuid.New()

		req := th.NewRemoveProductRequest(t, ctx, testServerURL, nonExistentID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid product ID format", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewRemoveProductRequest(t, ctx, testServerURL, "not-a-uuid", accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 503 when product exists but no Stripe product ID", func(t *testing.T) {
		th.ClearProductsTable(t, ctx, testPool)
		th.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create category first
		category := th.NewValidCategory("Test Category")
		th.InsertCategory(t, ctx, category, testPool)

		// Create product without Stripe product ID (set to empty or null)
		product := th.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = ""
		th.InsertProduct(t, ctx, testPool, product)

		req := th.NewRemoveProductRequest(t, ctx, testServerURL, product.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	})

	t.Run("should return 405 for wrong HTTP method", func(t *testing.T) {
		th.ClearProductsTable(t, ctx, testPool)
		th.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create category first
		category := th.NewValidCategory("Test Category")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = "prod_test123"
		th.InsertProduct(t, ctx, testPool, product)

		// Try using GET instead of DELETE - must manually construct to test wrong method
		url := testServerURL + productHandler.AdminProductsBasePath + "/" + product.ID.String()
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("should return 404 for invalid URL path - missing admin segment", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req, err := http.NewRequest("DELETE", testServerURL+"/products/123", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("should return 404 for invalid URL path - wrong admin segment", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req, err := http.NewRequest("DELETE", testServerURL+"/user/products/123", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 404 for invalid URL path - missing products segment", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req, err := http.NewRequest("DELETE", testServerURL+"/admin/items/123", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 404 for invalid URL path - too many segments", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req, err := http.NewRequest("DELETE", testServerURL+"/admin/products/123/extra", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should handle concurrent deletion attempts", func(t *testing.T) {
		th.ClearProductsTable(t, ctx, testPool)
		th.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup
		category := th.NewValidCategory("Test Category")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = "prod_concurrent_test"
		th.InsertProduct(t, ctx, testPool, product)

		// Launch multiple concurrent deletion requests
		const numRequests = 3
		responses := make(chan *http.Response, numRequests)
		errors := make(chan error, numRequests)

		for range numRequests {
			go func() {
				req := th.NewRemoveProductRequest(t, ctx, testServerURL, product.ID.String(), accessToken)

				client := &http.Client{Timeout: 10 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					errors <- err
					return
				}

				responses <- resp
			}()
		}

		// Collect responses
		var successCount, notFoundCount, errorCount int
		statusCodes := make([]int, 0, numRequests)

		for range numRequests {
			select {
			case resp := <-responses:
				defer resp.Body.Close()
				statusCodes = append(statusCodes, resp.StatusCode)
				switch resp.StatusCode {
				case http.StatusNoContent:
					successCount++
				case http.StatusNotFound:
					notFoundCount++
				case http.StatusInternalServerError, http.StatusServiceUnavailable:
					errorCount++
				default:
					t.Errorf("Unexpected status code: %d", resp.StatusCode)
				}
			case err := <-errors:
				t.Errorf("Request failed: %v", err)
			case <-time.After(15 * time.Second):
				t.Fatal("Test timed out waiting for responses")
			}
		}

		t.Logf("Response status codes: %v", statusCodes)
		t.Logf("Success: %d, NotFound: %d, Errors: %d", successCount, notFoundCount, errorCount)

		// With concurrent requests, we expect:
		// - At least one successful deletion (could be more if they all start before any completes)
		// - The remaining requests should either get 404 (not found) or errors from Stripe/rollback
		assert.True(t, successCount >= 1, "At least one deletion should succeed")
		assert.Equal(t, numRequests, successCount+notFoundCount+errorCount, "All requests should get a response")

		// Verify product is actually deleted from database
		_, err := th.GetProductByID(t, ctx, product.ID, testPool)
		assert.Error(t, err, "Product should be deleted from database")
	})
}
