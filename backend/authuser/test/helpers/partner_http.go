package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// NewCreatePartnerRequest creates an HTTP request for creating a partner
func NewCreatePartnerRequest(t *testing.T, ctx context.Context, serverURL string, request domain.CreatePartnerRequest) *http.Request {
	t.Helper()

	body, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal create partner request")

	req, err := http.NewRequestWithContext(ctx, "POST", serverURL+"/partners", bytes.NewReader(body))
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

	url := fmt.Sprintf("%s/admin/partners", serverURL)
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

// NewAddPartnerSpecializationRequest creates an HTTP request for adding a specialization to a partner
func NewAddPartnerSpecializationRequest(t *testing.T, ctx context.Context, serverURL string, partnerID, specializationID uuid.UUID) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s/admin/partners/%s/specializations/%s", serverURL, partnerID.String(), specializationID.String())
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	return req
}

// NewRemovePartnerSpecializationRequest creates an HTTP request for removing a specialization from a partner
func NewRemovePartnerSpecializationRequest(t *testing.T, ctx context.Context, serverURL string, partnerID, specializationID uuid.UUID) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s/admin/partners/%s/specializations/%s", serverURL, partnerID.String(), specializationID.String())
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	return req
}

// NewGetPartnerSpecializationsRequest creates an HTTP request for getting partner specializations
func NewGetPartnerSpecializationsRequest(t *testing.T, ctx context.Context, serverURL string, partnerID uuid.UUID) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s/partners/%s/specializations", serverURL, partnerID.String())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	return req
}

// NewCreateSpecializationRequest creates an HTTP request for creating a specialization
func NewCreateSpecializationRequest(t *testing.T, ctx context.Context, serverURL string, request domain.CreateSpecializationRequest) *http.Request {
	t.Helper()

	body, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal create specialization request")

	req, err := http.NewRequestWithContext(ctx, "POST", serverURL+"/admin/specializations", bytes.NewReader(body))
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewGetSpecializationByIDRequest creates an HTTP request for getting a specialization by ID
func NewGetSpecializationByIDRequest(t *testing.T, ctx context.Context, serverURL string, specializationID uuid.UUID) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s/admin/specializations/%s", serverURL, specializationID.String())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	return req
}

// NewGetAllSpecializationsRequest creates an HTTP request for getting all specializations
func NewGetAllSpecializationsRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s/specializations", serverURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	return req
}

// NewUpdateSpecializationRequest creates an HTTP request for updating a specialization
func NewUpdateSpecializationRequest(t *testing.T, ctx context.Context, serverURL string, specializationID uuid.UUID, request domain.UpdateSpecializationRequest) *http.Request {
	t.Helper()

	body, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal update specialization request")

	url := fmt.Sprintf("%s/admin/specializations/%s", serverURL, specializationID.String())
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(body))
	require.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewDeleteSpecializationRequest creates an HTTP request for deleting a specialization
func NewDeleteSpecializationRequest(t *testing.T, ctx context.Context, serverURL string, specializationID uuid.UUID) *http.Request {
	t.Helper()

	url := fmt.Sprintf("%s/admin/specializations/%s", serverURL, specializationID.String())
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	return req
}