package ports

import (
	"context"
	"io"

	"github.com/Leviosa-care/settings/internal/domain"
)

type SettingsRepository interface {
	// non encrypted
	GetString(ctx context.Context, key string) (*domain.Setting[string], error)
	GetInt(ctx context.Context, key string) (*domain.Setting[int], error)
	SetString(ctx context.Context, setting *domain.Setting[string]) error
	SetInt(ctx context.Context, setting *domain.Setting[int]) error
	// encrypted
	GetEncryptedSetting(ctx context.Context, key string) (*domain.SettingEncrypted[string], error)
	SetEncryptedSetting(ctx context.Context, setting *domain.SettingEncrypted[string]) error
}

type SettingsMedia interface {
	GetLogo(ctx context.Context) ([]byte, error)
	UploadLogo(ctx context.Context, key string, file io.Reader, size int64, contentType string) (string, error)
	// DeleteLogo(ctx context.Context)
}
