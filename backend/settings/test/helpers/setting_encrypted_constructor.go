package helpers

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/Leviosa-care/settings/internal/domain"
)

// helper: generate n random bytes
func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b) // crypto/rand, fine for test data
	return b
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
