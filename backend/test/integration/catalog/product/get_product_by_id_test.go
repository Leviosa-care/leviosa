package product_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	productHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/product"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME='^TestGetProductByID$' TEST_PATH=test/integration/catalog/product/get_product_by_id_test.go

func TestGetProductByID(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get product with prices and image", func(t *testing.T) {
		// Clean up tables
		clearTables(t, ctx)

		// Setup test data
		category := th.NewValidCategory("Electronics")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		// Create an active image for the product
		image := th.NewValidImage(product.ID)
		image.ParentType = domain.ProductType
		image.IsActive = true
		th.InsertImage(t, ctx, image, testPool)

		// Create active prices for the product
		price1 := th.NewValidPrice()
		price1.ProductID = product.ID
		price1.Amount = 1000 // $10.00
		price1.Currency = "USD"
		price1.Interval = "month"
		price1.IsActive = true
		th.InsertPrice(t, ctx, price1, testPool)

		price2 := th.NewValidPrice()
		price2.ProductID = product.ID
		price2.Amount = 10000 // $100.00
		price2.Currency = "USD"
		price2.Interval = "year"
		price2.IsActive = true
		th.InsertPrice(t, ctx, price2, testPool)

		req := th.NewGetProductByIDRequest(t, ctx, testServerURL, product.ID.String())
		resp, err := client.Do(req)
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
	})

	t.Run("should successfully get product without image", func(t *testing.T) {
		// Clean up tables
		clearTables(t, ctx)

		// Setup test data - no image this time
		category := th.NewValidCategory("Services")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Service Product", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		// Create prices but no image
		price := th.NewValidPrice()
		price.ProductID = product.ID
		th.InsertPrice(t, ctx, price, testPool)

		req := th.NewGetProductByIDRequest(t, ctx, testServerURL, product.ID.String())
		resp, err := client.Do(req)
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
	})

	t.Run("should successfully get product with inactive image", func(t *testing.T) {
		// Clean up tables
		clearTables(t, ctx)

		// Setup test data
		category := th.NewValidCategory("Books")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Book Product", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		// Create an inactive image
		image := th.NewValidImage(product.ID)
		image.ParentType = domain.ProductType
		image.IsActive = false
		th.InsertImage(t, ctx, image, testPool)

		// Create price
		price := th.NewValidPrice()
		price.ProductID = product.ID
		th.InsertPrice(t, ctx, price, testPool)

		req := th.NewGetProductByIDRequest(t, ctx, testServerURL, product.ID.String())
		resp, err := client.Do(req)
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
	})

	t.Run("should return 404 for non-existent product", func(t *testing.T) {
		// Clean up tables
		clearTables(t, ctx)

		// Use a valid UUID that doesn't exist in the database
		nonExistentID := uuid.New()

		req := th.NewGetProductByIDRequest(t, ctx, testServerURL, nonExistentID.String())
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errorResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		assert.NoError(t, err)

		assert.Contains(t, errorResp["error"].(string), nonExistentID.String())
	})

	t.Run("should return 400 for invalid UUID format", func(t *testing.T) {
		req := th.NewGetProductByIDRequest(t, ctx, testServerURL, "not-a-uuid")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errorResp struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		assert.NoError(t, err)

		assert.Contains(t, errorResp.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should return 404 for empty product ID", func(t *testing.T) {
		req := th.NewGetProductByIDRequest(t, ctx, testServerURL, "")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 404 for product ID with slash", func(t *testing.T) {
		req := th.NewGetProductByIDRequest(t, ctx, testServerURL, "123/456")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should fail to create product with non-existent category ID", func(t *testing.T) {
		// Clean up tables
		clearTables(t, ctx)

		// Try to insert a product with a non-existent category ID
		nonExistentCategoryID := uuid.New()
		product := th.NewValidProduct("Invalid Product", nonExistentCategoryID)

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
	})

	t.Run("should get product with multiple prices ordered by creation date", func(t *testing.T) {
		// Clean up tables
		clearTables(t, ctx)

		// Setup test data
		category := th.NewValidCategory("Software")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Software Product", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		// Create multiple prices with different statuses
		activePrice1 := th.NewValidPrice()
		activePrice1.ProductID = product.ID
		activePrice1.Amount = 500
		activePrice1.Interval = "month"
		activePrice1.IsActive = true
		activePrice1.CreatedAt = time.Now().Add(-2 * time.Hour)
		th.InsertPrice(t, ctx, activePrice1, testPool)

		activePrice2 := th.NewValidPrice()
		activePrice2.ProductID = product.ID
		activePrice2.Amount = 5000
		activePrice2.Interval = "year"
		activePrice2.IsActive = true
		activePrice2.CreatedAt = time.Now().Add(-1 * time.Hour)
		th.InsertPrice(t, ctx, activePrice2, testPool)

		inactivePrice := th.NewValidPrice()
		inactivePrice.ProductID = product.ID
		inactivePrice.Amount = 1000
		inactivePrice.IsActive = false
		th.InsertPrice(t, ctx, inactivePrice, testPool)

		req := th.NewGetProductByIDRequest(t, ctx, testServerURL, product.ID.String())
		resp, err := client.Do(req)
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
	})

	t.Run("should get product with no prices", func(t *testing.T) {
		// Clean up tables
		clearTables(t, ctx)

		// Setup test data without any prices
		category := th.NewValidCategory("Free Stuff")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Free Product", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		req := th.NewGetProductByIDRequest(t, ctx, testServerURL, product.ID.String())
		resp, err := client.Do(req)
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
		assert.Empty(t, result.Prices)
	})

	t.Run("should return 405 for invalid HTTP methods", func(t *testing.T) {
		// Clean up and setup minimal data
		clearTables(t, ctx)

		category := th.NewValidCategory("Test")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		// Test different HTTP methods - must manually construct to test wrong methods
		testCases := []string{"POST", "PUT", "DELETE", "PATCH"}

		for _, method := range testCases {
			t.Run(method, func(t *testing.T) {
				url := testServerURL + productHandler.ProductsBasePath + "/" + product.ID.String()
				req, err := http.NewRequest(method, url, nil)
				require.NoError(t, err)

				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				// Should return method not allowed
				assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
			})
		}
	})
}

// Benchmark test for performance verification
func BenchmarkGetProductByID(b *testing.B) {
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
	category := th.NewValidCategory("Benchmark")
	th.InsertCategory(setupHelper, ctx, category, testPool)

	product := th.NewValidProduct("Benchmark Product", category.ID)
	th.InsertProduct(setupHelper, ctx, testPool, product)

	image := th.NewValidImage(product.ID)
	image.ParentType = domain.ProductType
	image.IsActive = true
	th.InsertImage(setupHelper, ctx, image, testPool)

	price := th.NewValidPrice()
	price.ProductID = product.ID
	th.InsertPrice(setupHelper, ctx, price, testPool)

	url := testServerURL + productHandler.ProductsBasePath + "/" + product.ID.String()

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
