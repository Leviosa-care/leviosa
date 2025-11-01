package helpers

import (
	crand "crypto/rand"
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"
)

// helper: generate n random bytes
func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = crand.Read(b) // crypto/rand, fine for test data
	return b
}

func NewValidPlainSettingString(key, value string) *domain.Setting[string] {
	now := time.Now()
	return &domain.Setting[string]{
		Key:       key,
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewValidPlainSettingInt(key string, value int) *domain.Setting[int] {
	now := time.Now()
	return &domain.Setting[int]{
		Key:       key,
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Random string-valued setting
func NewRandomPlainSettingString() *domain.Setting[string] {
	now := time.Now()
	key := fmt.Sprintf("key_%d", mrand.Int())
	value := fmt.Sprintf("value_%d", mrand.Int())

	return &domain.Setting[string]{
		Key:       key,
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Random int-valued setting
func NewRandomPlainSettingInt() *domain.Setting[int] {
	now := time.Now()
	key := fmt.Sprintf("key_%d", mrand.Int())
	value := mrand.Intn(1000) // adjust range if needed

	return &domain.Setting[int]{
		Key:       key,
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ----- FACTORIES WITH PROVIDED VALUES -----

func NewValidEncryptedSettingString(key, value string) *domain.SettingEncrypted {
	now := time.Now()
	return &domain.SettingEncrypted{
		Key:       key,
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewValidEncryptedSettingInt(key string, value int) *domain.SettingEncrypted {
	now := time.Now()
	return &domain.SettingEncrypted{
		Key:       key,
		Value:     fmt.Sprintf("%d", value), // Convert int to string since SettingEncrypted.Value is string
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ----- RANDOM FACTORIES -----

func NewRandomEncryptedSettingString() *domain.SettingEncrypted {
	now := time.Now()
	key := fmt.Sprintf("key_%d", mrand.Int())
	value := fmt.Sprintf("value_%d", mrand.Int())

	return &domain.SettingEncrypted{
		Key:       key,
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewRandomEncryptedSettingInt() *domain.SettingEncrypted {
	now := time.Now()
	key := fmt.Sprintf("key_%d", mrand.Int())
	value := mrand.Intn(1000)

	return &domain.SettingEncrypted{
		Key:       key,
		Value:     fmt.Sprintf("%d", value), // Convert int to string since SettingEncrypted.Value is string
		CreatedAt: now,
		UpdatedAt: now,
	}
}
