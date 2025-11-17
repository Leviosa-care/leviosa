package category_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	categoryHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/category"
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

// make test-func TEST_NAME=TestModifyCategory TEST_PATH=test/integration/catalog/category/modify_category_test.go

func TestModifyCategory(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully modify a category's name and description with valid admin token", func(t *testing.T) {
		// Setup: Clear the table and insert a known category for this test
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		existingCategory := td.NewValidCategory("Original Name")
		td.InsertCategory(t, ctx, existingCategory, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		updatedName := "Updated Category Name"
		updatedDescription := "Updated description."

		requestBody := domain.UpdateCategoryRequest{
			Name:        &updatedName,
			Description: &updatedDescription,
		}

		req := th.NewModifyCategoryRequest(t, ctx, testServerURL, existingCategory.ID.String(), requestBody, accessToken)

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
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		nonExistentID := uuid.New().String()
		updatedName := "New Name"
		requestBody := domain.UpdateCategoryRequest{Name: &updatedName}

		req := th.NewModifyCategoryRequest(t, ctx, testServerURL, nonExistentID, requestBody, accessToken)

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

	t.Run("should return 409 Conflict if category name already exists", func(t *testing.T) {
		// Clean and setup for this specific test
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		existingCat1 := td.NewValidCategory("Unique Category 1")
		td.InsertCategory(t, ctx, existingCat1, testPool)

		existingCat2 := td.NewValidCategory("Unique Category 2")
		td.InsertCategory(t, ctx, existingCat2, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Attempt to update cat1's name to cat2's name
		requestBody := domain.UpdateCategoryRequest{
			Name: &existingCat2.Name,
		}

		req := th.NewModifyCategoryRequest(t, ctx, testServerURL, existingCat1.ID.String(), requestBody, accessToken)

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
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		updatedName := "New Name"
		requestBody := domain.UpdateCategoryRequest{Name: &updatedName}

		req := th.NewModifyCategoryRequest(t, ctx, testServerURL, "not-a-valid-uuid", requestBody, accessToken)

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
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		existingCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, existingCategory, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		requestBody := domain.UpdateCategoryRequest{}

		req := th.NewModifyCategoryRequest(t, ctx, testServerURL, existingCategory.ID.String(), requestBody, accessToken)

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
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		existingCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, existingCategory, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create request with invalid JSON by manually crafting the request body
		jsonBody := []byte(`{"name": 123}`)

		url := testServerURL + categoryHandler.AdminCategoriesBasePath + "/" + existingCategory.ID.String()
		req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)

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
		td.ClearCategoriesTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		existingCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, existingCategory, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create request with unknown field by manually crafting the request body
		jsonBody := []byte(`{"unknown_field": "some_value"}`)

		url := testServerURL + categoryHandler.AdminCategoriesBasePath + "/" + existingCategory.ID.String()
		req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)

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

	// Authentication test cases
	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)

		existingCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, existingCategory, testPool)

		updatedName := "Updated Category Name"
		requestBody := domain.UpdateCategoryRequest{Name: &updatedName}

		req := th.NewModifyCategoryRequest(t, ctx, testServerURL, existingCategory.ID.String(), requestBody, "")

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

		updatedName := "Updated Category Name"
		requestBody := domain.UpdateCategoryRequest{Name: &updatedName}

		req := th.NewModifyCategoryRequest(t, ctx, testServerURL, existingCategory.ID.String(), requestBody, accessToken)

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

		updatedName := "Updated Category Name"
		requestBody := domain.UpdateCategoryRequest{Name: &updatedName}

		req := th.NewModifyCategoryRequest(t, ctx, testServerURL, existingCategory.ID.String(), requestBody, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		td.ClearCategoriesTable(t, ctx, testPool)

		existingCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, existingCategory, testPool)

		updatedName := "Updated Category Name"
		requestBody := domain.UpdateCategoryRequest{Name: &updatedName}

		req := th.NewModifyCategoryRequest(t, ctx, testServerURL, existingCategory.ID.String(), requestBody, "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
