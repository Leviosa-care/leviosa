package application

import (
	"context"
	"io"
	"testing"

	settingsDomain "github.com/Leviosa-care/leviosa/backend/internal/settings/domain"
	settingsPorts "github.com/Leviosa-care/leviosa/backend/internal/settings/ports"
)

// mockSettingsService implements settingsPorts.SettingsService for testing.
type mockSettingsService struct {
	logoURL string
	logoErr error
}

var _ settingsPorts.SettingsService = (*mockSettingsService)(nil)

func (m *mockSettingsService) GetCompanyLogo(ctx context.Context) (*settingsDomain.GetCompanyLogoResponse, error) {
	if m.logoErr != nil {
		return nil, m.logoErr
	}
	return &settingsDomain.GetCompanyLogoResponse{LogoURL: m.logoURL}, nil
}

func (m *mockSettingsService) GetCompanyName(ctx context.Context) (*settingsDomain.GetCompanyNameResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) GetCompanyEmail(ctx context.Context) (*settingsDomain.GetCompanyEmailResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) GetCompanyTelephone(ctx context.Context) (*settingsDomain.GetCompanyTelephoneResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) GetCompanyLegalAddress(ctx context.Context) (*settingsDomain.GetCompanyLegalAddressResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) GetCompanyInstagram(ctx context.Context) (*settingsDomain.GetCompanyInstagramResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) GetOTPDuration(ctx context.Context) (*settingsDomain.GetOTPDurationResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) GetOTPLength(ctx context.Context) (*settingsDomain.GetOTPLengthResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) GetOTPMaxAttempts(ctx context.Context) (*settingsDomain.GetOTPMaxAttemptsResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) GetAccessTokenDuration(ctx context.Context) (*settingsDomain.GetAccessTokenDurationResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) GetRefreshTokenDuration(ctx context.Context) (*settingsDomain.GetRefreshTokenDurationResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) SetCompanyName(ctx context.Context, req *settingsDomain.SetCompanyNameRequest) (*settingsDomain.SetCompanyNameResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) SetCompanyEmail(ctx context.Context, req *settingsDomain.SetCompanyEmailRequest) (*settingsDomain.SetCompanyEmailResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) SetCompanyTelephone(ctx context.Context, req *settingsDomain.SetCompanyTelephoneRequest) (*settingsDomain.SetCompanyTelephoneResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) SetCompanyLegalAddress(ctx context.Context, req *settingsDomain.SetCompanyLegalAddressRequest) (*settingsDomain.SetCompanyLegalAddressResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) SetCompanyInstagram(ctx context.Context, req *settingsDomain.SetCompanyInstagramRequest) (*settingsDomain.SetCompanyInstagramResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) SetCompanyLogo(ctx context.Context, file io.Reader, fileSize int64, contentType string) (*settingsDomain.SetCompanyLogoResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) SetOTPDuration(ctx context.Context, req *settingsDomain.SetOTPDurationRequest) (*settingsDomain.SetOTPDurationResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) SetOTPLength(ctx context.Context, req *settingsDomain.SetOTPLengthRequest) (*settingsDomain.SetOTPLengthResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) SetOTPMaxAttempts(ctx context.Context, req *settingsDomain.SetOTPMaxAttemptsRequest) (*settingsDomain.SetOTPMaxAttemptsResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) SetAccessTokenDuration(ctx context.Context, req *settingsDomain.SetAccessTokenDurationRequest) (*settingsDomain.SetAccessTokenDurationResponse, error) {
	return nil, nil
}
func (m *mockSettingsService) SetRefreshTokenDuration(ctx context.Context, req *settingsDomain.SetRefreshTokenDurationRequest) (*settingsDomain.SetRefreshTokenDurationResponse, error) {
	return nil, nil
}

func TestSettingsProvider_GetCompanyLogo_returnsURLFromService(t *testing.T) {
	const wantURL = "https://s3.example.com/logo.png"

	svc := &mockSettingsService{logoURL: wantURL}
	provider := NewSettingsProvider(svc)

	got, err := provider.GetCompanyLogo(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != wantURL {
		t.Errorf("got logo URL %q, want %q", got, wantURL)
	}
}

func TestSettingsProvider_GetCompanyLogo_cachesResult(t *testing.T) {
	const wantURL = "https://s3.example.com/logo.png"

	svc := &mockSettingsService{logoURL: wantURL}
	provider := NewSettingsProvider(svc)

	// First call fetches from service.
	if _, err := provider.GetCompanyLogo(context.Background()); err != nil {
		t.Fatalf("first call: %v", err)
	}

	// Swap the mock URL to confirm the second call returns the cached value.
	svc.logoURL = "https://s3.example.com/different.png"

	got, err := provider.GetCompanyLogo(context.Background())
	if err != nil {
		t.Fatalf("second call: %v", err)
	}
	if got != wantURL {
		t.Errorf("got logo URL %q after cache, want cached %q", got, wantURL)
	}
}
