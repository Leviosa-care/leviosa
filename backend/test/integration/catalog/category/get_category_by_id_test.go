package category_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	categoryHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/category"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetCategoryByID TEST_PATH=test/integration/catalog/category/get_category_by_id_test.go

func TestGetCategoryByID(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Setup: Clear the table and insert known category and image for successful tests
	td.ClearCategoriesTable(t, ctx, testPool)
	existingCategory := td.NewValidCategory("Existing Test Category")
	td.InsertCategory(t, ctx, existingCategory, testPool)

	activeImage := td.NewValidImage(existingCategory.ID)
	activeImage.Title = "Test Image Title"
	activeImage.IsActive = true
	td.InsertImage(t, ctx, activeImage, testPool)

	type GetCategoryByIDResponse struct {
		Image    *domain.Image    `json:"image"`
		Category *domain.Category `json:"category"`
	}

	t.Run("should successfully get a category by ID", func(t *testing.T) {
		req := newGetCategoryByIDRequest(t, ctx, existingCategory.ID.String())
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody GetCategoryByIDResponse
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)

		// Post-conditions: Verify the returned category data matches the inserted one
		assert.Equal(t, existingCategory.ID, responseBody.Category.ID)
		assert.Equal(t, existingCategory.Name, responseBody.Category.Name)
		assert.Equal(t, existingCategory.Description, responseBody.Category.Description)
		assert.NotNil(t, responseBody.Image, "Expected a non-nil image object in the response")
		assert.Equal(t, activeImage.ID, responseBody.Image.ID)
		assert.True(t, responseBody.Image.IsActive, "Expected the image to be active")
	})

	t.Run("should successfully get a category without an active image", func(t *testing.T) {
		// Setup: Clear the table and insert a known category
		td.ClearCategoriesTable(t, ctx, testPool)
		existingCategory := td.NewValidCategory("Category Without Image")
		td.InsertCategory(t, ctx, existingCategory, testPool)

		// No image is inserted for this category

		req := newGetCategoryByIDRequest(t, ctx, existingCategory.ID.String())
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody GetCategoryByIDResponse
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)

		// Post-conditions: Verify the returned category and that the image is nil
		assert.Equal(t, existingCategory.ID, responseBody.Category.ID)
		assert.Equal(t, existingCategory.Name, responseBody.Category.Name)
		assert.Nil(t, responseBody.Image, "Expected a nil image object in the response")
	})

	t.Run("should return 404 Not Found if category ID does not exist", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		req := newGetCategoryByIDRequest(t, ctx, nonExistentID)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, "not found")
	})

	t.Run("should return 400 Bad Request for an invalid UUID format", func(t *testing.T) {
		req := newGetCategoryByIDRequest(t, ctx, "not-a-valid-uuid")
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})
}

// newGetCategoryByIDRequest creates a new HTTP request for the GetCategoryByID handler.
func newGetCategoryByIDRequest(t *testing.T, ctx context.Context, categoryID string) *http.Request {
	url := testServerURL + categoryHandler.CategoriesBasePath + "/" + categoryID
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	require.NoError(t, err)
	return req
}
