package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	priceHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/price"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Admin-Only Price Endpoints (All price endpoints require auth)
// ============================================================================

// NewGetPriceRequest creates a request to get a price by ID (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetPriceRequest(t *testing.T, ctx context.Context, serverURL string, priceID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+priceHandler.AdminPricesBasePath+"/"+priceID, nil)
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

// NewGetPricesByProductIDRequest creates a request to get all prices for a product (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetPricesByProductIDRequest(t *testing.T, ctx context.Context, serverURL string, productID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+priceHandler.AdminProductsBasePath+"/"+productID+priceHandler.PricesPath, nil)
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

// NewCreatePriceRequest creates a request to create a new price for a product (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewCreatePriceRequest(t *testing.T, ctx context.Context, serverURL string, productID string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+priceHandler.AdminProductsBasePath+"/"+productID+priceHandler.PricesPath, bytes.NewReader(jsonBody))
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

// NewUpdatePriceRequest creates a request to update a price (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewUpdatePriceRequest(t *testing.T, ctx context.Context, serverURL string, priceID string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, serverURL+priceHandler.AdminPricesBasePath+"/"+priceID, bytes.NewReader(jsonBody))
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
