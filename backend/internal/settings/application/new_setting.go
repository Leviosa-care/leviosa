package settings

import (
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"
)

func NewSetting[T any](key string, value T) *domain.Setting[T] {
	return &domain.Setting[T]{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
