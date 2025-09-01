package ports

import (
	"context"
	"io"

	"github.com/Leviosa-care/settings/internal/domain"
)

type SettingsService interface {
	// getters
	GetCompanyName(ctx context.Context) (*domain.GetCompanyNameResponse, error)
	GetCompanyEmail(ctx context.Context) (*domain.GetCompanyEmailResponse, error)
	GetCompanyTelephone(ctx context.Context) (*domain.GetCompanyTelephoneResponse, error)
	GetCompanyLegalAddress(ctx context.Context) (*domain.GetCompanyLegalAddressResponse, error)
	GetCompanyInstagram(ctx context.Context) (*domain.GetCompanyInstagramResponse, error)
	GetCompanyLogo(ctx context.Context) (*domain.GetCompanyLogoResponse, error)
	GetOTPDuration(ctx context.Context) (*domain.GetOTPDurationResponse, error)
	GetOTPLength(ctx context.Context) (*domain.GetOTPLengthResponse, error)
	GetOTPMaxAttempts(ctx context.Context) (*domain.GetOTPMaxAttemptsResponse, error)
	GetAccessTokenDuration(ctx context.Context) (*domain.GetAccessTokenDurationResponse, error)
	GetRefreshTokenDuration(ctx context.Context) (*domain.GetRefreshTokenDurationResponse, error)
	// setters
	SetCompanyName(ctx context.Context, request *domain.SetCompanyNameRequest) (*domain.SetCompanyNameResponse, error)
	SetCompanyEmail(ctx context.Context, request *domain.SetCompanyEmailRequest) (*domain.SetCompanyEmailResponse, error)
	SetCompanyTelephone(ctx context.Context, request *domain.SetCompanyTelephoneRequest) (*domain.SetCompanyTelephoneResponse, error)
	SetCompanyLegalAddress(ctx context.Context, request *domain.SetCompanyLegalAddressRequest) (*domain.SetCompanyLegalAddressResponse, error)
	SetCompanyInstagram(ctx context.Context, request *domain.SetCompanyInstagramRequest) (*domain.SetCompanyInstagramResponse, error)
	SetCompanyLogo(ctx context.Context, file io.Reader, fileSize int64, contentType string) (*domain.SetCompanyLogoResponse, error)
	SetOTPDuration(ctx context.Context, request *domain.SetOTPDurationRequest) (*domain.SetOTPDurationResponse, error)
	SetOTPLength(ctx context.Context, request *domain.SetOTPLengthRequest) (*domain.SetOTPLengthResponse, error)
	SetOTPMaxAttempts(ctx context.Context, request *domain.SetOTPMaxAttemptsRequest) (*domain.SetOTPMaxAttemptsResponse, error)
	SetAccessTokenDuration(ctx context.Context, request *domain.SetAccessTokenDurationRequest) (*domain.SetAccessTokenDurationResponse, error)
	SetRefreshTokenDuration(ctx context.Context, request *domain.SetRefreshTokenDurationRequest) (*domain.SetRefreshTokenDurationResponse, error)
	// publishers
	PublishSettingUpdate(ctx context.Context, key string, value any) error
}
