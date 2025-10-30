package category_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// newRemoveCategoryRequest creates a new HTTP request for the RemoveCategory handler.
func newRemoveCategoryRequest(t *testing.T, ctx context.Context, categoryID string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, testServerURL+"/admin/categories/"+categoryID, nil)
	assert.NoError(t, err)
	return req
}

func TestRemoveCategory(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully remove a category and its images", func(t *testing.T) {
		// Clean the database and S3 bucket to ensure a fresh start
		clearTables(t, ctx, testPool)

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

		req := newRemoveCategoryRequest(t, ctx, cat.ID.String())
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
		clearTables(t, ctx, testPool)

		nonExistentID := uuid.New().String()

		req := newRemoveCategoryRequest(t, ctx, nonExistentID)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrDomainNotFound.Error())
	})

	t.Run("should return 409 when the category still has associated products", func(t *testing.T) {
		// Clean the database to ensure isolation
		clearTables(t, ctx, testPool)

		// Setup: Create a category
		cat := td.NewValidCategory("Category with Products")
		td.InsertCategory(t, ctx, cat, testPool)

		// Setup: Create a product associated with the category
		prod := td.NewValidProduct("Product in Category", cat.ID)
		td.InsertProduct(t, ctx, testPool, prod)

		req := newRemoveCategoryRequest(t, ctx, cat.ID.String())
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Equal(t, "cannot delete category: it still has associated products", respBody.Error)
	})

	t.Run("should return 400 for an invalid category ID", func(t *testing.T) {
		req := newRemoveCategoryRequest(t, ctx, "not-a-uuid")
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
}
