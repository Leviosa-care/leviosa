package settings

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) GetCompanyTelephone(ctx context.Context) (string, error) {
	setting, err := s.repo.GetPhone(ctx)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrNotFound):
			return "", domain.NewNotFoundErr(err)
		case errors.Is(err, rp.ErrContext):
			return "", err
		case errors.Is(err, rp.ErrDatabase):
			return "", domain.NewQueryFailedErr(err)
		}
	}
	if err := s.crypto.DecryptStruct(ctx, setting); err != nil {
		return "", domain.NewNotDecryptedErr("settings phone number", err)
	}
	return setting.Value, nil
}

func (s *service) SetCompanyTelephone(ctx context.Context, phone string) error {
	if err := models.ValidateTelephone(phone); err != nil {
		return domain.NewInvalidValueErr(err.Error())
	}
	setting, err := NewSettingEncrypted(s, CompanyPhoneKey, phone)
	if err != nil {
		return domain.NewNotCreatedErr(err)
	}
	if err := s.crypto.ProcessStruct(ctx, setting); err != nil {
		return domain.NewNotEncryptedErr("settings phone number", err)
	}
	return s.repo.SetPhone(ctx, setting)
}
