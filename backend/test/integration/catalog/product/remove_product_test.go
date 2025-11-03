package product_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestRemoveProduct_Integration TEST_PATH=test/integration/catalog/product/remove_product_test.go

func TestRemoveProduct_Integration(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - Product removed successfully", func(t *testing.T) {
		// Setup
		defer td.ClearProductsTable(t, ctx, testPool)
		defer td.ClearPricesTable(t, ctx, testPool)
		defer td.ClearCategoriesTable(t, ctx, testPool)

		// Create category first
		category := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, category, testPool)

		// Create product with Stripe IDs
		product := td.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = "prod_test123"
		td.InsertProduct(t, ctx, testPool, product)

		// Create prices for the product
		price1 := td.NewValidPrice()
		price1.ProductID = product.ID
		price1.StripePriceID = "price_test123"
		td.InsertPrice(t, ctx, price1, testPool)

		price2 := td.NewValidPrice()
		price2.ProductID = product.ID
		price2.StripePriceID = "price_test456"
		td.InsertPrice(t, ctx, price2, testPool)

		// Make request
		url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assertions
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify product is deleted from database
		_, err = td.GetProductByID(t, ctx, product.ID, testPool)
		assert.Error(t, err, "Product should be deleted from database")
	})

	t.Run("Error - Invalid URL path format", func(t *testing.T) {
		testCases := []struct {
			name               string
			path               string
			expectedStatusCode int
		}{
			{
				name:               "Missing admin segment",
				path:               "/products/123",
				expectedStatusCode: http.StatusMethodNotAllowed,
			},
			{
				name:               "Wrong admin segment",
				path:               "/user/products/123",
				expectedStatusCode: http.StatusNotFound,
			},
			{
				name:               "Missing products segment",
				path:               "/admin/items/123",
				expectedStatusCode: http.StatusNotFound,
			},
			{
				name:               "Too many segments",
				path:               "/admin/products/123/extra",
				expectedStatusCode: http.StatusNotFound,
			},
			{
				name:               "Too few segments",
				path:               "/admin/products",
				expectedStatusCode: http.StatusMethodNotAllowed,
			},
			{
				name:               "Empty path",
				path:               "",
				expectedStatusCode: http.StatusNotFound,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				url := fmt.Sprintf("%s%s", testServerURL, tc.path)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)

				client := &http.Client{Timeout: 10 * time.Second}
				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			})
		}
	})

	t.Run("Error - Missing product ID", func(t *testing.T) {
		url := fmt.Sprintf("%s/admin/products/", testServerURL)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Error - Invalid product ID format", func(t *testing.T) {
		invalidIDs := []string{
			"not-a-uuid",
			"123",
			"invalid-uuid-format",
			"12345678-1234-1234-1234",
		}

		for _, invalidID := range invalidIDs {
			t.Run(fmt.Sprintf("Invalid ID: %s", invalidID), func(t *testing.T) {
				url := fmt.Sprintf("%s/admin/products/%s", testServerURL, invalidID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)

				client := &http.Client{Timeout: 10 * time.Second}
				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			})
		}
	})

	t.Run("Error - Product not found", func(t *testing.T) {
		defer td.ClearProductsTable(t, ctx, testPool)

		nonExistentID := uuid.New()
		url := fmt.Sprintf("%s/admin/products/%s", testServerURL, nonExistentID)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Error - Product exists but no Stripe product ID", func(t *testing.T) {
		// Setup
		defer td.ClearProductsTable(t, ctx, testPool)
		defer td.ClearCategoriesTable(t, ctx, testPool)

		// Create category first
		category := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, category, testPool)

		// Create product without Stripe product ID (set to empty or null)
		product := td.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = "" // No Stripe ID
		td.InsertProduct(t, ctx, testPool, product)

		// Make request
		url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	})

	t.Run("Success - Product with no prices", func(t *testing.T) {
		// Setup
		defer td.ClearProductsTable(t, ctx, testPool)
		defer td.ClearCategoriesTable(t, ctx, testPool)

		// Create category first
		category := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, category, testPool)

		// Create product with Stripe ID but no prices
		product := td.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = "prod_noprices123"
		td.InsertProduct(t, ctx, testPool, product)

		// Make request
		url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should still succeed even with no prices
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify product is deleted
		_, err = td.GetProductByID(t, ctx, product.ID, testPool)
		assert.Error(t, err, "Product should be deleted from database")
	})

	t.Run("Success - Product with single price", func(t *testing.T) {
		// Setup
		defer td.ClearProductsTable(t, ctx, testPool)
		defer td.ClearPricesTable(t, ctx, testPool)
		defer td.ClearCategoriesTable(t, ctx, testPool)

		// Create category first
		category := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, category, testPool)

		// Create product
		product := td.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = "prod_singleprice123"
		td.InsertProduct(t, ctx, testPool, product)

		// Create single price
		price := td.NewValidPrice()
		price.ProductID = product.ID
		price.StripePriceID = "price_single123"
		td.InsertPrice(t, ctx, price, testPool)

		// Make request
		url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify product is deleted
		_, err = td.GetProductByID(t, ctx, product.ID, testPool)
		assert.Error(t, err, "Product should be deleted from database")
	})

	t.Run("Success - Product with multiple prices", func(t *testing.T) {
		// Setup
		defer td.ClearProductsTable(t, ctx, testPool)
		defer td.ClearPricesTable(t, ctx, testPool)
		defer td.ClearCategoriesTable(t, ctx, testPool)

		// Create category first
		category := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, category, testPool)

		// Create product
		product := td.NewValidProduct("Test Product", category.ID)
		product.StripeProductID = "prod_multiprices123"
		td.InsertProduct(t, ctx, testPool, product)

		// Create multiple prices
		prices := make([]*domain.Price, 3)
		for i := range 3 {
			price := td.NewValidPrice()
			price.ProductID = product.ID
			price.StripePriceID = fmt.Sprintf("price_multi_%d", i)
			td.InsertPrice(t, ctx, price, testPool)
			prices[i] = price
		}

		// Make request
		url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify product is deleted
		_, err = td.GetProductByID(t, ctx, product.ID, testPool)
		assert.Error(t, err, "Product should be deleted from database")
	})
}

// make test-func TEST_NAME=TestRemoveProduct_HTTPMethod TEST_PATH=test/integration/catalog/product/remove_product_test.go

// Test for HTTP method validation (if your router supports it)
func TestRemoveProduct_HTTPMethod(t *testing.T) {
	ctx := context.Background()

	// Setup a valid product for testing different HTTP methods
	defer td.ClearProductsTable(t, ctx, testPool)
	defer td.ClearCategoriesTable(t, ctx, testPool)

	category := td.NewValidCategory("Test Category")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Test Product", category.ID)
	product.StripeProductID = "prod_method_test"
	td.InsertProduct(t, ctx, testPool, product)

	url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID)

	// Test unsupported HTTP methods
	unsupportedMethods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodHead,
		http.MethodOptions,
	}

	for _, method := range unsupportedMethods {
		t.Run(fmt.Sprintf("Method_%s", method), func(t *testing.T) {
			req, err := http.NewRequest(method, url, nil)
			require.NoError(t, err)

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Expecting 405 Method Not Allowed or similar
			assert.NotEqual(t, http.StatusNoContent, resp.StatusCode)
		})
	}
}

// make test-func TEST_NAME=TestRemoveProduct_ConcurrentDeletion TEST_PATH=test/integration/catalog/product/remove_product_test.go

// Test concurrent deletion attempts
func TestRemoveProduct_ConcurrentDeletion(t *testing.T) {
	ctx := context.Background()

	// Setup
	defer td.ClearProductsTable(t, ctx, testPool)
	defer td.ClearCategoriesTable(t, ctx, testPool)

	category := td.NewValidCategory("Test Category")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Test Product", category.ID)
	product.StripeProductID = "prod_concurrent_test"
	td.InsertProduct(t, ctx, testPool, product)

	url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID)

	// Launch multiple concurrent deletion requests
	const numRequests = 3
	responses := make(chan *http.Response, numRequests)
	errors := make(chan error, numRequests)

	for range numRequests {
		go func() {
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			if err != nil {
				errors <- err
				return
			}

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
	_, err := td.GetProductByID(t, ctx, product.ID, testPool)
	assert.Error(t, err, "Product should be deleted from database")
}
