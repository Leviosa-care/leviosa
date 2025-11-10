package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	couponHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/coupon"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Public Coupon Endpoints (No Auth Required)
// ============================================================================

// NewValidateCouponRequest creates a request to validate a coupon code (public endpoint for checkout)
func NewValidateCouponRequest(t *testing.T, ctx context.Context, serverURL string, requestBody interface{}) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+couponHandler.ValidateCouponEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewGetValidCouponsRequest creates a request to get all valid/active coupons (public endpoint)
func NewGetValidCouponsRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+couponHandler.GetValidCouponsEndpoint, nil)
	require.NoError(t, err)
	return req
}

// ============================================================================
// Admin-Only Coupon Endpoints (Auth Required)
// ============================================================================

// NewGetAllCouponsRequest creates a request to get all coupons (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetAllCouponsRequest(t *testing.T, ctx context.Context, serverURL string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+couponHandler.GetAllCouponsEndpoint, nil)
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

// NewGetCouponByIDRequest creates a request to get a coupon by ID (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetCouponByIDRequest(t *testing.T, ctx context.Context, serverURL string, couponID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+couponHandler.AdminCouponsBasePath+"/"+couponID, nil)
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

// NewGetCouponByStripeIDRequest creates a request to get a coupon by Stripe ID (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetCouponByStripeIDRequest(t *testing.T, ctx context.Context, serverURL string, stripeID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+couponHandler.AdminCouponsBasePath+couponHandler.StripePath+"/"+stripeID, nil)
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

// NewCreateCouponRequest creates a request to create a new coupon (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewCreateCouponRequest(t *testing.T, ctx context.Context, serverURL string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+couponHandler.CreateCouponEndpoint, bytes.NewReader(jsonBody))
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

// NewUpdateCouponRequest creates a request to update a coupon (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewUpdateCouponRequest(t *testing.T, ctx context.Context, serverURL string, couponID string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, serverURL+couponHandler.AdminCouponsBasePath+"/"+couponID, bytes.NewReader(jsonBody))
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

// NewDeactivateCouponRequest creates a request to deactivate a coupon (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewDeactivateCouponRequest(t *testing.T, ctx context.Context, serverURL string, couponID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+couponHandler.AdminCouponsBasePath+"/"+couponID+couponHandler.DeactivatePath, nil)
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

// NewDeleteCouponRequest creates a request to delete a coupon (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewDeleteCouponRequest(t *testing.T, ctx context.Context, serverURL string, couponID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, serverURL+couponHandler.AdminCouponsBasePath+"/"+couponID, nil)
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

// ============================================================================
// Response Parsing Helpers
// ============================================================================

// ParseCouponResponse parses a single coupon response from HTTP response body
func ParseCouponResponse(t *testing.T, resp *http.Response) *domain.CouponResponse {
	t.Helper()
	var coupon domain.CouponResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&coupon)
	require.NoError(t, err, "Failed to decode coupon response")
	return &coupon
}

// ParseCouponsResponse parses a list of coupons from HTTP response body
func ParseCouponsResponse(t *testing.T, resp *http.Response) []*domain.CouponResponse {
	t.Helper()
	var coupons []*domain.CouponResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&coupons)
	require.NoError(t, err, "Failed to decode coupons response")
	return coupons
}
