package domain

import (
	"context"
	"net/url"
	"strings"

	"github.com/Leviosa-care/core/validation"
	"github.com/hengadev/errsx"
)

type SetCompanyNameRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type SetCompanyNameResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type GetCompanyNameResponse struct {
	Name string `json:"name"`
}

type SetCompanyEmailRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type SetCompanyEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type GetCompanyEmailResponse struct {
	Email string `json:"email"`
}

type SetCompanyTelephoneRequest struct {
	Telephone string `json:"telephone" validate:"required,min=10,max=20"`
}

type SetCompanyTelephoneResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type GetCompanyTelephoneResponse struct {
	Telephone string `json:"telephone"`
}

type SetCompanyLegalAddressRequest struct {
	Address string `json:"address" validate:"required,min=1,max=500"`
}

type SetCompanyLegalAddressResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type GetCompanyLegalAddressResponse struct {
	Address string `json:"address"`
}

type SetCompanyInstagramRequest struct {
	Instagram string `json:"instagram" validate:"required,url,max=255"`
}

type SetCompanyInstagramResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type GetCompanyInstagramResponse struct {
	Instagram string `json:"instagram"`
}

type SetCompanyLogoRequest struct {
	ContentType string `json:"content_type" validate:"required,oneof=image/jpeg image/png image/gif"`
	FileSize    int64  `json:"file_size" validate:"required,min=1,max=5242880"`
}

type SetCompanyLogoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type GetCompanyLogoResponse struct {
	LogoURL     string `json:"logo_url"`
	ContentType string `json:"content_type"`
}

type SetOTPDurationRequest struct {
	Duration int `json:"duration" validate:"required,min=60,max=3600"`
}

type SetOTPDurationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type GetOTPDurationResponse struct {
	Duration int `json:"duration"`
}

type SetOTPLengthRequest struct {
	Length int `json:"length" validate:"required,min=4,max=10"`
}

type SetOTPLengthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type GetOTPLengthResponse struct {
	Length int `json:"length"`
}

type SetOTPMaxAttemptsRequest struct {
	MaxAttempts int `json:"max_attempts" validate:"required,min=1,max=10"`
}

type SetOTPMaxAttemptsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type GetOTPMaxAttemptsResponse struct {
	MaxAttempts int `json:"max_attempts"`
}

type SetAccessTokenDurationRequest struct {
	Duration int `json:"duration" validate:"required,min=1,max=240"`
}

type SetAccessTokenDurationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type GetAccessTokenDurationResponse struct {
	Duration int `json:"duration"`
}

type SetRefreshTokenDurationRequest struct {
	Duration int `json:"duration" validate:"required,min=1,max=720"`
}

type SetRefreshTokenDurationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type GetRefreshTokenDurationResponse struct {
	Duration int `json:"duration"`
}

type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
	Code    string            `json:"code,omitempty"`
}

func (d SetCompanyNameRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if d.Name == "" {
		errs.Set("name_required", "company name is required")
	}
	if len(strings.TrimSpace(d.Name)) == 0 {
		errs.Set("name_empty", "company name cannot be empty or whitespace only")
	}
	if len(d.Name) > 255 {
		errs.Set("name_length", "company name cannot exceed 255 characters")
	}
	return errs.AsError()
}

func (d SetCompanyEmailRequest) Valid(ctx context.Context) error {
	return validation.ValidateEmail(d.Email)
}

func (d SetCompanyTelephoneRequest) Valid(ctx context.Context) error {
	return validation.ValidatePhone(d.Telephone)
}

func (d SetCompanyLegalAddressRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if d.Address == "" {
		errs.Set("address_required", "legal address is required")
	}
	if len(strings.TrimSpace(d.Address)) == 0 {
		errs.Set("address_empty", "legal address cannot be empty or whitespace only")
	}
	if len(d.Address) > 500 {
		errs.Set("address_length", "legal address cannot exceed 500 characters")
	}
	return errs.AsError()
}

func (d SetCompanyInstagramRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if d.Instagram == "" {
		errs.Set("instagram_required", "instagram link is required")
	}
	if len(d.Instagram) > 255 {
		errs.Set("instagram_length", "instagram link cannot exceed 255 characters")
	}

	// Stricter URL validation
	u, err := url.ParseRequestURI(d.Instagram)
	if err != nil {
		errs.Set("instagram_format", "invalid URL format for instagram link")
	} else {
		// Check scheme - only http/https allowed
		if u.Scheme != "http" && u.Scheme != "https" {
			errs.Set("instagram_format", "instagram link must use http or https")
		}
		// Check host exists (not empty)
		if u.Host == "" {
			errs.Set("instagram_format", "instagram link must have a valid domain")
		}
	}

	return errs.AsError()
}

func (d SetCompanyLogoRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if d.ContentType == "" {
		errs.Set("content_type_required", "content type is required")
	}
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
	validType := false
	for _, allowedType := range allowedTypes {
		if d.ContentType == allowedType {
			validType = true
			break
		}
	}
	if !validType {
		errs.Set("content_type_invalid", "content type must be image/jpeg, image/png, or image/gif")
	}
	if d.FileSize <= 0 {
		errs.Set("file_size_invalid", "file size must be greater than 0")
	}
	if d.FileSize > 5242880 {
		errs.Set("file_size_too_large", "file size cannot exceed 5MB")
	}
	return errs.AsError()
}

func (d SetOTPDurationRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if d.Duration < 60 {
		errs.Set("duration_min", "OTP duration must be at least 60 seconds")
	}
	if d.Duration > 3600 {
		errs.Set("duration_max", "OTP duration cannot exceed 3600 seconds (1 hour)")
	}
	return errs.AsError()
}

func (d SetOTPLengthRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if d.Length < 4 {
		errs.Set("length_min", "OTP length must be at least 4 digits")
	}
	if d.Length > 10 {
		errs.Set("length_max", "OTP length cannot exceed 10 digits")
	}
	return errs.AsError()
}

func (d SetOTPMaxAttemptsRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if d.MaxAttempts < 1 {
		errs.Set("max_attempts_min", "OTP max attempts must be at least 1")
	}
	if d.MaxAttempts > 10 {
		errs.Set("max_attempts_max", "OTP max attempts cannot exceed 10")
	}
	return errs.AsError()
}

func (d SetAccessTokenDurationRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if d.Duration < 1 {
		errs.Set("duration_min", "access token duration must be at least 1 minute")
	}
	if d.Duration > 240 {
		errs.Set("duration_max", "access token duration cannot exceed 240 minutes (4 hours)")
	}
	return errs.AsError()
}

func (d SetRefreshTokenDurationRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if d.Duration < 1 {
		errs.Set("duration_min", "refresh token duration must be at least 1 hour")
	}
	if d.Duration > 720 {
		errs.Set("duration_max", "refresh token duration cannot exceed 720 hours (30 days)")
	}
	return errs.AsError()
}
