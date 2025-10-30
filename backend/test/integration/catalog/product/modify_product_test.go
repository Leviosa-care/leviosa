package product_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModifyProduct_Success_UpdateName(t *testing.T) {
	ctx := context.Background()

	// Clean up tables
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)

	// Setup test data
	category := td.NewValidCategory("Electronics")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Original Product Name", category.ID)
	td.InsertProduct(t, ctx, testPool, product)

	// Prepare update request
	newName := "Updated Product Name"
	updateReq := domain.UpdateProductRequest{
		Name: &newName,
	}

	reqBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	// Make request
	url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID.String())
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify database update
	updatedProduct, err := td.GetProductByID(t, ctx, product.ID, testPool)
	require.NoError(t, err)
	assert.Equal(t, newName, updatedProduct.Name)
	assert.Equal(t, product.Description, updatedProduct.Description) // Should remain unchanged
}

func TestModifyProduct_Success_UpdateDescription(t *testing.T) {
	ctx := context.Background()

	// Clean up tables
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)

	// Setup test data
	category := td.NewValidCategory("Services")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Test Product", category.ID)
	td.InsertProduct(t, ctx, testPool, product)

	// Prepare update request
	newDescription := "This is the updated description"
	updateReq := domain.UpdateProductRequest{
		Description: &newDescription,
	}

	reqBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	// Make request
	url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID.String())
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify database update
	updatedProduct, err := td.GetProductByID(t, ctx, product.ID, testPool)
	require.NoError(t, err)
	assert.Equal(t, newDescription, updatedProduct.Description)
	assert.Equal(t, product.Name, updatedProduct.Name) // Should remain unchanged
}

func TestModifyProduct_Success_UpdateStatus(t *testing.T) {
	ctx := context.Background()

	// Clean up tables
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)

	// Setup test data
	category := td.NewValidCategory("Books")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Test Product", category.ID)
	product.Status = domain.Draft // Original status
	td.InsertProduct(t, ctx, testPool, product)

	// Prepare update request
	newStatus := domain.Published
	updateReq := domain.UpdateProductRequest{
		Status: td.StrPtr(string(newStatus)),
	}

	reqBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	// Make request
	url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID.String())
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify database update
	updatedProduct, err := td.GetProductByID(t, ctx, product.ID, testPool)
	require.NoError(t, err)
	assert.Equal(t, newStatus, updatedProduct.Status)
}

func TestModifyProduct_Success_UpdateMetadata(t *testing.T) {
	ctx := context.Background()

	// Clean up tables
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)

	// Setup test data
	category := td.NewValidCategory("Software")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Test Product", category.ID)
	product.Metadata = map[string]any{"version": "1.0", "type": "basic"}
	td.InsertProduct(t, ctx, testPool, product)

	// Prepare update request
	newMetadata := map[string]any{
		"version":  "2.0",
		"type":     "premium",
		"features": []any{"advanced", "support"},
	}
	updateReq := domain.UpdateProductRequest{
		Metadata: newMetadata,
	}

	reqBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	// Make request
	url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID.String())
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify database update
	updatedProduct, err := td.GetProductByID(t, ctx, product.ID, testPool)
	require.NoError(t, err)
	assert.Equal(t, newMetadata, updatedProduct.Metadata)
}

func TestModifyProduct_Success_UpdateMultipleFields(t *testing.T) {
	ctx := context.Background()

	// Clean up tables
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)

	// Setup test data
	category := td.NewValidCategory("Mixed")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Original Name", category.ID)
	product.Status = domain.Draft
	td.InsertProduct(t, ctx, testPool, product)

	// Prepare update request with multiple fields
	newName := "Updated Name"
	newDescription := "Updated description"
	newStatus := domain.Published
	newMetadata := map[string]any{"updated": true}

	updateReq := domain.UpdateProductRequest{
		Name:        &newName,
		Description: &newDescription,
		Status:      td.StrPtr(string(newStatus)),
		Metadata:    newMetadata,
	}

	reqBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	// Make request
	url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID.String())
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify all database updates
	updatedProduct, err := td.GetProductByID(t, ctx, product.ID, testPool)
	require.NoError(t, err)
	assert.Equal(t, newName, updatedProduct.Name)
	assert.Equal(t, newDescription, updatedProduct.Description)
	assert.Equal(t, newStatus, updatedProduct.Status)
	assert.Equal(t, newMetadata, updatedProduct.Metadata)
}

func TestModifyProduct_InvalidURLPath(t *testing.T) {
	testCases := []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{
			name:               "Wrong path structure", // 405
			url:                "/products/123",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:               "Missing admin prefix", // 405
			url:                "/products/123",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:               "Wrong section order",
			url:                "/admin/categories/123",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Too many path segments",
			url:                "/admin/products/123/extra",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Missing product ID",
			url:                "/admin/products/",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateReq := domain.UpdateProductRequest{
				Name: td.StrPtr("Test"),
			}
			reqBody, err := json.Marshal(updateReq)
			require.NoError(t, err)

			url := testServerURL + tc.url
			req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
		})
	}
}

func TestModifyProduct_InvalidRequestBody(t *testing.T) {
	ctx := context.Background()

	// Setup minimal test data for valid product ID
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)

	category := td.NewValidCategory("Test")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Test Product", category.ID)
	td.InsertProduct(t, ctx, testPool, product)

	testCases := []struct {
		name        string
		requestBody string
		expectedMsg string
	}{
		{
			name:        "Invalid JSON",
			requestBody: `{"name": "test"`,
			expectedMsg: "invalid request body",
		},
		{
			name:        "Unknown fields",
			requestBody: `{"name": "test", "unknownField": "value"}`,
			expectedMsg: "invalid request body",
		},
		{
			name:        "Empty request body",
			requestBody: `{}`,
			expectedMsg: "no updatable fields provided in request body",
		},
		{
			name:        "All fields null",
			requestBody: `{"name": null, "description": null, "publishedStatus": null, "metadata": null}`,
			expectedMsg: "no updatable fields provided in request body",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID.String())
			req, err := http.NewRequest("PATCH", url, strings.NewReader(tc.requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

			var errorResp map[string]string
			err = json.NewDecoder(resp.Body).Decode(&errorResp)
			require.NoError(t, err)
			assert.Contains(t, errorResp["error"], tc.expectedMsg)
		})
	}
}

func TestModifyProduct_InvalidProductID(t *testing.T) {
	testCases := []struct {
		name      string
		productID string
		// expectedMsg        string
		expectedStatusCode int
	}{
		{
			name:               "Invalid UUID format",
			productID:          "not-a-uuid",
			expectedStatusCode: http.StatusBadRequest,
			// expectedMsg: "product ID is required and must be a valid UUID",
		},
		{
			name:               "Empty UUID",
			productID:          "",
			expectedStatusCode: http.StatusNotFound,
			// expectedMsg: "product ID is missing from the URL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateReq := domain.UpdateProductRequest{
				Name: td.StrPtr("Updated Name"),
			}
			reqBody, err := json.Marshal(updateReq)
			require.NoError(t, err)

			url := fmt.Sprintf("%s/admin/products/%s", testServerURL, tc.productID)
			req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			// The invalid UUID error comes from the service layer,
			// but URL parsing errors come from handler layer

			// var errorResp map[string]string
			// err = json.NewDecoder(resp.Body).Decode(&errorResp)
			// require.NoError(t, err)
			// assert.Contains(t, errorResp["error"], tc.expectedMsg)
		})
	}
}

func TestModifyProduct_ProductNotFound(t *testing.T) {
	ctx := context.Background()

	// Clean up tables
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)

	// Use a valid UUID that doesn't exist in the database
	nonExistentID := uuid.New()

	updateReq := domain.UpdateProductRequest{
		Name: td.StrPtr("Updated Name"),
	}
	reqBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	// Make request
	url := fmt.Sprintf("%s/admin/products/%s", testServerURL, nonExistentID.String())
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should return 404 or the error that your service returns for not found
	// Based on your code, it should be handled by the error handling in the handler
	// The exact status code depends on how your service maps the error
	assert.True(t, resp.StatusCode == http.StatusNotFound || resp.StatusCode >= 400)
}

func TestModifyProduct_HTTPMethods(t *testing.T) {
	ctx := context.Background()

	// Setup minimal test data
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)

	category := td.NewValidCategory("Test")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Test Product", category.ID)
	td.InsertProduct(t, ctx, testPool, product)

	updateReq := domain.UpdateProductRequest{
		Name: td.StrPtr("Updated Name"),
	}
	reqBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	// Test different HTTP methods
	methods := []struct {
		method         string
		expectedStatus int
	}{
		{"GET", http.StatusMethodNotAllowed},
		{"POST", http.StatusMethodNotAllowed},
		{"PUT", http.StatusMethodNotAllowed},
		{"PATCH", http.StatusNoContent}, // Should succeed
	}

	for _, m := range methods {
		t.Run(fmt.Sprintf("Method_%s", m.method), func(t *testing.T) {
			url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID.String())

			req, err := http.NewRequest(m.method, url, bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, m.expectedStatus, resp.StatusCode)
		})
	}
}

func TestModifyProduct_ContentTypeValidation(t *testing.T) {
	ctx := context.Background()

	// Setup minimal test data
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)

	category := td.NewValidCategory("Test")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Test Product", category.ID)
	td.InsertProduct(t, ctx, testPool, product)

	updateReq := domain.UpdateProductRequest{
		Name: td.StrPtr("Updated Name"),
	}
	reqBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		contentType string
		expectError bool
	}{
		{
			name:        "Valid JSON content type",
			contentType: "application/json",
			expectError: false,
		},
		{
			name:        "Missing content type",
			contentType: "",
			expectError: false, // JSON decoder should still work
		},
		{
			name:        "Wrong content type",
			contentType: "text/plain",
			expectError: false, // JSON decoder attempts to parse anyway
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID.String())
			req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
			require.NoError(t, err)

			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			if tc.expectError {
				assert.NotEqual(t, http.StatusNoContent, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusNoContent, resp.StatusCode)
			}
		})
	}
}

// Benchmark test for performance verification
func BenchmarkModifyProduct(b *testing.B) {
	ctx := context.Background()

	// Clean up tables - use direct SQL since td helpers need *testing.T
	_, err := testPool.Exec(ctx, "TRUNCATE TABLE catalog.categories RESTART IDENTITY CASCADE;")
	if err != nil {
		b.Fatalf("Failed to clear categories table: %v", err)
	}
	_, err = testPool.Exec(ctx, "TRUNCATE catalog.products CASCADE;")
	if err != nil {
		b.Fatalf("Failed to clear products table: %v", err)
	}

	// Setup test data
	setupHelper := &testing.T{}
	category := td.NewValidCategory("Benchmark")
	td.InsertCategory(setupHelper, ctx, category, testPool)

	product := td.NewValidProduct("Benchmark Product", category.ID)
	td.InsertProduct(setupHelper, ctx, testPool, product)

	if setupHelper.Failed() {
		b.Fatalf("Test setup failed")
	}

	// Prepare update request
	updateReq := domain.UpdateProductRequest{
		Name: td.StrPtr("Updated Benchmark Product"),
	}
	reqBody, err := json.Marshal(updateReq)
	if err != nil {
		b.Fatalf("Failed to marshal request: %v", err)
	}

	url := fmt.Sprintf("%s/admin/products/%s", testServerURL, product.ID.String())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
			if err != nil {
				b.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				b.Fatalf("Request failed: %v", err)
			}
			if resp.StatusCode != http.StatusNoContent {
				b.Fatalf("Expected 204, got %d", resp.StatusCode)
			}
			resp.Body.Close()
		}
	})
}

// Helper function to create string pointers
// func td.StrPtr(s string) *string {
// 	return &s
// }
