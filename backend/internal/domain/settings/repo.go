package settings

import (
	"context"
)

type Reader interface {
	GetString(ctx context.Context, key string) (*Setting[string], error)
	GetInt(ctx context.Context, key string) (*Setting[int], error)
}

type encryptedReadWriter interface {
	GetPhone(ctx context.Context) (*SettingEncrypted[string], error)
	SetPhone(ctx context.Context, setting *SettingEncrypted[string]) error
}

type writer interface {
	SetString(ctx context.Context, setting *Setting[string]) error
	SetInt(ctx context.Context, setting *Setting[int]) error
}

type MediaReader interface {
	GetLogo(ctx context.Context) ([]byte, error)
}

type mediaWriter interface {
	SetLogo(ctx context.Context, logo []byte) error
}

type mediaReadWriter interface {
	MediaReader
	mediaWriter
}

type ReadWriter interface {
	Reader
	writer
	encryptedReadWriter
}
