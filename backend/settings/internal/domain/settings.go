package domain

import (
	"time"

	"github.com/google/uuid"
)

type Setting[T any] struct {
	ID        uuid.UUID `json:"-"`
	Key       string    `json:"-"`
	Value     T         `json:"value"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type SettingEncrypted[T any] struct {
	ID             uuid.UUID `json:"-"`
	Key            string    `json:"-"`
	Value          T         `json:"value" encx:"encrypt"`
	ValueEncrypted []byte    `json:"-"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
	DEK            []byte    `json:"-"`
	DEKEncrypted   []byte    `json:"-"`
	KeyVersion     int       `json:"-"`
}
