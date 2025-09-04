package ports

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByEmailHash(ctx context.Context, emailHash string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
}
