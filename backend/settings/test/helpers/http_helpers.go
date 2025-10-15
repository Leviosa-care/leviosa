package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/Leviosa-care/core/contracts/settings"
	httpEndpoints "github.com/Leviosa-care/settings/internal/adapters/http"
	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/stretchr/testify/require"
)

// HTTP Request Helpers for Company Settings

func NewGetCompanyNameRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.GetCompanyNameEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyNameRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetCompanyNameRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+httpEndpoints.SetCompanyNameEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyEmailRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.GetCompanyEmailEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyEmailRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetCompanyEmailRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+httpEndpoints.SetCompanyEmailEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyPhoneRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.AdminGetCompanyPhoneEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyPhoneRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetCompanyTelephoneRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+httpEndpoints.SetCompanyPhoneEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyAddressRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.GetCompanyAddressEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyAddressRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetCompanyLegalAddressRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+httpEndpoints.SetCompanyAddressEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyInstagramRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.GetCompanyInstagramEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyInstagramRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetCompanyInstagramRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+httpEndpoints.SetCompanyInstagramEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyLogoRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.GetCompanyLogoEndpoint, nil)
	require.NoError(t, err)
	return req
}

// HTTP Request Helpers for OTP Settings

func NewGetOTPDurationRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.AdminGetOTPDurationEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetOTPDurationRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetOTPDurationRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+httpEndpoints.AdminSetOTPDurationEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetOTPLengthRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.AdminGetOTPLengthEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetOTPLengthRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetOTPLengthRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+httpEndpoints.AdminSetOTPLengthEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetOTPMaxAttemptsRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.AdminGetOTPMaxAttemptsEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetOTPMaxAttemptsRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetOTPMaxAttemptsRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+httpEndpoints.AdminSetOTPMaxAttemptsEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// Bulk Settings Request Helper
func NewBulkSettingsRequest(t *testing.T, ctx context.Context, serverURL string, keys []string) *http.Request {
	paramkeys := strings.Join(keys, ",")
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		serverURL+httpEndpoints.AdminBulkEndpoint+"?keys="+paramkeys,
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")
	return req
}

// Internal Service-to-Service Endpoint Helpers

func NewInternalGetCompanyNameRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.InternalGetCompanyNameEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewInternalGetCompanyEmailRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.InternalGetCompanyEmailEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewInternalGetCompanyPhoneRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.InternalGetCompanyPhoneEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewInternalGetCompanyAddressRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.InternalGetCompanyAddressEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewInternalGetOTPDurationRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+httpEndpoints.InternalGetOTPDurationEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewInternalBulkSettingsRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	url := httpEndpoints.InternalBulkEndpoint + "?keys=" + settings.CompanyLogo
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+url, nil)
	require.NoError(t, err)
	return req
}
