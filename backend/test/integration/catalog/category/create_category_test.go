package category_test

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

// make test-func TEST_NAME=TestCreateCategory TEST_PATH=test/integration/catalog/category/create_category_test.go

func TestCreateCategory(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully create a new category", func(t *testing.T) {
		// Clean the database to ensure isolation for this test
		td.ClearCategoriesTable(t, ctx, testPool)

		requestBody := domain.CreateCategoryRequest{
			Name:        "New Category",
			Description: "A great new category.",
			Metadata:    map[string]any{"key1": "value1"},
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCategoryRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var respBody struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)

		// Post-conditions: Verify the response content
		assert.NotEmpty(t, respBody.ID)
		assert.Equal(t, "Category created successfully!", respBody.Message)

		// Post-conditions: Verify the category exists in the database
		createdID, _ := uuid.Parse(respBody.ID)
		cat, getErr := td.GetCategoryByID(t, ctx, createdID, testPool)
		assert.NoError(t, getErr, "Failed to retrieve the newly created category from the database")
		assert.Equal(t, strings.ToLower(requestBody.Name), cat.Name)
		assert.Equal(t, requestBody.Description, cat.Description)
		assert.Equal(t, domain.Draft, cat.Status)
		assert.Equal(t, requestBody.Metadata, cat.Metadata)
	})

	t.Run("should return 409 Conflict if category name already exists", func(t *testing.T) {
		// Clean the database for this test
		td.ClearCategoriesTable(t, ctx, testPool)

		// Setup: Insert a category with a known name
		knownName := "Existing Category Name"
		existingCategory := td.NewValidCategory(knownName)
		existingCategory.Name = strings.ToLower(existingCategory.Name) // that operation is done by the domain function, so I have to replicate it when directly inserting to the database
		td.InsertCategory(t, ctx, existingCategory, testPool)

		requestBody := domain.CreateCategoryRequest{
			Name:        knownName, // Use the same name
			Description: "Another description.",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCategoryRequest(t, ctx, jsonBody)
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
		assert.Contains(t, respBody.Error, errs.ErrAlreadyExists.Error())
	})

	t.Run("should return 400 Bad Request for empty name", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)
		requestBody := domain.CreateCategoryRequest{
			Name:        "", // Empty name
			Description: "Valid description.",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCategoryRequest(t, ctx, jsonBody)
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
		assert.Contains(t, respBody.Error, "category name cannot be empty")
	})

	t.Run("should return 400 Bad Request for invalid JSON body", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)
		req := newCreateCategoryRequest(t, ctx, []byte(`{"name": "test", "description": "test", "metadata": "not_an_object"}`))
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

	t.Run("should return 415 Unsupported Media Type for incorrect content type", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)
		requestBody := domain.CreateCategoryRequest{
			Name:        "Test",
			Description: "Test description.",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newCreateCategoryRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "text/plain") // Incorrect content type

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, "unsupported media type")
	})

	t.Run("should return 400 Bad Request for unknown fields", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)
		jsonBody := []byte(`{"name": "Unknown Field Test", "description": "test", "unknown_field": "some_value"}`)

		req := newCreateCategoryRequest(t, ctx, jsonBody)
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
		assert.Contains(t, respBody.Error, "unknown field \"unknown_field\"")
	})
}

// Helper function to create a new HTTP request for the CreateCategory handler.
func newCreateCategoryRequest(t *testing.T, ctx context.Context, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/categories", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}
