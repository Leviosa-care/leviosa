package settings

import (
	"time"

	"github.com/Leviosa-care/settings/internal/domain"
)

func NewSetting[T any](key string, value T) *domain.Setting[T] {
	return &domain.Setting[T]{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
