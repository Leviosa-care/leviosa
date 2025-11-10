package product_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	productHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/product"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

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
	_ = validRequest

	t.Run("should successfully create product with price with valid admin token", func(t *testing.T) {
		clearTables(t, ctx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Pre-requisite: Create a category in the database as the product depends on it.
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)

		// Set the category ID in our request to the one we just created.
		validRequest.Product.CategoryID = parentCategory.ID.String()

		// --- Make the HTTP Request ---
		req := th.NewCreateProductWithPriceRequest(t, ctx, testServerURL, validRequest, accessToken)

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

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		invalidBody := `{"product": {"product_name": "Test"}, "price": "invalid"}`
		req, err := http.NewRequestWithContext(ctx, "POST", testServerURL+productHandler.CreateProductWithPriceEndpoint, strings.NewReader(invalidBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)

		client := &http.Client{}
		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		var response map[string]string
		json.NewDecoder(res.Body).Decode(&response)
		assert.Contains(t, response["error"], errs.ErrInvalidValue.Error())
	})

	t.Run("should return 400 Bad Request for invalid input data", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

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

		req := th.NewCreateProductWithPriceRequest(t, ctx, testServerURL, invalidRequest, accessToken)

		client := &http.Client{}
		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		var response map[string]string
		json.NewDecoder(res.Body).Decode(&response)
		assert.Contains(t, response["error"], errs.ErrInvalidValue.Error())
	})

	t.Run("should return 404 Not Found if the category ID does not exist", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Request with a non-existent category ID
		requestWithBadCategory := domain.CreateProductWithPriceRequest{
			Product: domain.CreateProductRequest{
				CategoryID: uuid.New().String(), // Non-existent ID
				Name:       "Product with bad category",
			},
			Price: validPriceRequest,
		}

		req := th.NewCreateProductWithPriceRequest(t, ctx, testServerURL, requestWithBadCategory, accessToken)

		client := &http.Client{}
		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		var response map[string]string
		json.NewDecoder(res.Body).Decode(&response)
		assert.Contains(t, response["error"], errs.ErrDomainNotFound.Error())
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		clearTables(t, ctx)

		// Pre-requisite: Create a category in the database as the product depends on it.
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)

		// Set the category ID in our request to the one we just created.
		validRequest.Product.CategoryID = parentCategory.ID.String()

		req := th.NewCreateProductWithPriceRequest(t, ctx, testServerURL, validRequest, "")

		client := &http.Client{}
		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		clearTables(t, ctx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create expired admin session
		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Pre-requisite: Create a category in the database as the product depends on it.
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)

		// Set the category ID in our request to the one we just created.
		validRequest.Product.CategoryID = parentCategory.ID.String()

		req := th.NewCreateProductWithPriceRequest(t, ctx, testServerURL, validRequest, accessToken)

		client := &http.Client{}
		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		clearTables(t, ctx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create standard user (not admin)
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Pre-requisite: Create a category in the database as the product depends on it.
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)

		// Set the category ID in our request to the one we just created.
		validRequest.Product.CategoryID = parentCategory.ID.String()

		req := th.NewCreateProductWithPriceRequest(t, ctx, testServerURL, validRequest, accessToken)

		client := &http.Client{}
		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusForbidden, res.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		clearTables(t, ctx)

		// Pre-requisite: Create a category in the database as the product depends on it.
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)

		// Set the category ID in our request to the one we just created.
		validRequest.Product.CategoryID = parentCategory.ID.String()

		req := th.NewCreateProductWithPriceRequest(t, ctx, testServerURL, validRequest, "invalid-token-12345")

		client := &http.Client{}
		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})
}
