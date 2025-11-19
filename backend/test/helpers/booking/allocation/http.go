package allocationHelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// NewCreateDedicatedAllocationRequest creates an HTTP request for creating a dedicated allocation
func NewCreateDedicatedAllocationRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	request domain.CreateDedicatedAllocationRequest,
	accessToken string,
) *http.Request {
	t.Helper()

	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		serverURL+"/allocations/dedicated",
		bytes.NewReader(jsonBody),
	)
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

// NewCreateSharedAllocationRequest creates an HTTP request for creating a shared allocation
func NewCreateSharedAllocationRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	request domain.CreateSharedAllocationRequest,
	accessToken string,
) *http.Request {
	t.Helper()

	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		serverURL+"/allocations/shared",
		bytes.NewReader(jsonBody),
	)
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

// NewGetAllocationRequest creates an HTTP request for getting an allocation by ID
func NewGetAllocationRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	allocationID uuid.UUID,
	accessToken string,
) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		serverURL+"/allocations/"+allocationID.String(),
		nil,
	)
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

// NewUpdateDedicatedAllocationRequest creates an HTTP request for updating a dedicated allocation
func NewUpdateDedicatedAllocationRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	allocationID uuid.UUID,
	request domain.UpdateDedicatedAllocationRequest,
	accessToken string,
) *http.Request {
	t.Helper()

	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		serverURL+"/allocations/"+allocationID.String()+"/dedicated",
		bytes.NewReader(jsonBody),
	)
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

// NewDeleteAllocationRequest creates an HTTP request for deleting (soft-deleting) an allocation
func NewDeleteAllocationRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	allocationID uuid.UUID,
	accessToken string,
) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		serverURL+"/allocations/"+allocationID.String(),
		nil,
	)
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
