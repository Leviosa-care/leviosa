package settings

import (
	"context"
	"errors"
	"io"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/core/errs"
	"github.com/hengadev/errsx"
)

func (s *SettingsService) SetCompanyLogo(ctx context.Context, file io.Reader, fileSize int64, contentType string) (*domain.SetCompanyLogoResponse, error) {
	request := &domain.SetCompanyLogoRequest{
		ContentType: contentType,
		FileSize:    fileSize,
	}

	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// file validation
	var fileErrs errsx.Map
	if file == nil {
		fileErrs.Set("file", "image file is required")
	}
	if fileErrs != nil {
		return nil, errs.NewInvalidValueErr(fileErrs.Error())
	}

	imageKey, err := CreateLogoPrefix("logo", contentType)
	if err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	if _, err = s.media.UploadLogo(ctx, imageKey, file, fileSize, contentType); err != nil {
		return nil, errs.NewExternalServiceErr(err, "get company logo")
	}

	if err := s.PublishSettingUpdate(ctx, settings.CompanyLogo, imageKey); err != nil {
		return nil, err
	}

	return &domain.SetCompanyLogoResponse{Success: true}, nil
}

func (s *SettingsService) GetCompanyLogo(ctx context.Context) (*domain.GetCompanyLogoResponse, error) {
	logo, err := s.media.GetLogo(ctx)
	_ = logo
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewNotFoundErr(err, "company logo")
		}
		return nil, err
	}

	// Note: This returns raw bytes. You may want to return a URL instead
	// depending on your frontend requirements
	return &domain.GetCompanyLogoResponse{
		LogoURL:     "/settings/logo", // or construct actual URL
		ContentType: "image/jpeg",     // you may need to store this
	}, nil
}
