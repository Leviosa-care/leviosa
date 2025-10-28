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

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// NewCreatePartnerRequest creates an HTTP request for creating a partner
func NewCreatePartnerRequest(t *testing.T, ctx context.Context, serverURL string, request domain.CreatePartnerRequest) *http.Request {
	t.Helper()

	body, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal create partner request")

	req, err := http.NewRequestWithContext(ctx, "POST", serverURL+partnerEndpoints.CreatePartnerEndpoint, bytes.NewReader(body))
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewGetPartnerByIDRequest creates an HTTP request for getting a partner by ID
func NewGetPartnerByIDRequest(t *testing.T, ctx context.Context, serverURL string, partnerID uuid.UUID) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s/admin/partners/%s", serverURL, partnerID.String())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	return req
}

// NewGetPartnerByUserIDRequest creates an HTTP request for getting a partner by user ID
func NewGetPartnerByUserIDRequest(t *testing.T, ctx context.Context, serverURL string, userID uuid.UUID) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s/admin/partners/user/%s", serverURL, userID.String())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

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
func NewUpdatePartnerRequest(t *testing.T, ctx context.Context, serverURL string, partnerID uuid.UUID, request domain.UpdatePartnerRequest) *http.Request {
	t.Helper()

	body, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal update partner request")

	url := fmt.Sprintf("%s/partners/%s", serverURL, partnerID.String())
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(body))
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewDeletePartnerRequest creates an HTTP request for deleting a partner
func NewDeletePartnerRequest(t *testing.T, ctx context.Context, serverURL string, partnerID uuid.UUID) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s/admin/partners/%s", serverURL, partnerID.String())
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	return req
}

// NewVerifyPartnerRequest creates an HTTP request for verifying a partner
func NewVerifyPartnerRequest(t *testing.T, ctx context.Context, serverURL string, partnerID uuid.UUID) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s/admin/partners/%s/verify", serverURL, partnerID.String())
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	return req
}