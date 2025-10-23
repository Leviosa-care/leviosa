package promotionCode_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Helper function to create a new HTTP request for the CreatePromotionCode handler.
func newCreatePromotionCodeRequest(t *testing.T, ctx context.Context, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/promotion-codes", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the ValidatePromotionCode handler.
func newValidatePromotionCodeRequest(t *testing.T, ctx context.Context, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/promotion-codes/validate", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetPromotionCodeByID handler.
func newGetPromotionCodeByIDRequest(t *testing.T, ctx context.Context, promotionCodeID string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/admin/promotion-codes/"+promotionCodeID, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetPromotionCodeByCode handler.
func newGetPromotionCodeByCodeRequest(t *testing.T, ctx context.Context, code string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/admin/promotion-codes/code/"+code, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetPromotionCodeWithCoupon handler.
func newGetPromotionCodeWithCouponRequest(t *testing.T, ctx context.Context, code string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/promotion-codes/code/"+code, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetAllPromotionCodes handler.
func newGetAllPromotionCodesRequest(t *testing.T, ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/admin/promotion-codes", nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetActivePromotionCodes handler.
func newGetActivePromotionCodesRequest(t *testing.T, ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/admin/promotion-codes/active", nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the UpdatePromotionCode handler.
func newUpdatePromotionCodeRequest(t *testing.T, ctx context.Context, promotionCodeID string, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, testServerURL+"/admin/promotion-codes/"+promotionCodeID, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the DeactivatePromotionCode handler.
func newDeactivatePromotionCodeRequest(t *testing.T, ctx context.Context, promotionCodeID string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/promotion-codes/"+promotionCodeID+"/deactivate", nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the DeletePromotionCode handler.
func newDeletePromotionCodeRequest(t *testing.T, ctx context.Context, promotionCodeID string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, testServerURL+"/admin/promotion-codes/"+promotionCodeID, nil)
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