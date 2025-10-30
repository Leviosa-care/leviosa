package coupon_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Helper function to create a new HTTP request for the CreateCoupon handler.
func newCreateCouponRequest(t *testing.T, ctx context.Context, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/coupons", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the ValidateCoupon handler.
func newValidateCouponRequest(t *testing.T, ctx context.Context, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/coupons/validate", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetCouponByID handler.
func newGetCouponByIDRequest(t *testing.T, ctx context.Context, couponID string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/admin/coupons/"+couponID, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetCouponByStripeID handler.
func newGetCouponByStripeIDRequest(t *testing.T, ctx context.Context, stripeID string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/admin/coupons/stripe/"+stripeID, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetAllCoupons handler.
func newGetAllCouponsRequest(t *testing.T, ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/admin/coupons", nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetValidCoupons handler (admin).
func newGetValidCouponsAdminRequest(t *testing.T, ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/admin/coupons/valid", nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetValidCoupons handler (public).
func newGetValidCouponsPublicRequest(t *testing.T, ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/coupons/valid", nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the UpdateCoupon handler.
func newUpdateCouponRequest(t *testing.T, ctx context.Context, couponID string, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, testServerURL+"/admin/coupons/"+couponID, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the DeactivateCoupon handler.
func newDeactivateCouponRequest(t *testing.T, ctx context.Context, couponID string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/coupons/"+couponID+"/deactivate", nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the DeleteCoupon handler.
func newDeleteCouponRequest(t *testing.T, ctx context.Context, couponID string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, testServerURL+"/admin/coupons/"+couponID, nil)
	require.NoError(t, err)
	return req
}

// Helper function to decode JSON response body.
func decodeJSONResponse(t *testing.T, resp *http.Response, target interface{}) {
	t.Helper()
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(target)
	require.NoError(t, err, "Failed to decode JSON response")
}
