package product_test

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

// make test-func TEST_NAME='^TestGetAdminAllProducts$$' TEST_PATH=test/integration/catalog/product/get_admin_all_products_test.go

func TestGetAdminAllProducts(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully return empty list when no products exist with valid admin token", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewGetAdminAllProductsRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var products []*domain.ProductAggregator
		err = json.NewDecoder(resp.Body).Decode(&products)
		require.NoError(t, err)
		assert.Empty(t, products, "Expected empty product list")
	})

	t.Run("should successfully return products with prices and images with valid admin token", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create test category
		category := td.NewValidCategory("Electronics")
		td.InsertCategory(t, ctx, category, testPool)

		// Create test products
		product1 := td.NewValidProduct("Laptop", category.ID)
		product2 := td.NewValidProduct("Mouse", category.ID)
		td.InsertProduct(t, ctx, testPool, product1)
		td.InsertProduct(t, ctx, testPool, product2)

		// Create prices for products
		price1 := td.NewValidPrice()
		price1.ProductID = product1.ID
		price1.Amount = 99999 // $999.99
		price1.Currency = "USD"
		price1.Interval = "one_time"

		price2 := td.NewValidPrice()
		price2.ProductID = product1.ID
		price2.Amount = 89999 // $899.99 (sale price)
		price2.Currency = "USD"
		price2.Interval = "one_time"
		price2.IsActive = false

		price3 := td.NewValidPrice()
		price3.ProductID = product2.ID
		price3.Amount = 2999 // $29.99
		price3.Currency = "USD"
		price3.Interval = "one_time"

		td.InsertPrice(t, ctx, price1, testPool)
		td.InsertPrice(t, ctx, price2, testPool)
		td.InsertPrice(t, ctx, price3, testPool)

		// Create images for products
		image1 := td.NewValidImage(product1.ID)
		image1.ParentType = domain.ProductType
		image1.Title = "Laptop Image"
		image1.IsActive = true

		image2 := td.NewValidImage(product2.ID)
		image2.ParentType = domain.ProductType
		image2.Title = "Mouse Image"
		image2.IsActive = true

		td.InsertImage(t, ctx, image1, testPool)
		td.InsertImage(t, ctx, image2, testPool)

		req := th.NewGetAdminAllProductsRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var products []*domain.ProductAggregator
		err = json.NewDecoder(resp.Body).Decode(&products)
		require.NoError(t, err)

		assert.Len(t, products, 2, "Expected 2 products")

		// Verify products are sorted by created_at DESC (newest first)
		assert.True(t, products[0].Product.CreatedAt.After(products[1].Product.CreatedAt) ||
			products[0].Product.CreatedAt.Equal(products[1].Product.CreatedAt),
			"Products should be sorted by created_at DESC")

		// Find laptop product
		var laptopProduct *domain.ProductAggregator
		for _, p := range products {
			if p.Product.Name == "Laptop" {
				laptopProduct = p
				break
			}
		}
		assert.NotNil(t, laptopProduct, "Laptop product should exist")

		// Verify laptop has 2 prices
		assert.Len(t, laptopProduct.Prices, 2, "Laptop should have 2 prices")

		// Verify laptop has image
		assert.NotNil(t, laptopProduct.Image, "Laptop should have an image")
		assert.Equal(t, "Laptop Image", laptopProduct.Image.Title)
		assert.True(t, laptopProduct.Image.IsActive)

		// Find mouse product
		var mouseProduct *domain.ProductAggregator
		for _, p := range products {
			if p.Product.Name == "Mouse" {
				mouseProduct = p
				break
			}
		}
		assert.NotNil(t, mouseProduct, "Mouse product should exist")

		// Verify mouse has 1 price
		assert.Len(t, mouseProduct.Prices, 1, "Mouse should have 1 price")
		assert.Equal(t, int(2999), mouseProduct.Prices[0].Amount)

		// Verify mouse has image
		assert.NotNil(t, mouseProduct.Image, "Mouse should have an image")
		assert.Equal(t, "Mouse Image", mouseProduct.Image.Title)
	})

	t.Run("should successfully return products without prices or images with valid admin token", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create test category
		category := td.NewValidCategory("Books")
		td.InsertCategory(t, ctx, category, testPool)

		// Create product without prices or images
		product := td.NewValidProduct("Programming Book", category.ID)
		td.InsertProduct(t, ctx, testPool, product)

		req := th.NewGetAdminAllProductsRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var products []*domain.ProductAggregator
		err = json.NewDecoder(resp.Body).Decode(&products)
		require.NoError(t, err)

		assert.Len(t, products, 1, "Expected 1 product")

		product1 := products[0]
		assert.Equal(t, "Programming Book", product1.Product.Name)
		assert.Empty(t, product1.Prices, "Product should have no prices")
		assert.Nil(t, product1.Image, "Product should have no image")
		assert.NotNil(t, product1.Product.Category, "Product should have category")
		assert.Equal(t, "Books", product1.Product.Category.Name)
	})

	t.Run("should successfully return products with only inactive images (no images in response) with valid admin token", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create test category
		category := td.NewValidCategory("Clothing")
		td.InsertCategory(t, ctx, category, testPool)

		// Create product
		product := td.NewValidProduct("T-Shirt", category.ID)
		td.InsertProduct(t, ctx, testPool, product)

		// Create inactive image
		image := td.NewValidImage(product.ID)
		image.ParentType = domain.ProductType
		image.IsActive = false // Inactive image
		td.InsertImage(t, ctx, image, testPool)

		req := th.NewGetAdminAllProductsRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var products []*domain.ProductAggregator
		err = json.NewDecoder(resp.Body).Decode(&products)
		require.NoError(t, err)

		assert.Len(t, products, 1, "Expected 1 product")

		product1 := products[0]
		assert.Equal(t, "T-Shirt", product1.Product.Name)
		assert.Nil(t, product1.Image, "Product should have no active image")
	})

	t.Run("should successfully handle multiple products with mixed data with valid admin token", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create multiple categories
		cat1 := td.NewValidCategory("Electronics")
		cat2 := td.NewValidCategory("Books")
		td.InsertCategory(t, ctx, cat1, testPool)
		td.InsertCategory(t, ctx, cat2, testPool)

		// Create products with different combinations of data
		// Product 1: Has prices and image
		product1 := td.NewValidProduct("Smartphone", cat1.ID)
		td.InsertProduct(t, ctx, testPool, product1)

		price1 := td.NewValidPrice()
		price1.ProductID = product1.ID
		td.InsertPrice(t, ctx, price1, testPool)

		image1 := td.NewValidImage(product1.ID)
		image1.ParentType = domain.ProductType
		image1.IsActive = true
		td.InsertImage(t, ctx, image1, testPool)

		// Product 2: Has prices but no image
		product2 := td.NewValidProduct("Tablet", cat1.ID)
		// Make this product created slightly before product1 for ordering test
		product2.CreatedAt = product1.CreatedAt.Add(-1 * time.Minute)
		td.InsertProduct(t, ctx, testPool, product2)

		price2 := td.NewValidPrice()
		price2.ProductID = product2.ID
		td.InsertPrice(t, ctx, price2, testPool)

		// Product 3: Has image but no prices
		product3 := td.NewValidProduct("Novel", cat2.ID)
		// Make this the oldest product
		product3.CreatedAt = product1.CreatedAt.Add(-2 * time.Minute)
		td.InsertProduct(t, ctx, testPool, product3)

		image3 := td.NewValidImage(product3.ID)
		image3.ParentType = domain.ProductType
		image3.IsActive = true
		td.InsertImage(t, ctx, image3, testPool)

		// Product 4: Has neither prices nor images
		product4 := td.NewValidProduct("Dictionary", cat2.ID)
		// Make this the second newest
		product4.CreatedAt = product1.CreatedAt.Add(-30 * time.Second)
		td.InsertProduct(t, ctx, testPool, product4)

		req := th.NewGetAdminAllProductsRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var products []*domain.ProductAggregator
		err = json.NewDecoder(resp.Body).Decode(&products)
		require.NoError(t, err)

		assert.Len(t, products, 4, "Expected 4 products")

		// Verify ordering (newest first)
		expectedOrder := []string{"Smartphone", "Dictionary", "Tablet", "Novel"}
		for i, expectedName := range expectedOrder {
			assert.Equal(t, expectedName, products[i].Product.Name,
				"Product at index %d should be %s", i, expectedName)
		}

		// Verify data combinations
		productMap := make(map[string]*domain.ProductAggregator)
		for _, p := range products {
			productMap[p.Product.Name] = p
		}

		// Smartphone: has price and image
		smartphone := productMap["Smartphone"]
		assert.Len(t, smartphone.Prices, 1, "Smartphone should have 1 price")
		assert.NotNil(t, smartphone.Image, "Smartphone should have image")

		// Tablet: has price but no image
		tablet := productMap["Tablet"]
		assert.Len(t, tablet.Prices, 1, "Tablet should have 1 price")
		assert.Nil(t, tablet.Image, "Tablet should have no image")

		// Novel: has image but no price
		novel := productMap["Novel"]
		assert.Empty(t, novel.Prices, "Novel should have no prices")
		assert.NotNil(t, novel.Image, "Novel should have image")

		// Dictionary: has neither
		dictionary := productMap["Dictionary"]
		assert.Empty(t, dictionary.Prices, "Dictionary should have no prices")
		assert.Nil(t, dictionary.Image, "Dictionary should have no image")
	})

	t.Run("should successfully return correct response format with valid admin token", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup minimal data
		category := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, category, testPool)

		product := td.NewValidProduct("Test Product", category.ID)
		td.InsertProduct(t, ctx, testPool, product)

		req := th.NewGetAdminAllProductsRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify content type
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response structure
		var products []*domain.ProductAggregator
		err = json.NewDecoder(resp.Body).Decode(&products)
		assert.NoError(t, err)

		assert.Len(t, products, 1)

		// Verify all expected fields are present
		assert.NotEmpty(t, products[0].Product.ID)
		assert.Equal(t, "Test Product", products[0].Product.Name)
		assert.NotNil(t, products[0].Product.Category)
		assert.Equal(t, "Test Category", products[0].Product.Category.Name)
		assert.Empty(t, products[0].Prices)
		// Image can be nil if no active images exist
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		clearTables(t, ctx)

		req := th.NewGetAdminAllProductsRequest(t, ctx, testServerURL, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		req := th.NewGetAdminAllProductsRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		req := th.NewGetAdminAllProductsRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		clearTables(t, ctx)

		req := th.NewGetAdminAllProductsRequest(t, ctx, testServerURL, "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
