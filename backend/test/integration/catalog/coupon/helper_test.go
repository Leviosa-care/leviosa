package coupon_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	couponHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/coupon"

	"github.com/stretchr/testify/require"
)

// buildEndpointWithID replaces the {id} placeholder with the actual ID
func buildEndpointWithID(template, id string) string {
	return strings.Replace(template, "{id}", id, 1)
}

// buildEndpointWithStripeID replaces the {stripeId} placeholder with the actual Stripe ID
func buildEndpointWithStripeID(template, stripeID string) string {
	return strings.Replace(template, "{stripeId}", stripeID, 1)
}

// Helper function to create a new HTTP request for the CreateCoupon handler.
func newCreateCouponRequest(t *testing.T, ctx context.Context, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+couponHandler.CreateCouponEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the ValidateCoupon handler.
func newValidateCouponRequest(t *testing.T, ctx context.Context, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+couponHandler.ValidateCouponEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetCouponByID handler.
func newGetCouponByIDRequest(t *testing.T, ctx context.Context, couponID string) *http.Request {
	endpoint := buildEndpointWithID(couponHandler.GetCouponByIDEndpoint, couponID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+endpoint, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetCouponByStripeID handler.
func newGetCouponByStripeIDRequest(t *testing.T, ctx context.Context, stripeID string) *http.Request {
	endpoint := buildEndpointWithStripeID(couponHandler.GetCouponByStripeIDEndpoint, stripeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+endpoint, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetAllCoupons handler.
func newGetAllCouponsRequest(t *testing.T, ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+couponHandler.GetAllCouponsEndpoint, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetValidCoupons handler (admin).
func newGetValidCouponsAdminRequest(t *testing.T, ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+couponHandler.GetAllCouponsEndpoint+"/valid", nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetValidCoupons handler (public).
func newGetValidCouponsPublicRequest(t *testing.T, ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+couponHandler.GetValidCouponsEndpoint, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the UpdateCoupon handler.
func newUpdateCouponRequest(t *testing.T, ctx context.Context, couponID string, jsonBody []byte) *http.Request {
	endpoint := buildEndpointWithID(couponHandler.UpdateCouponEndpoint, couponID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, testServerURL+endpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the DeactivateCoupon handler.
func newDeactivateCouponRequest(t *testing.T, ctx context.Context, couponID string) *http.Request {
	endpoint := buildEndpointWithID(couponHandler.DeactivateCouponEndpoint, couponID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+endpoint, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the DeleteCoupon handler.
func newDeleteCouponRequest(t *testing.T, ctx context.Context, couponID string) *http.Request {
	endpoint := buildEndpointWithID(couponHandler.DeleteCouponEndpoint, couponID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, testServerURL+endpoint, nil)
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
