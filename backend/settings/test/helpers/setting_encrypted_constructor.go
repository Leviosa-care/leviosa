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

func NewValidEncryptedSettingString(key, value string) *domain.SettingEncrypted[string] {
	now := time.Now()
	return &domain.SettingEncrypted[string]{
		Key:            key,
		Value:          value,
		ValueEncrypted: randomBytes(16),
		CreatedAt:      now,
		UpdatedAt:      now,
		DEK:            randomBytes(32),
		DEKEncrypted:   randomBytes(48),
		KeyVersion:     1,
	}
}

func NewValidEncryptedSettingInt(key string, value int) *domain.SettingEncrypted[int] {
	now := time.Now()
	return &domain.SettingEncrypted[int]{
		Key:            key,
		Value:          value,
		ValueEncrypted: randomBytes(16),
		CreatedAt:      now,
		UpdatedAt:      now,
		DEK:            randomBytes(32),
		DEKEncrypted:   randomBytes(48),
		KeyVersion:     1,
	}
}

// ----- RANDOM FACTORIES -----

func NewRandomEncryptedSettingString() *domain.SettingEncrypted[string] {
	now := time.Now()
	key := fmt.Sprintf("key_%d", mrand.Int())
	value := fmt.Sprintf("value_%d", mrand.Int())

	return &domain.SettingEncrypted[string]{
		Key:            key,
		Value:          value,
		ValueEncrypted: randomBytes(16),
		CreatedAt:      now,
		UpdatedAt:      now,
		DEK:            randomBytes(32),
		DEKEncrypted:   randomBytes(48),
		KeyVersion:     mrand.Intn(5) + 1, // 1..5
	}
}

func NewRandomEncryptedSettingInt() *domain.SettingEncrypted[int] {
	now := time.Now()
	key := fmt.Sprintf("key_%d", mrand.Int())
	value := mrand.Intn(1000)

	return &domain.SettingEncrypted[int]{
		Key:            key,
		Value:          value,
		ValueEncrypted: randomBytes(16),
		CreatedAt:      now,
		UpdatedAt:      now,
		DEK:            randomBytes(32),
		DEKEncrypted:   randomBytes(48),
		KeyVersion:     mrand.Intn(5) + 1,
	}
}
