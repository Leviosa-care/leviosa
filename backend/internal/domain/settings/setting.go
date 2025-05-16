package settings

import (
	"fmt"
	"time"
)

type Setting[T any] struct {
	ID        string    `json:"-"`
	Key       string    `json:"-"`
	Value     T         `json:"value"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type SettingEncrypted[T any] struct {
	ID             string    `json:"-"`
	Key            string    `json:"-"`
	Value          T         `json:"value" encx:"encrypt"`
	ValueEncrypted []byte    `json:"-"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
	DEK            []byte    `json:"-"`
	DEKEncrypted   []byte    `json:"-"`
	KeyVersion     int       `json:"-"`
}

func NewSetting[T any](key string, value T) *Setting[T] {
	return &Setting[T]{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewSettingEncrypted[T any](svc *service, key string, value T) (*SettingEncrypted[T], error) {
	dek, err := svc.crypto.GenerateDEK()
	if err != nil {
		return nil, fmt.Errorf("failed to generate DEK for OTP: %w", err)
	}
	return &SettingEncrypted[T]{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DEK:       dek,
	}, nil
}
