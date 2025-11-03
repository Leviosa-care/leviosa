package price_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	priceHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/price"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/require"
)

// makeRequest is a helper function to send HTTP requests during tests
func makeRequest(t *testing.T, method, endpoint string, body interface{}) (*http.Response, []byte) {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	fullURL := testServerURL + endpoint

	req, err := http.NewRequest(method, fullURL, reqBody)
	require.NoError(t, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	responseBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp, responseBody
}

// buildEndpointWithID replaces the {id} placeholder with the actual ID
func buildEndpointWithID(template, id string) string {
	return strings.Replace(template, "{id}", id, 1)
}

// setupTestProduct creates a product for price testing and returns its ID
func setupTestProduct(t *testing.T, ctx context.Context) string {
	t.Helper()

	// Clear tables first
	td.ClearProductsTable(t, ctx, testPool)
	td.ClearPricesTable(t, ctx, testPool)

	// Clear categories table to avoid constraint violations
	td.ClearCategoriesTable(t, ctx, testPool)

	// Create a test category first (required for product)
	category := td.NewValidCategory("test category")
	td.InsertCategory(t, ctx, category, testPool)

	// Insert a test product
	product := td.NewValidProduct("Test Product", category.ID)
	td.InsertProduct(t, ctx, testPool, product)

	return product.ID.String()
}

// createPriceViaAPI is a helper to create a price through the HTTP API
func createPriceViaAPI(t *testing.T, productID string, request *domain.CreatePriceRequest) (*http.Response, []byte) {
	t.Helper()
	endpoint := buildEndpointWithID(priceHandler.CreatePriceEndpoint, productID)
	return makeRequest(t, "POST", endpoint, request)
}

// getPriceViaAPI is a helper to get a price through the HTTP API
func getPriceViaAPI(t *testing.T, priceID string) (*http.Response, []byte) {
	t.Helper()
	endpoint := buildEndpointWithID(priceHandler.GetPriceEndpoint, priceID)
	return makeRequest(t, "GET", endpoint, nil)
}

// getPricesByProductIDViaAPI is a helper to get prices by product ID through the HTTP API
func getPricesByProductIDViaAPI(t *testing.T, productID string) (*http.Response, []byte) {
	t.Helper()
	endpoint := buildEndpointWithID(priceHandler.GetPricesByProductIDEndpoint, productID)
	return makeRequest(t, "GET", endpoint, nil)
}

// updatePriceViaAPI is a helper to update a price through the HTTP API
func updatePriceViaAPI(t *testing.T, priceID string, request *domain.UpdatePriceRequest) (*http.Response, []byte) {
	t.Helper()
	endpoint := buildEndpointWithID(priceHandler.UpdatePriceEndpoint, priceID)
	return makeRequest(t, "PATCH", endpoint, request)
}

// assertErrorResponse checks that the response is an error with expected status code
func assertErrorResponse(t *testing.T, resp *http.Response, expectedStatus int) {
	t.Helper()
	require.Equal(t, expectedStatus, resp.StatusCode, "Expected status code %d, got %d", expectedStatus, resp.StatusCode)

	// Check content type - both text/plain (route not found) and application/json (handler errors) are valid for 404s
	contentType := resp.Header.Get("Content-Type")
	if expectedStatus == 404 {
		// 404s can be either route not found (text/plain) or handler not found (application/json)
		require.True(t, contentType == "text/plain; charset=utf-8" || contentType == "application/json",
			"Expected text/plain or application/json for 404, got %s", contentType)
	} else {
		require.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	}
}

// assertSuccessResponse checks that the response is successful with expected status code
func assertSuccessResponse(t *testing.T, resp *http.Response, expectedStatus int) {
	t.Helper()
	require.Equal(t, expectedStatus, resp.StatusCode, "Expected status code %d, got %d", expectedStatus, resp.StatusCode)
	require.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}
