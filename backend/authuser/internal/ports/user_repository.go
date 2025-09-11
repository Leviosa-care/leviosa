package ports

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	ExistsByEmailHash(ctx context.Context, emailHash string) (bool, error)
	GetUserByEmailHash(ctx context.Context, emailHash string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	GetPendingUsers(ctx context.Context) ([]*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	GetUserByGoogleID(ctx context.Context, googleID string) (*domain.User, error)
	GetUserByAppleID(ctx context.Context, appleID string) (*domain.User, error)
	ExistsByGoogleID(ctx context.Context, googleID string) (bool, error)
	ExistsByAppleID(ctx context.Context, appleID string) (bool, error)
}
