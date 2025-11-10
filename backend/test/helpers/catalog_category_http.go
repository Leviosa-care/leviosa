package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	categoryHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/category"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Public Category Endpoints (No Auth Required)
// ============================================================================

// NewGetAllPublishedCategoriesRequest creates a request to get all published categories (public endpoint)
func NewGetAllPublishedCategoriesRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+categoryHandler.GetAllPublishedCategoriesEndpoint, nil)
	require.NoError(t, err)
	return req
}

// NewGetCategoryByIDRequest creates a request to get a category by ID (public endpoint)
func NewGetCategoryByIDRequest(t *testing.T, ctx context.Context, serverURL string, categoryID string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+categoryHandler.CategoriesBasePath+"/"+categoryID, nil)
	require.NoError(t, err)
	return req
}

// ============================================================================
// Admin-Only Category Endpoints (Auth Required)
// ============================================================================

// NewGetAdminAllCategoriesRequest creates a request to get all categories including drafts (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetAdminAllCategoriesRequest(t *testing.T, ctx context.Context, serverURL string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+categoryHandler.GetAdminAllCategoriesEndpoint, nil)
	require.NoError(t, err)

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewCreateCategoryRequest creates a request to create a new category (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewCreateCategoryRequest(t *testing.T, ctx context.Context, serverURL string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+categoryHandler.CreateCategoryEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewModifyCategoryRequest creates a request to modify an existing category (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewModifyCategoryRequest(t *testing.T, ctx context.Context, serverURL string, categoryID string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, serverURL+categoryHandler.AdminCategoriesBasePath+"/"+categoryID, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewRemoveCategoryRequest creates a request to remove a category (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewRemoveCategoryRequest(t *testing.T, ctx context.Context, serverURL string, categoryID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, serverURL+categoryHandler.AdminCategoriesBasePath+"/"+categoryID, nil)
	require.NoError(t, err)

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}
