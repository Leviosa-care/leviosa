package product_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateProductWithPrice TEST_PATH=test/integration/catalog/product/create_product_with_price_test.go
func TestCreateProductWithPrice(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Sample valid requests for the tests
	parentCategoryID := uuid.New().String()
	productNameRequest := "Integration Test Product"
	validProductRequest := domain.CreateProductRequest{
		CategoryID:        parentCategoryID,
		Name:              productNameRequest,
		Description:       "A product for integration testing purposes.",
		Duration:          30,
		Availability:      domain.Hybrid,
		BufferTime:        10,
		CancellationHours: 24,
	}
	validPriceRequest := domain.CreatePriceRequest{
		Amount:   1200,
		Currency: "EUR",
		Interval: "month",
	}

	// This is the combined request body we'll use for successful tests.
	validRequest := domain.CreateProductWithPriceRequest{
		Product: validProductRequest,
		Price:   validPriceRequest,
	}

	t.Run("should successfully create a product and a price and persist to DB", func(t *testing.T) {
		clearTables(t, ctx)

		// Pre-requisite: Create a category in the database as the product depends on it.
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)

		// Set the category ID in our request to the one we just created.
		validRequest.Product.CategoryID = parentCategory.ID.String()
		requestBody, _ := json.Marshal(validRequest)

		// --- Make the HTTP Request ---
		req := newCreateProductWithPriceRequest(t, ctx, requestBody)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		// --- Assert the HTTP Response ---
		assert.Equal(t, http.StatusCreated, res.StatusCode, "expected a 201 Created status code")

		var responseBody struct {
			ProductID string `json:"product_id"`
			PriceID   string `json:"price_id"`
			Message   string `json:"message"`
		}
		err = json.NewDecoder(res.Body).Decode(&responseBody)
		assert.NoError(t, err, "failed to decode response body")

		assert.NotEmpty(t, responseBody.ProductID, "response should contain a product ID")
		assert.NotEmpty(t, responseBody.PriceID, "response should contain a price ID")
		assert.Equal(t, "Product created successfully!", responseBody.Message)

		// --- Verify Persistence in the Database ---
		// Check if the product was created
		var productName string
		err = testPool.QueryRow(ctx, "SELECT name FROM catalog.products WHERE id = $1", responseBody.ProductID).Scan(&productName)
		assert.NoError(t, err, "failed to query product from database")
		assert.Equal(t, strings.ToLower(productNameRequest), productName, "product name in DB should match request")

		// Check if the price was created
		var priceAmount int64
		err = testPool.QueryRow(ctx, "SELECT amount FROM catalog.prices WHERE id = $1 AND product_id = $2", responseBody.PriceID, responseBody.ProductID).Scan(&priceAmount)
		assert.NoError(t, err, "failed to query price from database")
		assert.Equal(t, int64(1200), priceAmount, "price amount in DB should match request")
	})

	t.Run("should return 400 Bad Request for an invalid JSON body", func(t *testing.T) {
		clearTables(t, ctx)

		invalidBody := `{"product": {"product_name": "Test"}, "price": "invalid"}`
		req, err := http.NewRequestWithContext(ctx, "POST", testServerURL+"/admin/products", strings.NewReader(invalidBody))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		var response map[string]string
		json.NewDecoder(res.Body).Decode(&response)
		assert.Contains(t, response["error"], errs.ErrInvalidValue.Error())
	})

	t.Run("should return 400 Bad Request for invalid input data", func(t *testing.T) {
		clearTables(t, ctx)

		// Pre-requisite: Create a category
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)

		// Invalid request: price amount is negative
		invalidRequest := domain.CreateProductWithPriceRequest{
			Product: validProductRequest,
			Price: domain.CreatePriceRequest{
				Amount:   -100,
				Currency: "USD",
				Interval: "month",
			},
		}
		invalidRequest.Product.CategoryID = parentCategory.ID.String()
		requestBody, _ := json.Marshal(invalidRequest)

		req := newCreateProductWithPriceRequest(t, ctx, requestBody)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		var response map[string]string
		json.NewDecoder(res.Body).Decode(&response)
		assert.Contains(t, response["error"], errs.ErrInvalidValue.Error())
	})

	t.Run("should return 404 Not Found if the category ID does not exist", func(t *testing.T) {
		clearTables(t, ctx)

		// Request with a non-existent category ID
		requestWithBadCategory := domain.CreateProductWithPriceRequest{
			Product: domain.CreateProductRequest{
				CategoryID: uuid.New().String(), // Non-existent ID
				Name:       "Product with bad category",
			},
			Price: validPriceRequest,
		}
		requestBody, _ := json.Marshal(requestWithBadCategory)

		req := newCreateProductWithPriceRequest(t, ctx, requestBody)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		var response map[string]string
		json.NewDecoder(res.Body).Decode(&response)
		assert.Contains(t, response["error"], errs.ErrDomainNotFound.Error())
	})
}

// Helper function to create a new HTTP request for the CreateCategory handler.
func newCreateProductWithPriceRequest(t *testing.T, ctx context.Context, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/products", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}
