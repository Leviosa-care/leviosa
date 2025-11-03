package promotionCode_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	promotionCodeHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/promotion_code"

	"github.com/stretchr/testify/require"
)

// buildEndpointWithID replaces the {id} placeholder with the actual ID
func buildEndpointWithID(template, id string) string {
	return strings.Replace(template, "{id}", id, 1)
}

// buildEndpointWithCode replaces the {code} placeholder with the actual code
func buildEndpointWithCode(template, code string) string {
	return strings.Replace(template, "{code}", code, 1)
}

// Helper function to create a new HTTP request for the CreatePromotionCode handler.
func newCreatePromotionCodeRequest(t *testing.T, ctx context.Context, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+promotionCodeHandler.CreatePromotionCodeEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the ValidatePromotionCode handler.
func newValidatePromotionCodeRequest(t *testing.T, ctx context.Context, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+promotionCodeHandler.ValidatePromotionCodeEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetPromotionCodeByID handler.
func newGetPromotionCodeByIDRequest(t *testing.T, ctx context.Context, promotionCodeID string) *http.Request {
	endpoint := buildEndpointWithID(promotionCodeHandler.GetPromotionCodeByIDEndpoint, promotionCodeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+endpoint, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetPromotionCodeByCode handler.
func newGetPromotionCodeByCodeRequest(t *testing.T, ctx context.Context, code string) *http.Request {
	endpoint := buildEndpointWithCode(promotionCodeHandler.GetPromotionCodeByCodeEndpoint, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+endpoint, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetPromotionCodeWithCoupon handler.
func newGetPromotionCodeWithCouponRequest(t *testing.T, ctx context.Context, code string) *http.Request {
	endpoint := buildEndpointWithCode(promotionCodeHandler.GetPromotionCodeWithCouponEndpoint, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+endpoint, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetAllPromotionCodes handler.
func newGetAllPromotionCodesRequest(t *testing.T, ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+promotionCodeHandler.GetAllPromotionCodesEndpoint, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the GetActivePromotionCodes handler.
func newGetActivePromotionCodesRequest(t *testing.T, ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+promotionCodeHandler.GetActivePromotionCodesEndpoint, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the UpdatePromotionCode handler.
func newUpdatePromotionCodeRequest(t *testing.T, ctx context.Context, promotionCodeID string, jsonBody []byte) *http.Request {
	endpoint := buildEndpointWithID(promotionCodeHandler.UpdatePromotionCodeEndpoint, promotionCodeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, testServerURL+endpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the DeactivatePromotionCode handler.
func newDeactivatePromotionCodeRequest(t *testing.T, ctx context.Context, promotionCodeID string) *http.Request {
	endpoint := buildEndpointWithID(promotionCodeHandler.DeactivatePromotionCodeEndpoint, promotionCodeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+endpoint, nil)
	require.NoError(t, err)
	return req
}

// Helper function to create a new HTTP request for the DeletePromotionCode handler.
func newDeletePromotionCodeRequest(t *testing.T, ctx context.Context, promotionCodeID string) *http.Request {
	endpoint := buildEndpointWithID(promotionCodeHandler.DeletePromotionCodeEndpoint, promotionCodeID)
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