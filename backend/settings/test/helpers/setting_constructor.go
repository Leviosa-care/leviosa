package helpers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Leviosa-care/settings/internal/domain"
)

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
	key := fmt.Sprintf("key_%d", rand.Int())
	value := fmt.Sprintf("value_%d", rand.Int())

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
	key := fmt.Sprintf("key_%d", rand.Int())
	value := rand.Intn(1000) // adjust range if needed

	return &domain.Setting[int]{
		Key:       key,
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
