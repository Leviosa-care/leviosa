package settings

import (
	"fmt"
	"time"

	"github.com/Leviosa-care/settings/internal/domain"
)

func NewSettingEncrypted[T any](svc *SettingsService, key string, value T) (*domain.SettingEncrypted[T], error) {
	dek, err := svc.crypto.GenerateDEK()
	if err != nil {
		return nil, fmt.Errorf("failed to generate DEK for OTP: %w", err)
	}
	return &domain.SettingEncrypted[T]{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DEK:       dek,
	}, nil
}
