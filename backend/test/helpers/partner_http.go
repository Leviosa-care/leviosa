package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	partnerEndpoints "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/partner"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// NewGetPartnerByIDRequest creates an HTTP request for getting a partner by ID
func NewGetPartnerByIDRequest(t *testing.T, ctx context.Context, serverURL string, partnerID uuid.UUID) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s%s/%s", serverURL, partnerEndpoints.PartnersBasePath, partnerID.String())
	println("the URL is :", url)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

	require.NoError(t, err, "Failed to create HTTP request")

	return req
}

// NewGetPartnerMeRequest creates an HTTP request for getting authenticated partner's own profile
func NewGetPartnerMeRequest(t *testing.T, ctx context.Context, serverURL string, accessToken string) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s%s", serverURL, partnerEndpoints.GetPartnerMeEndpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewGetAllPartnersRequest creates an HTTP request for getting all partners
func NewGetAllPartnersRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	t.Helper()

	url := serverURL + partnerEndpoints.GetAllPartnersEndpoint
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	return req
}

// NewUpdatePartnerRequest creates an HTTP request for updating a partner
func NewUpdatePartnerRequest(t *testing.T, ctx context.Context, serverURL string, partnerID uuid.UUID, request domain.UpdatePartnerRequest, accessToken string) *http.Request {
	t.Helper()

	body, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal update partner request")

	url := fmt.Sprintf("%s/partners/%s", serverURL, partnerID.String())
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(body))
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewDeletePartnerRequest creates an HTTP request for deleting a partner
func NewDeletePartnerRequest(t *testing.T, ctx context.Context, serverURL string, partnerID uuid.UUID, accessToken string) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s/admin/partners/%s", serverURL, partnerID.String())
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewVerifyPartnerRequest creates an HTTP request for verifying a partner
func NewVerifyPartnerRequest(t *testing.T, ctx context.Context, serverURL string, partnerID uuid.UUID, accessToken string) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s/admin/partners/%s/verify", serverURL, partnerID.String())
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	// Add session cookie if provided
	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// ParsePartnerResponse parses a single partner from HTTP response body
func ParsePartnerResponse(t *testing.T, resp *http.Response) *domain.PartnerResponse {
	t.Helper()

	var partner domain.PartnerResponse
	err := json.NewDecoder(resp.Body).Decode(&partner)
	require.NoError(t, err, "Failed to parse partner response")

	return &partner
}

// ParsePartnersListResponse parses a list of partners from HTTP response body
func ParsePartnersListResponse(t *testing.T, resp *http.Response) []*domain.PartnerResponse {
	t.Helper()

	var partners []*domain.PartnerResponse
	err := json.NewDecoder(resp.Body).Decode(&partners)
	require.NoError(t, err, "Failed to parse partners list response")

	return partners
}

// NewGetAllPartnersByCategoryRequest creates an HTTP request for getting all partners by category ID
func NewGetAllPartnersByCategoryRequest(t *testing.T, ctx context.Context, serverURL string, categoryID uuid.UUID) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s%s/%s", serverURL, partnerEndpoints.PartnersBasePath+partnerEndpoints.CategoryPath, categoryID.String())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	return req
}

// NewGetAllPartnersByCategoryRequestWithInvalidID creates an HTTP request with an invalid category ID for testing error handling
func NewGetAllPartnersByCategoryRequestWithInvalidID(t *testing.T, ctx context.Context, serverURL string, invalidID string) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s%s/%s", serverURL, partnerEndpoints.PartnersBasePath+partnerEndpoints.CategoryPath, invalidID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	return req
}

// NewGetAllPartnersByCategoriesRequest creates an HTTP request for getting all partners by multiple category IDs
func NewGetAllPartnersByCategoriesRequest(t *testing.T, ctx context.Context, serverURL string, categoryIDs []uuid.UUID) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s%s", serverURL, partnerEndpoints.PartnersBasePath+partnerEndpoints.CategoryPath)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	// Add category_id query parameters
	q := req.URL.Query()
	for _, categoryID := range categoryIDs {
		q.Add("category_id", categoryID.String())
	}
	req.URL.RawQuery = q.Encode()

	return req
}

// NewGetAllPartnersByCategoriesRequestWithStrings creates an HTTP request with string category IDs for testing
func NewGetAllPartnersByCategoriesRequestWithStrings(t *testing.T, ctx context.Context, serverURL string, categoryIDs []string) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s%s", serverURL, partnerEndpoints.PartnersBasePath+partnerEndpoints.CategoryPath)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	// Add category_id query parameters
	q := req.URL.Query()
	for _, categoryID := range categoryIDs {
		q.Add("category_id", categoryID)
	}
	req.URL.RawQuery = q.Encode()

	return req
}
