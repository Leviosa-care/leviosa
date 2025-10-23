package product_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllPublishedProducts_Integration(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		setupData      func(t *testing.T, ctx context.Context)
		expectedStatus int
		validateFunc   func(t *testing.T, products []*domain.ProductAggregator)
	}{
		{
			name: "should return empty list when no published products exist",
			setupData: func(t *testing.T, ctx context.Context) {
				// Clean all tables
				td.ClearProductsTable(t, ctx, testPool)
				td.ClearPricesTable(t, ctx, testPool)
				td.ClearImagesTable(t, ctx, testPool)
				td.ClearCategoriesTable(t, ctx, testPool)
			},
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, products []*domain.ProductAggregator) {
				assert.Empty(t, products, "Expected empty product list when no products exist")
			},
		},
		{
			name: "should return empty list when only draft products exist",
			setupData: func(t *testing.T, ctx context.Context) {
				// Clean tables first
				td.ClearProductsTable(t, ctx, testPool)
				td.ClearPricesTable(t, ctx, testPool)
				td.ClearImagesTable(t, ctx, testPool)
				td.ClearCategoriesTable(t, ctx, testPool)

				// Create test category
				category := td.NewValidCategory("Electronics")
				td.InsertCategory(t, ctx, category, testPool)

				// Create draft products only
				draftProduct1 := td.NewValidProduct("Draft Laptop", category.ID)
				draftProduct1.Status = domain.Draft // Explicitly set to draft
				draftProduct2 := td.NewValidProduct("Draft Mouse", category.ID)
				draftProduct2.Status = domain.Draft

				td.InsertProduct(t, ctx, testPool, draftProduct1)
				td.InsertProduct(t, ctx, testPool, draftProduct2)
			},
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, products []*domain.ProductAggregator) {
				assert.Empty(t, products, "Expected empty product list when only draft products exist")
			},
		},
		{
			name: "should return only published products with prices and images",
			setupData: func(t *testing.T, ctx context.Context) {
				// Clean tables first
				td.ClearProductsTable(t, ctx, testPool)
				td.ClearPricesTable(t, ctx, testPool)
				td.ClearImagesTable(t, ctx, testPool)
				td.ClearCategoriesTable(t, ctx, testPool)

				// Create test category
				category := td.NewValidCategory("Electronics")
				td.InsertCategory(t, ctx, category, testPool)

				// Create published products
				publishedProduct1 := td.NewValidProduct("Published Laptop", category.ID)
				publishedProduct1.Status = domain.Published
				publishedProduct2 := td.NewValidProduct("Published Mouse", category.ID)
				publishedProduct2.Status = domain.Published

				// Create draft product (should not appear in results)
				draftProduct := td.NewValidProduct("Draft Product", category.ID)
				draftProduct.Status = domain.Draft

				td.InsertProduct(t, ctx, testPool, publishedProduct1)
				td.InsertProduct(t, ctx, testPool, publishedProduct2)
				td.InsertProduct(t, ctx, testPool, draftProduct)

				// Create prices for published products
				price1 := td.NewValidPrice()
				price1.ProductID = publishedProduct1.ID
				price1.Amount = 99999 // $999.99
				price1.Currency = "USD"
				price1.Interval = "one_time"

				price2 := td.NewValidPrice()
				price2.ProductID = publishedProduct1.ID
				price2.Amount = 89999 // $899.99 (sale price)
				price2.Currency = "USD"
				price2.Interval = "one_time"
				price2.IsActive = false

				price3 := td.NewValidPrice()
				price3.ProductID = publishedProduct2.ID
				price3.Amount = 2999 // $29.99
				price3.Currency = "USD"
				price3.Interval = "one_time"

				// Also create a price for draft product (should not appear)
				priceDraft := td.NewValidPrice()
				priceDraft.ProductID = draftProduct.ID
				priceDraft.Amount = 5000

				td.InsertPrice(t, ctx, price1, testPool)
				td.InsertPrice(t, ctx, price2, testPool)
				td.InsertPrice(t, ctx, price3, testPool)
				td.InsertPrice(t, ctx, priceDraft, testPool)

				// Create images for published products
				image1 := td.NewValidImage(publishedProduct1.ID)
				image1.ParentType = domain.ProductType
				image1.Title = "Published Laptop Image"
				image1.IsActive = true

				image2 := td.NewValidImage(publishedProduct2.ID)
				image2.ParentType = domain.ProductType
				image2.Title = "Published Mouse Image"
				image2.IsActive = true

				// Create image for draft product (should not appear)
				imageDraft := td.NewValidImage(draftProduct.ID)
				imageDraft.ParentType = domain.ProductType
				imageDraft.Title = "Draft Product Image"
				imageDraft.IsActive = true

				td.InsertImage(t, ctx, image1, testPool)
				td.InsertImage(t, ctx, image2, testPool)
				td.InsertImage(t, ctx, imageDraft, testPool)
			},
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, products []*domain.ProductAggregator) {
				require.Len(t, products, 2, "Expected 2 published products only")

				// Verify all returned products are published
				for _, p := range products {
					assert.Equal(t, domain.Published, p.Product.Status,
						"All returned products should be published")
				}

				// Verify no draft products are returned
				for _, p := range products {
					assert.NotEqual(t, "Draft Product", p.Product.Name,
						"Draft products should not be returned")
				}

				// Verify products are sorted by created_at DESC (newest first)
				assert.True(t, products[0].Product.CreatedAt.After(products[1].Product.CreatedAt) ||
					products[0].Product.CreatedAt.Equal(products[1].Product.CreatedAt),
					"Products should be sorted by created_at DESC")

				// Find laptop product
				var laptopProduct *domain.ProductAggregator
				for _, p := range products {
					if p.Product.Name == "Published Laptop" {
						laptopProduct = p
						break
					}
				}
				require.NotNil(t, laptopProduct, "Published Laptop should exist")

				// Verify laptop has 2 prices
				assert.Len(t, laptopProduct.Prices, 2, "Laptop should have 2 prices")

				// Verify laptop has image
				require.NotNil(t, laptopProduct.Image, "Laptop should have an image")
				assert.Equal(t, "Published Laptop Image", laptopProduct.Image.Title)
				assert.True(t, laptopProduct.Image.IsActive)

				// Find mouse product
				var mouseProduct *domain.ProductAggregator
				for _, p := range products {
					if p.Product.Name == "Published Mouse" {
						mouseProduct = p
						break
					}
				}
				require.NotNil(t, mouseProduct, "Published Mouse should exist")

				// Verify mouse has 1 price
				assert.Len(t, mouseProduct.Prices, 1, "Mouse should have 1 price")
				assert.Equal(t, int(2999), mouseProduct.Prices[0].Amount)

				// Verify mouse has image
				require.NotNil(t, mouseProduct.Image, "Mouse should have an image")
				assert.Equal(t, "Published Mouse Image", mouseProduct.Image.Title)
			},
		},
		{
			name: "should return published products without prices or images",
			setupData: func(t *testing.T, ctx context.Context) {
				// Clean tables
				td.ClearProductsTable(t, ctx, testPool)
				td.ClearPricesTable(t, ctx, testPool)
				td.ClearImagesTable(t, ctx, testPool)
				td.ClearCategoriesTable(t, ctx, testPool)

				// Create test category
				category := td.NewValidCategory("Books")
				td.InsertCategory(t, ctx, category, testPool)

				// Create published product without prices or images
				publishedProduct := td.NewValidProduct("Published Programming Book", category.ID)
				publishedProduct.Status = domain.Published

				// Create draft product (should not appear)
				draftProduct := td.NewValidProduct("Draft Book", category.ID)
				draftProduct.Status = domain.Draft

				td.InsertProduct(t, ctx, testPool, publishedProduct)
				td.InsertProduct(t, ctx, testPool, draftProduct)
			},
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, products []*domain.ProductAggregator) {
				require.Len(t, products, 1, "Expected 1 published product")

				product := products[0]
				assert.Equal(t, "Published Programming Book", product.Product.Name)
				assert.Equal(t, domain.Published, product.Product.Status)
				assert.Empty(t, product.Prices, "Product should have no prices")
				assert.Nil(t, product.Image, "Product should have no image")
				assert.NotNil(t, product.Product.Category, "Product should have category")
				assert.Equal(t, "Books", product.Product.Category.Name)
			},
		},
		{
			name: "should return published products with only inactive images (no images in response)",
			setupData: func(t *testing.T, ctx context.Context) {
				// Clean tables
				td.ClearProductsTable(t, ctx, testPool)
				td.ClearPricesTable(t, ctx, testPool)
				td.ClearImagesTable(t, ctx, testPool)
				td.ClearCategoriesTable(t, ctx, testPool)

				// Create test category
				category := td.NewValidCategory("Clothing")
				td.InsertCategory(t, ctx, category, testPool)

				// Create published product
				publishedProduct := td.NewValidProduct("Published T-Shirt", category.ID)
				publishedProduct.Status = domain.Published
				td.InsertProduct(t, ctx, testPool, publishedProduct)

				// Create inactive image for published product
				image := td.NewValidImage(publishedProduct.ID)
				image.ParentType = domain.ProductType
				image.IsActive = false // Inactive image
				td.InsertImage(t, ctx, image, testPool)
			},
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, products []*domain.ProductAggregator) {
				require.Len(t, products, 1, "Expected 1 published product")

				product := products[0]
				assert.Equal(t, "Published T-Shirt", product.Product.Name)
				assert.Equal(t, domain.Published, product.Product.Status)
				assert.Nil(t, product.Image, "Product should have no active image")
			},
		},
		{
			name: "should handle mixed status products with various data combinations",
			setupData: func(t *testing.T, ctx context.Context) {
				// Clean tables
				td.ClearProductsTable(t, ctx, testPool)
				td.ClearPricesTable(t, ctx, testPool)
				td.ClearImagesTable(t, ctx, testPool)
				td.ClearCategoriesTable(t, ctx, testPool)

				// Create multiple categories
				cat1 := td.NewValidCategory("Electronics")
				cat2 := td.NewValidCategory("Books")
				td.InsertCategory(t, ctx, cat1, testPool)
				td.InsertCategory(t, ctx, cat2, testPool)

				baseTime := time.Now()

				// Published Product 1: Has prices and image (newest)
				publishedProduct1 := td.NewValidProduct("Published Smartphone", cat1.ID)
				publishedProduct1.Status = domain.Published
				publishedProduct1.CreatedAt = baseTime
				td.InsertProduct(t, ctx, testPool, publishedProduct1)

				price1 := td.NewValidPrice()
				price1.ProductID = publishedProduct1.ID
				td.InsertPrice(t, ctx, price1, testPool)

				image1 := td.NewValidImage(publishedProduct1.ID)
				image1.ParentType = domain.ProductType
				image1.IsActive = true
				td.InsertImage(t, ctx, image1, testPool)

				// Draft Product: Has prices and image but should NOT appear
				draftProduct := td.NewValidProduct("Draft Tablet", cat1.ID)
				draftProduct.Status = domain.Draft
				draftProduct.CreatedAt = baseTime.Add(-10 * time.Second)
				td.InsertProduct(t, ctx, testPool, draftProduct)

				draftPrice := td.NewValidPrice()
				draftPrice.ProductID = draftProduct.ID
				td.InsertPrice(t, ctx, draftPrice, testPool)

				draftImage := td.NewValidImage(draftProduct.ID)
				draftImage.ParentType = domain.ProductType
				draftImage.IsActive = true
				td.InsertImage(t, ctx, draftImage, testPool)

				// Published Product 2: Has prices but no image
				publishedProduct2 := td.NewValidProduct("Published Laptop", cat1.ID)
				publishedProduct2.Status = domain.Published
				publishedProduct2.CreatedAt = baseTime.Add(-30 * time.Second)
				td.InsertProduct(t, ctx, testPool, publishedProduct2)

				price2 := td.NewValidPrice()
				price2.ProductID = publishedProduct2.ID
				td.InsertPrice(t, ctx, price2, testPool)

				// Published Product 3: Has image but no prices (oldest)
				publishedProduct3 := td.NewValidProduct("Published Novel", cat2.ID)
				publishedProduct3.Status = domain.Published
				publishedProduct3.CreatedAt = baseTime.Add(-2 * time.Minute)
				td.InsertProduct(t, ctx, testPool, publishedProduct3)

				image3 := td.NewValidImage(publishedProduct3.ID)
				image3.ParentType = domain.ProductType
				image3.IsActive = true
				td.InsertImage(t, ctx, image3, testPool)

				// Published Product 4: Has neither prices nor images
				publishedProduct4 := td.NewValidProduct("Published Dictionary", cat2.ID)
				publishedProduct4.Status = domain.Published
				publishedProduct4.CreatedAt = baseTime.Add(-1 * time.Minute)
				td.InsertProduct(t, ctx, testPool, publishedProduct4)
			},
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, products []*domain.ProductAggregator) {
				require.Len(t, products, 4, "Expected 4 published products only")

				// Verify all returned products are published
				for _, p := range products {
					assert.Equal(t, domain.Published, p.Product.Status,
						"All returned products should be published, got %s for %s",
						p.Product.Status, p.Product.Name)
				}

				// Verify no draft products are returned
				for _, p := range products {
					assert.NotEqual(t, "Draft Tablet", p.Product.Name,
						"Draft products should not be returned")
				}

				// Verify ordering (newest first)
				expectedOrder := []string{"Published Smartphone", "Published Laptop", "Published Dictionary", "Published Novel"}
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
				smartphone := productMap["Published Smartphone"]
				assert.Len(t, smartphone.Prices, 1, "Smartphone should have 1 price")
				assert.NotNil(t, smartphone.Image, "Smartphone should have image")

				// Laptop: has price but no image
				laptop := productMap["Published Laptop"]
				assert.Len(t, laptop.Prices, 1, "Laptop should have 1 price")
				assert.Nil(t, laptop.Image, "Laptop should have no image")

				// Novel: has image but no price
				novel := productMap["Published Novel"]
				assert.Empty(t, novel.Prices, "Novel should have no prices")
				assert.NotNil(t, novel.Image, "Novel should have image")

				// Dictionary: has neither
				dictionary := productMap["Published Dictionary"]
				assert.Empty(t, dictionary.Prices, "Dictionary should have no prices")
				assert.Nil(t, dictionary.Image, "Dictionary should have no image")
			},
		},
		{
			name: "should handle all possible product statuses correctly",
			setupData: func(t *testing.T, ctx context.Context) {
				// Clean tables
				td.ClearProductsTable(t, ctx, testPool)
				td.ClearPricesTable(t, ctx, testPool)
				td.ClearImagesTable(t, ctx, testPool)
				td.ClearCategoriesTable(t, ctx, testPool)

				// Create test category
				category := td.NewValidCategory("Mixed Status")
				td.InsertCategory(t, ctx, category, testPool)

				// Create products with different statuses
				publishedProduct := td.NewValidProduct("Published Product", category.ID)
				publishedProduct.Status = domain.Published

				draftProduct := td.NewValidProduct("Draft Product", category.ID)
				draftProduct.Status = domain.Draft

				// If you have other statuses like Archived, Inactive, etc., add them here
				// archivedProduct := td.NewValidProduct("Archived Product", category.ID)
				// archivedProduct.Status = domain.Archived

				td.InsertProduct(t, ctx, testPool, publishedProduct)
				td.InsertProduct(t, ctx, testPool, draftProduct)
			},
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, products []*domain.ProductAggregator) {
				require.Len(t, products, 1, "Expected only 1 published product")

				product := products[0]
				assert.Equal(t, "Published Product", product.Product.Name)
				assert.Equal(t, domain.Published, product.Product.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test data
			tt.setupData(t, ctx)

			// Make HTTP request to the published products endpoint
			url := fmt.Sprintf("%s/products", testServerURL) // Adjust endpoint as needed
			resp, err := http.Get(url)
			require.NoError(t, err, "Failed to make GET request")
			defer resp.Body.Close()

			// Verify status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode, "Unexpected status code")

			if tt.expectedStatus == http.StatusOK {
				// Parse response
				var products []*domain.ProductAggregator
				err = json.NewDecoder(resp.Body).Decode(&products)
				require.NoError(t, err, "Failed to decode response")

				// Run validation
				tt.validateFunc(t, products)
			}
		})
	}
}

func TestGetAllPublishedProducts_ResponseFormat(t *testing.T) {
	ctx := context.Background()

	// Setup minimal data
	td.ClearProductsTable(t, ctx, testPool)
	td.ClearPricesTable(t, ctx, testPool)
	td.ClearImagesTable(t, ctx, testPool)
	td.ClearCategoriesTable(t, ctx, testPool)

	category := td.NewValidCategory("Test Category")
	td.InsertCategory(t, ctx, category, testPool)

	// Create a published product
	product := td.NewValidProduct("Test Published Product", category.ID)
	product.Status = domain.Published // Ensure it's published
	td.InsertProduct(t, ctx, testPool, product)

	// Make request
	url := fmt.Sprintf("%s/products", testServerURL) // Adjust endpoint as needed
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verify content type
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response structure
	var products []*domain.ProductAggregator
	err = json.NewDecoder(resp.Body).Decode(&products)
	require.NoError(t, err)

	require.Len(t, products, 1)

	// Verify all expected fields are present
	assert.NotEmpty(t, products[0].Product.ID)
	assert.Equal(t, "Test Published Product", products[0].Product.Name)
	assert.Equal(t, domain.Published, products[0].Product.Status)
	assert.NotNil(t, products[0].Product.Category)
	assert.Equal(t, "Test Category", products[0].Product.Category.Name)
	assert.Empty(t, products[0].Prices) // Should be empty slice, not nil
	// Image can be nil if no active images exist
}

func TestGetAllPublishedProducts_ErrorScenarios(t *testing.T) {
	ctx := context.Background()

	t.Run("should handle invalid HTTP methods", func(t *testing.T) {
		url := fmt.Sprintf("%s/products", testServerURL)

		// Test POST method (should be method not allowed if your router is strict)
		resp, err := http.Post(url, "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Depending on your router setup, this might be 404 or 405
		assert.True(t, resp.StatusCode == http.StatusNotFound ||
			resp.StatusCode == http.StatusMethodNotAllowed,
			"Expected 404 or 405 for invalid method, got %d", resp.StatusCode)
	})

	t.Run("should return empty list when products exist but none are published", func(t *testing.T) {
		// Clean tables
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)
		td.ClearImagesTable(t, ctx, testPool)
		td.ClearCategoriesTable(t, ctx, testPool)

		// Create test category
		category := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, category, testPool)

		// Create multiple products, but all are draft
		for i := range 5 {
			product := td.NewValidProduct(fmt.Sprintf("Draft Product %d", i+1), category.ID)
			product.Status = domain.Draft
			td.InsertProduct(t, ctx, testPool, product)
		}

		// Make request
		url := fmt.Sprintf("%s/products", testServerURL)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return empty list
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var products []*domain.ProductAggregator
		err = json.NewDecoder(resp.Body).Decode(&products)
		require.NoError(t, err)

		assert.Empty(t, products, "Should return empty list when no published products exist")
	})
}
