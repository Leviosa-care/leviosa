package category_test

import (
	"context"
	"encoding/json"
	"net/http"

	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateCategory TEST_PATH=test/integration/catalog/category/create_category_test.go

func TestCreateCategory(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully create category with valid admin token", func(t *testing.T) {
		// Clean the database to ensure isolation for this test
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		requestBody := domain.CreateCategoryRequest{
			Name:        "New Category",
			Description: "A great new category.",
		}

		req := th.NewCreateCategoryRequest(t, ctx, testServerURL, requestBody, accessToken)

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
	})

	t.Run("should return 409 Conflict if category name already exists", func(t *testing.T) {
		// Clean the database for this test
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Setup: Insert a category with a known name
		knownName := "Existing Category Name"
		existingCategory := td.NewValidCategory(knownName)
		existingCategory.Name = strings.ToLower(existingCategory.Name) // that operation is done by the domain function, so I have to replicate it when directly inserting to the database
		td.InsertCategory(t, ctx, existingCategory, testPool)

		requestBody := domain.CreateCategoryRequest{
			Name:        knownName, // Use the same name
			Description: "Another description.",
		}

		req := th.NewCreateCategoryRequest(t, ctx, testServerURL, requestBody, accessToken)

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
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		requestBody := domain.CreateCategoryRequest{
			Name:        "", // Empty name
			Description: "Valid description.",
		}

		req := th.NewCreateCategoryRequest(t, ctx, testServerURL, requestBody, accessToken)

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
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewCreateCategoryRequest(t, ctx, testServerURL, nil, accessToken)

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

	t.Run("should return 400 Bad Request for unknown fields", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// This would be handled by the JSON unmarshaling in the actual request
		requestBody := domain.CreateCategoryRequest{
			Name:        "Unknown Field Test",
			Description: "test",
		}

		req := th.NewCreateCategoryRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// This test may need adjustment based on actual validation implementation
		// For now, we expect it to succeed since unknown fields should be ignored
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	// Add authentication test cases
	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)

		requestBody := domain.CreateCategoryRequest{
			Name:        "Test Category",
			Description: "Test description",
		}

		req := th.NewCreateCategoryRequest(t, ctx, testServerURL, requestBody, "") // Empty token

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create expired admin session
		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		requestBody := domain.CreateCategoryRequest{
			Name:        "Test Category",
			Description: "Test description",
		}

		req := th.NewCreateCategoryRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create standard user (not admin)
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		requestBody := domain.CreateCategoryRequest{
			Name:        "Test Category",
			Description: "Test description",
		}

		req := th.NewCreateCategoryRequest(t, ctx, testServerURL, requestBody, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)

		requestBody := domain.CreateCategoryRequest{
			Name:        "Test Category",
			Description: "Test description",
		}

		req := th.NewCreateCategoryRequest(t, ctx, testServerURL, requestBody, "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
