package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/Leviosa-care/settings/internal/domain"
	"github.com/stretchr/testify/require"
)

// HTTP Request Helpers for Company Settings

func NewGetCompanyNameRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/settings/name", nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyNameRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetCompanyNameRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/admin/settings/name", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyEmailRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/settings/email", nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyEmailRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetCompanyEmailRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/admin/settings/email", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyPhoneRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/admin/settings/phone", nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyPhoneRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetCompanyTelephoneRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/admin/settings/phone", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyAddressRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/settings/address", nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyAddressRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetCompanyLegalAddressRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/admin/settings/address", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyInstagramRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/settings/instagram", nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyInstagramRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetCompanyInstagramRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/admin/settings/instagram", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyLogoRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/settings/logo", nil)
	require.NoError(t, err)
	return req
}

// HTTP Request Helpers for OTP Settings

func NewGetOTPDurationRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/admin/settings/otp/duration", nil)
	require.NoError(t, err)
	return req
}

func NewSetOTPDurationRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetOTPDurationRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/admin/settings/otp/duration", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetOTPLengthRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/admin/settings/otp/length", nil)
	require.NoError(t, err)
	return req
}

func NewSetOTPLengthRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetOTPLengthRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/admin/settings/otp/length", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetOTPMaxAttemptsRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/admin/settings/otp/max-attempts", nil)
	require.NoError(t, err)
	return req
}

func NewSetOTPMaxAttemptsRequest(t *testing.T, ctx context.Context, serverURL string, request domain.SetOTPMaxAttemptsRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/admin/settings/otp/max-attempts", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// Bulk Settings Request Helper

func NewBulkSettingsRequest(t *testing.T, ctx context.Context, serverURL string, keys []string) *http.Request {
	keysParam := strings.Join(keys, ",")
	url := serverURL + "/settings/bulk?keys=" + keysParam
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	require.NoError(t, err)
	return req
}
