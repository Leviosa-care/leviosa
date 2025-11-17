package category_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestRemoveCategory TEST_PATH=test/integration/catalog/category/remove_category_test.go

func TestRemoveCategory(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully remove a category and its images with valid admin token", func(t *testing.T) {
		// Clean the database and S3 bucket to ensure a fresh start
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup: Create a category
		cat := td.NewValidCategory("Category to Remove")
		td.InsertCategory(t, ctx, cat, testPool)

		// Setup: Create an image associated with the category
		img := td.NewValidImage(cat.ID)
		img.Title = "Test Image"

		td.InsertImage(t, ctx, img, testPool)
		// Assume the image file is also in S3 (e.g., from a previous UploadImage call)
		// For this test, we can just verify the database entry is gone.
		// A full test would also verify the S3 file is gone, but that requires
		// mocking the S3 service or pre-populating the bucket which is complex.
		// We'll rely on the image service's internal tests to cover S3 deletion logic.

		req := th.NewRemoveCategoryRequest(t, ctx, testServerURL, cat.ID.String(), accessToken)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Post-conditions: Verify the category is no longer in the database
		var categoryCount int
		err = testPool.QueryRow(ctx, "SELECT count(*) FROM catalog.categories WHERE id = $1", cat.ID).Scan(&categoryCount)
		assert.NoError(t, err)
		assert.Equal(t, 0, categoryCount, "Expected category to be removed from the database")

		// Post-conditions: Verify the image is no longer in the database
		var imageCount int
		err = testPool.QueryRow(ctx, "SELECT count(*) FROM catalog.images WHERE parent_id = $1", cat.ID).Scan(&imageCount)
		assert.NoError(t, err)
		assert.Equal(t, 0, imageCount, "Expected images to be removed from the database")
	})

	t.Run("should return 404 when the category does not exist", func(t *testing.T) {
		// Clean the database to ensure isolation
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		nonExistentID := uuid.New().String()

		req := th.NewRemoveCategoryRequest(t, ctx, testServerURL, nonExistentID, accessToken)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrRepositoryNotFound.Error())
	})

	t.Run("should return 409 when the category still has associated products", func(t *testing.T) {
		// Clean the database to ensure isolation
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup: Create a category
		cat := td.NewValidCategory("Category with Products")
		td.InsertCategory(t, ctx, cat, testPool)

		// Setup: Create a product associated with the category
		prod := td.NewValidProduct("Product in Category", cat.ID)
		td.InsertProduct(t, ctx, testPool, prod)

		req := th.NewRemoveCategoryRequest(t, ctx, testServerURL, cat.ID.String(), accessToken)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrConflict.Error())
	})

	t.Run("should return 400 for an invalid category ID", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewRemoveCategoryRequest(t, ctx, testServerURL, "not-a-uuid", accessToken)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, "parent ID must be a valid UUID")
	})

	// Authentication test cases
	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)

		existingCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, existingCategory, testPool)

		req := th.NewRemoveCategoryRequest(t, ctx, testServerURL, existingCategory.ID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		existingCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, existingCategory, testPool)

		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		req := th.NewRemoveCategoryRequest(t, ctx, testServerURL, existingCategory.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		existingCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, existingCategory, testPool)

		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		req := th.NewRemoveCategoryRequest(t, ctx, testServerURL, existingCategory.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)

		existingCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, existingCategory, testPool)

		req := th.NewRemoveCategoryRequest(t, ctx, testServerURL, existingCategory.ID.String(), "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
