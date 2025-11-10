package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	productHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/product"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Public Product Endpoints (No Auth Required)
// ============================================================================

// NewGetAllPublishedProductsRequest creates a request to get all published products (public endpoint)
func NewGetAllPublishedProductsRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+productHandler.GetAllPublishedProductsEndpoint, nil)
	require.NoError(t, err)
	return req
}

// NewGetProductByIDRequest creates a request to get a product by ID (public endpoint)
func NewGetProductByIDRequest(t *testing.T, ctx context.Context, serverURL string, productID string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+productHandler.ProductsBasePath+"/"+productID, nil)
	require.NoError(t, err)
	return req
}

// ============================================================================
// Admin-Only Product Endpoints (Auth Required)
// ============================================================================

// NewGetAdminAllProductsRequest creates a request to get all products including drafts (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetAdminAllProductsRequest(t *testing.T, ctx context.Context, serverURL string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+productHandler.GetAdminAllProductsEndpoint, nil)
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

// NewCreateProductWithPriceRequest creates a request to create a new product with price (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewCreateProductWithPriceRequest(t *testing.T, ctx context.Context, serverURL string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+productHandler.CreateProductWithPriceEndpoint, bytes.NewReader(jsonBody))
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

// NewModifyProductRequest creates a request to modify an existing product (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewModifyProductRequest(t *testing.T, ctx context.Context, serverURL string, productID string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, serverURL+productHandler.AdminProductsBasePath+"/"+productID, bytes.NewReader(jsonBody))
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

// NewRemoveProductRequest creates a request to remove a product (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewRemoveProductRequest(t *testing.T, ctx context.Context, serverURL string, productID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, serverURL+productHandler.AdminProductsBasePath+"/"+productID, nil)
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
