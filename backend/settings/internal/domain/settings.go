package domain

import (
	"time"
)

//go:generate encx-gen generate .

// Setting is a plain key-value setting without encryption
type Setting[T any] struct {
	ID        string    `json:"-"`
	Key       string    `json:"-"`
	Value     T         `json:"value"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// SettingEncrypted is the clean source struct for encrypted string settings
// ENCX will generate SettingEncryptedEncx with encrypted fields
type SettingEncrypted struct {
	ID        string    `json:"-"`
	Key       string    `json:"-"`
	Value     string    `json:"value" encx:"encrypt"`  // Clean - no companion fields!
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
