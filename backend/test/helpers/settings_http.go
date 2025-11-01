package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	settingsDomain "github.com/Leviosa-care/leviosa/backend/internal/settings/domain"
	settingsEndpoints "github.com/Leviosa-care/leviosa/backend/internal/settings/interface/http"
	"github.com/stretchr/testify/require"
)

func NewGetCompanyNameRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.GetCompanyNameEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyNameRequest(t *testing.T, ctx context.Context, serverURL string, request settingsDomain.SetCompanyNameRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+settingsEndpoints.SetCompanyNameEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyEmailRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.GetCompanyEmailEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyEmailRequest(t *testing.T, ctx context.Context, serverURL string, request settingsDomain.SetCompanyEmailRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+settingsEndpoints.SetCompanyEmailEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyPhoneRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.AdminGetCompanyPhoneEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyPhoneRequest(t *testing.T, ctx context.Context, serverURL string, request settingsDomain.SetCompanyTelephoneRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+settingsEndpoints.SetCompanyPhoneEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyAddressRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.GetCompanyAddressEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyAddressRequest(t *testing.T, ctx context.Context, serverURL string, request settingsDomain.SetCompanyLegalAddressRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+settingsEndpoints.SetCompanyAddressEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyInstagramRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.GetCompanyInstagramEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetCompanyInstagramRequest(t *testing.T, ctx context.Context, serverURL string, request settingsDomain.SetCompanyInstagramRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+settingsEndpoints.SetCompanyInstagramEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetCompanyLogoRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.GetCompanyLogoEndpoint, nil)
	require.NoError(t, err)
	return req
}

// HTTP Request Helpers for OTP Settings

func NewGetOTPDurationRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.AdminGetOTPDurationEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetOTPDurationRequest(t *testing.T, ctx context.Context, serverURL string, request settingsDomain.SetOTPDurationRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+settingsEndpoints.AdminSetOTPDurationEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetOTPLengthRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.AdminGetOTPLengthEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetOTPLengthRequest(t *testing.T, ctx context.Context, serverURL string, request settingsDomain.SetOTPLengthRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+settingsEndpoints.AdminSetOTPLengthEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewGetOTPMaxAttemptsRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.AdminGetOTPMaxAttemptsEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewSetOTPMaxAttemptsRequest(t *testing.T, ctx context.Context, serverURL string, request settingsDomain.SetOTPMaxAttemptsRequest) *http.Request {
	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+settingsEndpoints.AdminSetOTPMaxAttemptsEndpoint, bytes.NewReader(jsonBody))
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
		serverURL+settingsEndpoints.AdminBulkEndpoint+"?keys="+paramkeys,
		nil,
	)
	require.NoError(t, err, "Failed to create HTTP request")
	return req
}

// Internal Service-to-Service Endpoint Helpers

func NewInternalGetCompanyNameRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.InternalGetCompanyNameEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewInternalGetCompanyEmailRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.InternalGetCompanyEmailEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewInternalGetCompanyPhoneRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.InternalGetCompanyPhoneEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewInternalGetCompanyAddressRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.InternalGetCompanyAddressEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewInternalGetOTPDurationRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+settingsEndpoints.InternalGetOTPDurationEndpoint, nil)
	require.NoError(t, err)
	return req
}

func NewInternalBulkSettingsRequest(t *testing.T, ctx context.Context, serverURL string) *http.Request {
	url := settingsEndpoints.InternalBulkEndpoint + "?keys=" + settings.CompanyLogo
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+url, nil)
	require.NoError(t, err)
	return req
}
