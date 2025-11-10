package product_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	productHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/product"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME='^TestModifyProduct$' TEST_PATH=test/integration/catalog/product/modify_product_test.go

func TestModifyProduct(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully modify product name with valid admin token", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup test data
		category := th.NewValidCategory("Electronics")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Original Product Name", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		// Prepare update request
		newName := "Updated Product Name"
		updateReq := domain.UpdateProductRequest{
			Name: &newName,
		}

		req := th.NewModifyProductRequest(t, ctx, testServerURL, product.ID.String(), updateReq, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify database update
		updatedProduct, err := th.GetProductByID(t, ctx, product.ID, testPool)
		assert.NoError(t, err)
		assert.Equal(t, newName, updatedProduct.Name)
		assert.Equal(t, product.Description, updatedProduct.Description)
	})

	t.Run("should successfully modify product description with valid admin token", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup test data
		category := th.NewValidCategory("Services")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		// Prepare update request
		newDescription := "This is the updated description"
		updateReq := domain.UpdateProductRequest{
			Description: &newDescription,
		}

		req := th.NewModifyProductRequest(t, ctx, testServerURL, product.ID.String(), updateReq, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify database update
		updatedProduct, err := th.GetProductByID(t, ctx, product.ID, testPool)
		assert.NoError(t, err)
		assert.Equal(t, newDescription, updatedProduct.Description)
		assert.Equal(t, product.Name, updatedProduct.Name)
	})

	t.Run("should successfully modify product status with valid admin token", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup test data
		category := th.NewValidCategory("Books")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		product.Status = domain.Draft
		th.InsertProduct(t, ctx, testPool, product)

		// Prepare update request
		newStatus := domain.Published
		updateReq := domain.UpdateProductRequest{
			Status: th.StrPtr(string(newStatus)),
		}

		req := th.NewModifyProductRequest(t, ctx, testServerURL, product.ID.String(), updateReq, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify database update
		updatedProduct, err := th.GetProductByID(t, ctx, product.ID, testPool)
		assert.NoError(t, err)
		assert.Equal(t, newStatus, updatedProduct.Status)
	})

	t.Run("should successfully modify product metadata with valid admin token", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup test data
		category := th.NewValidCategory("Software")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		product.Metadata = map[string]any{"version": "1.0", "type": "basic"}
		th.InsertProduct(t, ctx, testPool, product)

		// Prepare update request
		newMetadata := map[string]any{
			"version":  "2.0",
			"type":     "premium",
			"features": []any{"advanced", "support"},
		}
		updateReq := domain.UpdateProductRequest{
			Metadata: newMetadata,
		}

		req := th.NewModifyProductRequest(t, ctx, testServerURL, product.ID.String(), updateReq, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify database update
		updatedProduct, err := th.GetProductByID(t, ctx, product.ID, testPool)
		assert.NoError(t, err)
		assert.Equal(t, newMetadata, updatedProduct.Metadata)
	})

	t.Run("should successfully modify multiple fields with valid admin token", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup test data
		category := th.NewValidCategory("Mixed")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Original Name", category.ID)
		product.Status = domain.Draft
		th.InsertProduct(t, ctx, testPool, product)

		// Prepare update request with multiple fields
		newName := "Updated Name"
		newDescription := "Updated description"
		newStatus := domain.Published
		newMetadata := map[string]any{"updated": true}

		updateReq := domain.UpdateProductRequest{
			Name:        &newName,
			Description: &newDescription,
			Status:      th.StrPtr(string(newStatus)),
			Metadata:    newMetadata,
		}

		req := th.NewModifyProductRequest(t, ctx, testServerURL, product.ID.String(), updateReq, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify all database updates
		updatedProduct, err := th.GetProductByID(t, ctx, product.ID, testPool)
		assert.NoError(t, err)
		assert.Equal(t, newName, updatedProduct.Name)
		assert.Equal(t, newDescription, updatedProduct.Description)
		assert.Equal(t, newStatus, updatedProduct.Status)
		assert.Equal(t, newMetadata, updatedProduct.Metadata)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		clearTables(t, ctx)

		category := th.NewValidCategory("Test")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		updateReq := domain.UpdateProductRequest{
			Name: th.StrPtr("Updated Name"),
		}

		req := th.NewModifyProductRequest(t, ctx, testServerURL, product.ID.String(), updateReq, "")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		category := th.NewValidCategory("Test")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		updateReq := domain.UpdateProductRequest{
			Name: th.StrPtr("Updated Name"),
		}

		req := th.NewModifyProductRequest(t, ctx, testServerURL, product.ID.String(), updateReq, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		category := th.NewValidCategory("Test")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		updateReq := domain.UpdateProductRequest{
			Name: th.StrPtr("Updated Name"),
		}

		req := th.NewModifyProductRequest(t, ctx, testServerURL, product.ID.String(), updateReq, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		clearTables(t, ctx)

		category := th.NewValidCategory("Test")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		updateReq := domain.UpdateProductRequest{
			Name: th.StrPtr("Updated Name"),
		}

		req := th.NewModifyProductRequest(t, ctx, testServerURL, product.ID.String(), updateReq, "invalid-token-12345")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 404 when product not found", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		nonExistentID := uuid.New()

		updateReq := domain.UpdateProductRequest{
			Name: th.StrPtr("Updated Name"),
		}

		req := th.NewModifyProductRequest(t, ctx, testServerURL, nonExistentID.String(), updateReq, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid request body", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		category := th.NewValidCategory("Test")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		// Invalid JSON request - need to manually construct with raw body
		req, err := http.NewRequest("PATCH", testServerURL+productHandler.AdminProductsBasePath+"/"+product.ID.String(), strings.NewReader(`{"name": "test"`))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid product ID format", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		updateReq := domain.UpdateProductRequest{
			Name: th.StrPtr("Updated Name"),
		}

		req := th.NewModifyProductRequest(t, ctx, testServerURL, "not-a-uuid", updateReq, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 405 for wrong HTTP method", func(t *testing.T) {
		clearTables(t, ctx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		category := th.NewValidCategory("Test")
		th.InsertCategory(t, ctx, category, testPool)

		product := th.NewValidProduct("Test Product", category.ID)
		th.InsertProduct(t, ctx, testPool, product)

		// Try using GET instead of PATCH - must manually construct to test wrong method
		req, err := http.NewRequest("GET", testServerURL+productHandler.AdminProductsBasePath+"/"+product.ID.String(), nil)
		require.NoError(t, err)

		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})
}
