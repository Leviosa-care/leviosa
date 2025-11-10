package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	promotionCodeHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/promotion_code"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Public Promotion Code Endpoints (No Auth Required)
// ============================================================================

// NewValidatePromotionCodeRequest creates a request to validate a promotion code (public endpoint for checkout)
func NewValidatePromotionCodeRequest(t *testing.T, ctx context.Context, serverURL string, requestBody interface{}) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+promotionCodeHandler.ValidatePromotionCodeEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewGetPromotionCodeWithCouponRequest creates a request to get a promotion code with coupon by code (public endpoint)
func NewGetPromotionCodeWithCouponRequest(t *testing.T, ctx context.Context, serverURL string, code string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+promotionCodeHandler.PromotionCodesBasePath+promotionCodeHandler.CodePath+"/"+code, nil)
	require.NoError(t, err)
	return req
}

// ============================================================================
// Admin-Only Promotion Code Endpoints (Auth Required)
// ============================================================================

// NewGetAllPromotionCodesRequest creates a request to get all promotion codes (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetAllPromotionCodesRequest(t *testing.T, ctx context.Context, serverURL string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+promotionCodeHandler.GetAllPromotionCodesEndpoint, nil)
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

// NewGetActivePromotionCodesRequest creates a request to get active promotion codes (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetActivePromotionCodesRequest(t *testing.T, ctx context.Context, serverURL string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+promotionCodeHandler.GetActivePromotionCodesEndpoint, nil)
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

// NewGetPromotionCodeByIDRequest creates a request to get a promotion code by ID (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetPromotionCodeByIDRequest(t *testing.T, ctx context.Context, serverURL string, promotionCodeID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+promotionCodeHandler.AdminPromotionCodesBasePath+"/"+promotionCodeID, nil)
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

// NewGetPromotionCodeByCodeRequest creates a request to get a promotion code by code string (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewGetPromotionCodeByCodeRequest(t *testing.T, ctx context.Context, serverURL string, code string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+promotionCodeHandler.AdminPromotionCodesBasePath+promotionCodeHandler.CodePath+"/"+code, nil)
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

// NewCreatePromotionCodeRequest creates a request to create a new promotion code (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewCreatePromotionCodeRequest(t *testing.T, ctx context.Context, serverURL string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+promotionCodeHandler.CreatePromotionCodeEndpoint, bytes.NewReader(jsonBody))
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

// NewUpdatePromotionCodeRequest creates a request to update a promotion code (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewUpdatePromotionCodeRequest(t *testing.T, ctx context.Context, serverURL string, promotionCodeID string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, serverURL+promotionCodeHandler.AdminPromotionCodesBasePath+"/"+promotionCodeID, bytes.NewReader(jsonBody))
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

// NewDeactivatePromotionCodeRequest creates a request to deactivate a promotion code (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewDeactivatePromotionCodeRequest(t *testing.T, ctx context.Context, serverURL string, promotionCodeID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+promotionCodeHandler.AdminPromotionCodesBasePath+"/"+promotionCodeID+promotionCodeHandler.DeactivatePath, nil)
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

// NewDeletePromotionCodeRequest creates a request to delete a promotion code (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewDeletePromotionCodeRequest(t *testing.T, ctx context.Context, serverURL string, promotionCodeID string, accessToken string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, serverURL+promotionCodeHandler.AdminPromotionCodesBasePath+"/"+promotionCodeID, nil)
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

// ParseValidatePromotionCodeResponse parses a promotion code validation response from HTTP response body
func ParseValidatePromotionCodeResponse(t *testing.T, resp *http.Response) *domain.ValidatePromotionCodeResponse {
	t.Helper()
	var response domain.ValidatePromotionCodeResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode promotion code validation response")
	return &response
}

// ParsePromotionCodeResponse parses a single promotion code response from HTTP response body
func ParsePromotionCodeResponse(t *testing.T, resp *http.Response) *domain.PromotionCodeWithCouponResponse {
	t.Helper()
	var response domain.PromotionCodeWithCouponResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode promotion code response")
	return &response
}

// ParsePromotionCodesResponse parses a list of promotion codes from HTTP response body
func ParsePromotionCodesResponse(t *testing.T, resp *http.Response) []*domain.PromotionCodeWithCouponResponse {
	t.Helper()
	var responses []*domain.PromotionCodeWithCouponResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&responses)
	require.NoError(t, err, "Failed to decode promotion codes response")
	return responses
}

// ParseSimplePromotionCodeResponse parses a single simple promotion code response from HTTP response body
func ParseSimplePromotionCodeResponse(t *testing.T, resp *http.Response) *domain.PromotionCodeResponse {
	t.Helper()
	var response domain.PromotionCodeResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&response)
	require.NoError(t, err, "Failed to decode simple promotion code response")
	return &response
}

// ParseSimplePromotionCodesResponse parses a list of simple promotion codes from HTTP response body
func ParseSimplePromotionCodesResponse(t *testing.T, resp *http.Response) []*domain.PromotionCodeResponse {
	t.Helper()
	var responses []*domain.PromotionCodeResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&responses)
	require.NoError(t, err, "Failed to decode simple promotion codes response")
	return responses
}
