package category_test

import (
	"bytes"
	"context"
	"encoding/json"
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

// make test-func TEST_NAME=TestModifyCategory TEST_PATH=test/integration/catalog/category/modify_category_test.go

func TestModifyCategory(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Setup: Clear the table and insert a known category for successful tests
	td.ClearCategoriesTable(t, ctx, testPool)
	existingCategory := td.NewValidCategory("Original Name")
	td.InsertCategory(t, ctx, existingCategory, testPool)

	t.Run("should successfully modify a category's name and description", func(t *testing.T) {
		updatedName := "Updated Category Name"
		updatedDescription := "Updated description."

		requestBody := domain.UpdateCategoryRequest{
			Name:        &updatedName,
			Description: &updatedDescription,
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newModifyCategoryRequest(t, ctx, existingCategory.ID.String(), jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Post-conditions: Verify the category was updated in the database
		updatedCategory, getErr := td.GetCategoryByID(t, ctx, existingCategory.ID, testPool)
		assert.NoError(t, getErr, "Failed to retrieve the updated category from the database")
		assert.Equal(t, updatedName, updatedCategory.Name)
		assert.Equal(t, updatedDescription, updatedCategory.Description)
	})

	t.Run("should return 404 Not Found if category ID does not exist", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		updatedName := "New Name"
		requestBody := domain.UpdateCategoryRequest{Name: &updatedName}
		jsonBody, _ := json.Marshal(requestBody)

		req := newModifyCategoryRequest(t, ctx, nonExistentID, jsonBody)
		req.Header.Set("Content-Type", "application/json")

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

	t.Run("should return 409 Conflict if category name already exists", func(t *testing.T) {
		// Clean and setup for this specific test
		td.ClearCategoriesTable(t, ctx, testPool)

		existingCat1 := td.NewValidCategory("Unique Category 1")
		td.InsertCategory(t, ctx, existingCat1, testPool)

		existingCat2 := td.NewValidCategory("Unique Category 2")
		td.InsertCategory(t, ctx, existingCat2, testPool)

		// Attempt to update cat1's name to cat2's name
		requestBody := domain.UpdateCategoryRequest{
			Name: &existingCat2.Name,
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newModifyCategoryRequest(t, ctx, existingCat1.ID.String(), jsonBody)
		req.Header.Set("Content-Type", "application/json")

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

	t.Run("should return 400 Bad Request for an invalid UUID in the URL", func(t *testing.T) {
		updatedName := "New Name"
		requestBody := domain.UpdateCategoryRequest{Name: &updatedName}
		jsonBody, _ := json.Marshal(requestBody)

		req := newModifyCategoryRequest(t, ctx, "not-a-valid-uuid", jsonBody)
		req.Header.Set("Content-Type", "application/json")

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

	t.Run("should return 400 Bad Request for an empty request body", func(t *testing.T) {
		req := newModifyCategoryRequest(t, ctx, existingCategory.ID.String(), []byte("{}"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, "no updatable fields provided")
	})

	t.Run("should return 400 Bad Request for invalid JSON body", func(t *testing.T) {
		req := newModifyCategoryRequest(t, ctx, existingCategory.ID.String(), []byte(`{"name": 123}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, "invalid request body")
	})

	t.Run("should return 400 Bad Request for unknown fields in request body", func(t *testing.T) {
		jsonBody := []byte(`{"unknown_field": "some_value"}`)

		req := newModifyCategoryRequest(t, ctx, existingCategory.ID.String(), jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, "invalid request body")
		assert.Contains(t, respBody.Error, "unknown field")
	})
}

// newModifyCategoryRequest creates a new HTTP request for the ModifyCategory handler.
func newModifyCategoryRequest(t *testing.T, ctx context.Context, categoryID string, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, testServerURL+"/admin/categories/"+categoryID, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}
