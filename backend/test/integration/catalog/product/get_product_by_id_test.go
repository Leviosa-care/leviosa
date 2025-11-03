package product_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetProductByID_Success TEST_PATH=test/integration/catalog/product/get_product_by_id_test.go

func TestGetProductByID_Success(t *testing.T) {
	ctx := context.Background()

	// Clean up tables
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)
	td.ClearImagesTable(t, ctx, testPool)
	td.ClearPricesTable(t, ctx, testPool)

	// Setup test data
	category := td.NewValidCategory("Electronics")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Test Product", category.ID)
	td.InsertProduct(t, ctx, testPool, product)

	// Create an active image for the product
	image := td.NewValidImage(product.ID)
	image.ParentType = domain.ProductType
	image.IsActive = true
	td.InsertImage(t, ctx, image, testPool)

	// Create active prices for the product
	price1 := td.NewValidPrice()
	price1.ProductID = product.ID
	price1.Amount = 1000 // $10.00
	price1.Currency = "USD"
	price1.Interval = "month"
	price1.IsActive = true
	td.InsertPrice(t, ctx, price1, testPool)

	price2 := td.NewValidPrice()
	price2.ProductID = product.ID
	price2.Amount = 10000 // $100.00
	price2.Currency = "USD"
	price2.Interval = "year"
	price2.IsActive = true
	td.InsertPrice(t, ctx, price2, testPool)

	// Make request
	url := fmt.Sprintf("%s/products/%s", testServerURL, product.ID.String())
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var result domain.ProductAggregator
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	// Verify product data
	assert.Equal(t, product.ID, result.Product.ID)
	assert.Equal(t, product.Name, result.Product.Name)
	assert.Equal(t, product.Description, result.Product.Description)
	assert.Equal(t, category.ID, result.Product.Category.ID)

	// Verify category data is populated
	assert.Equal(t, category.Name, result.Product.Category.Name)

	// Verify image data
	assert.NotNil(t, result.Image)
	assert.Equal(t, image.ID, result.Image.ID)
	assert.Equal(t, product.ID, result.Image.ParentID)
	assert.Equal(t, domain.ProductType, result.Image.ParentType)
	assert.True(t, result.Image.IsActive)

	// Verify prices data
	assert.Len(t, result.Prices, 2)
	priceIDs := []uuid.UUID{result.Prices[0].ID, result.Prices[1].ID}
	assert.Contains(t, priceIDs, price1.ID)
	assert.Contains(t, priceIDs, price2.ID)

	// Verify all prices are active
	for _, price := range result.Prices {
		assert.True(t, price.IsActive)
		assert.Equal(t, product.ID, price.ProductID)
	}
}

// make test-func TEST_NAME=TestGetProductByID_SuccessWithoutImage TEST_PATH=test/integration/catalog/product/get_product_by_id_test.go

func TestGetProductByID_SuccessWithoutImage(t *testing.T) {
	ctx := context.Background()

	// Clean up tables
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)
	td.ClearImagesTable(t, ctx, testPool)
	td.ClearPricesTable(t, ctx, testPool)

	// Setup test data - no image this time
	category := td.NewValidCategory("Services")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Service Product", category.ID)
	td.InsertProduct(t, ctx, testPool, product)

	// Create prices but no image
	price := td.NewValidPrice()
	price.ProductID = product.ID
	td.InsertPrice(t, ctx, price, testPool)

	// Make request
	url := fmt.Sprintf("%s/products/%s", testServerURL, product.ID.String())
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result domain.ProductAggregator
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	// Verify product data is still present
	assert.Equal(t, product.ID, result.Product.ID)
	assert.Equal(t, product.Name, result.Product.Name)

	// Verify image is nil (no active image found)
	assert.Nil(t, result.Image)

	// Verify prices are still present
	assert.Len(t, result.Prices, 1)
	assert.Equal(t, price.ID, result.Prices[0].ID)
}

// make test-func TEST_NAME=TestGetProductByID_SuccessWithInactiveImage TEST_PATH=test/integration/catalog/product/get_product_by_id_test.go

func TestGetProductByID_SuccessWithInactiveImage(t *testing.T) {
	ctx := context.Background()

	// Clean up tables
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)
	td.ClearImagesTable(t, ctx, testPool)
	td.ClearPricesTable(t, ctx, testPool)

	// Setup test data
	category := td.NewValidCategory("Books")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Book Product", category.ID)
	td.InsertProduct(t, ctx, testPool, product)

	// Create an inactive image
	image := td.NewValidImage(product.ID)
	image.ParentType = domain.ProductType
	image.IsActive = false // This should not be returned
	td.InsertImage(t, ctx, image, testPool)

	// Create price
	price := td.NewValidPrice()
	price.ProductID = product.ID
	td.InsertPrice(t, ctx, price, testPool)

	// Make request
	url := fmt.Sprintf("%s/products/%s", testServerURL, product.ID.String())
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result domain.ProductAggregator
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	// Verify image is nil because it's inactive
	assert.Nil(t, result.Image)

	// Verify other data is present
	assert.Equal(t, product.ID, result.Product.ID)
	assert.Len(t, result.Prices, 1)
}

// make test-func TEST_NAME=TestGetProductByID_InvalidProductID TEST_PATH=test/integration/catalog/product/get_product_by_id_test.go

func TestGetProductByID_InvalidProductID(t *testing.T) {
	testCases := []struct {
		name           string
		productID      string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Empty product ID",
			productID:      "",
			expectedStatus: http.StatusNotFound,
			expectedError:  "",
		},
		{
			name:           "Product ID with slash",
			productID:      "123/456",
			expectedStatus: http.StatusNotFound,
			expectedError:  "",
		},
		{
			name:           "Invalid UUID format",
			productID:      "not-a-uuid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  errs.ErrInvalidValue.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/products/%s", testServerURL, tc.productID)
			resp, err := http.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedError != "" {
				var errorResp struct {
					Error string `json:"error"`
				}
				err = json.NewDecoder(resp.Body).Decode(&errorResp)
				assert.NoError(t, err)

				assert.Contains(t, errorResp.Error, tc.expectedError)
			}
		})
	}
}

// make test-func TEST_NAME=TestGetProductByID_ProductNotFound TEST_PATH=test/integration/catalog/product/get_product_by_id_test.go

func TestGetProductByID_ProductNotFound(t *testing.T) {
	ctx := context.Background()

	// Clean up tables
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)
	td.ClearImagesTable(t, ctx, testPool)
	td.ClearPricesTable(t, ctx, testPool)

	// Use a valid UUID that doesn't exist in the database
	nonExistentID := uuid.New()

	// Make request
	url := fmt.Sprintf("%s/products/%s", testServerURL, nonExistentID.String())
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var errorResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	assert.NoError(t, err)

	assert.Contains(t, errorResp["error"].(string), fmt.Sprintf("product with ID %s", nonExistentID))
}

// make test-func TEST_NAME=TestGetProductByID_CategoryForeignKeyConstraint TEST_PATH=test/integration/catalog/product/get_product_by_id_test.go

func TestGetProductByID_CategoryForeignKeyConstraint(t *testing.T) {
	ctx := context.Background()

	// Clean up tables
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)
	td.ClearImagesTable(t, ctx, testPool)
	td.ClearPricesTable(t, ctx, testPool)

	// Try to insert a product with a non-existent category ID
	// This should fail due to foreign key constraint
	nonExistentCategoryID := uuid.New()
	product := td.NewValidProduct("Invalid Product", nonExistentCategoryID)

	// This should fail due to foreign key constraint
	query := `
	INSERT INTO catalog.products (
		id, name, description, category_id, duration,
		created_at, updated_at, status, availability,
		buffer_time, cancellation_hours, stripe_product_id, metadata
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	)`

	_, err := testPool.Exec(ctx, query,
		product.ID, product.Name, product.Description, product.CategoryID,
		product.Duration, product.CreatedAt, product.UpdatedAt, product.Status,
		product.Availability, product.BufferTime, product.CancellationHours,
		product.StripeProductID, []byte("{}"),
	)

	// Should fail with foreign key violation
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "foreign key")
}

// make test-func TEST_NAME=TestGetProductByID_WithMultiplePrices TEST_PATH=test/integration/catalog/product/get_product_by_id_test.go

func TestGetProductByID_WithMultiplePrices(t *testing.T) {
	ctx := context.Background()

	// Clean up tables
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)
	td.ClearImagesTable(t, ctx, testPool)
	td.ClearPricesTable(t, ctx, testPool)

	// Setup test data
	category := td.NewValidCategory("Software")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Software Product", category.ID)
	td.InsertProduct(t, ctx, testPool, product)

	// Create multiple prices with different statuses
	activePrice1 := td.NewValidPrice()
	activePrice1.ProductID = product.ID
	activePrice1.Amount = 500
	activePrice1.Interval = "month"
	activePrice1.IsActive = true
	activePrice1.CreatedAt = time.Now().Add(-2 * time.Hour) // Older
	td.InsertPrice(t, ctx, activePrice1, testPool)

	activePrice2 := td.NewValidPrice()
	activePrice2.ProductID = product.ID
	activePrice2.Amount = 5000
	activePrice2.Interval = "year"
	activePrice2.IsActive = true
	activePrice2.CreatedAt = time.Now().Add(-1 * time.Hour) // Newer
	td.InsertPrice(t, ctx, activePrice2, testPool)

	inactivePrice := td.NewValidPrice()
	inactivePrice.ProductID = product.ID
	inactivePrice.Amount = 1000
	inactivePrice.IsActive = false // Should not be returned
	td.InsertPrice(t, ctx, inactivePrice, testPool)

	// Make request
	url := fmt.Sprintf("%s/products/%s", testServerURL, product.ID.String())
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result domain.ProductAggregator
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	// Should only return active prices
	assert.Len(t, result.Prices, 2)

	// Verify all returned prices are active
	for _, price := range result.Prices {
		assert.True(t, price.IsActive)
		assert.Equal(t, product.ID, price.ProductID)
	}

	// Verify prices are ordered by creation date (newest first)
	assert.True(t, result.Prices[0].CreatedAt.After(result.Prices[1].CreatedAt))
}

// make test-func TEST_NAME=TestGetProductByID_WithNoPrices TEST_PATH=test/integration/catalog/product/get_product_by_id_test.go

func TestGetProductByID_WithNoPrices(t *testing.T) {
	ctx := context.Background()

	// Clean up tables
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)
	td.ClearImagesTable(t, ctx, testPool)
	td.ClearPricesTable(t, ctx, testPool)

	// Setup test data without any prices
	category := td.NewValidCategory("Free Stuff")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Free Product", category.ID)
	td.InsertProduct(t, ctx, testPool, product)

	// Make request
	url := fmt.Sprintf("%s/products/%s", testServerURL, product.ID.String())
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result domain.ProductAggregator
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	// Should have product data but empty prices array
	assert.Equal(t, product.ID, result.Product.ID)
	assert.Equal(t, product.Name, result.Product.Name)
	assert.Nil(t, result.Image)
	assert.Empty(t, result.Prices) // Empty slice, not nil
}

// make test-func TEST_NAME=TestGetProductByID_HTTPMethodNotAllowed TEST_PATH=test/integration/catalog/product/get_product_by_id_test.go

func TestGetProductByID_HTTPMethodNotAllowed(t *testing.T) {
	ctx := context.Background()

	// Clean up and setup minimal data
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)

	category := td.NewValidCategory("Test")
	td.InsertCategory(t, ctx, category, testPool)

	product := td.NewValidProduct("Test Product", category.ID)
	td.InsertProduct(t, ctx, testPool, product)

	// Test different HTTP methods
	testCases := []string{"POST", "PUT", "DELETE", "PATCH"}

	for _, method := range testCases {
		t.Run(fmt.Sprintf("Method_%s", method), func(t *testing.T) {
			url := fmt.Sprintf("%s/products/%s", testServerURL, product.ID.String())

			req, err := http.NewRequest(method, url, nil)
			require.NoError(t, err)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should return method not allowed (assuming your router is configured correctly)
			assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		})
	}
}

// Benchmark test for performance verification
func BenchmarkGetProductByID(b *testing.B) {
	ctx := context.Background()

	// Clean up tables
	// Clean up tables - use direct SQL since td helpers need *testing.T
	_, err := testPool.Exec(ctx, "TRUNCATE TABLE catalog.categories RESTART IDENTITY CASCADE;")
	if err != nil {
		b.Fatalf("Failed to clear categories table: %v", err)
	}
	_, err = testPool.Exec(ctx, "TRUNCATE catalog.products CASCADE;")
	if err != nil {
		b.Fatalf("Failed to clear products table: %v", err)
	}
	_, err = testPool.Exec(ctx, "TRUNCATE TABLE catalog.images RESTART IDENTITY CASCADE;")
	if err != nil {
		b.Fatalf("Failed to clear images table: %v", err)
	}
	_, err = testPool.Exec(ctx, "TRUNCATE TABLE catalog.prices RESTART IDENTITY CASCADE;")
	if err != nil {
		b.Fatalf("Failed to clear prices table: %v", err)
	}

	// Setup test data - create a wrapper testing.T for the helpers
	setupHelper := &testing.T{}
	// Setup test data
	category := td.NewValidCategory("Benchmark")
	td.InsertCategory(setupHelper, ctx, category, testPool)

	product := td.NewValidProduct("Benchmark Product", category.ID)
	td.InsertProduct(setupHelper, ctx, testPool, product)

	image := td.NewValidImage(product.ID)
	image.ParentType = domain.ProductType
	image.IsActive = true
	td.InsertImage(setupHelper, ctx, image, testPool)

	price := td.NewValidPrice()
	price.ProductID = product.ID
	td.InsertPrice(setupHelper, ctx, price, testPool)

	url := fmt.Sprintf("%s/products/%s", testServerURL, product.ID.String())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := http.Get(url)
			if err != nil {
				b.Fatalf("Request failed: %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				b.Fatalf("Expected 200, got %d", resp.StatusCode)
			}
			resp.Body.Close()
		}
	})
}
