package ports

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	ExistsByEmailHash(ctx context.Context, emailHash string) (bool, error)
	GetUserByEmailHash(ctx context.Context, emailHash string) (*domain.UserEncx, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.UserEncx, error)
	GetPendingUsers(ctx context.Context) ([]*domain.UserEncx, error)
	GetAllUsers(ctx context.Context) ([]*domain.UserEncx, error)
	CreateUser(ctx context.Context, user *domain.UserEncx) error
	UpdateUser(ctx context.Context, user *domain.UserEncx) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	GetUserByGoogleID(ctx context.Context, googleID string) (*domain.UserEncx, error)
	GetUserByAppleID(ctx context.Context, appleID string) (*domain.UserEncx, error)
	ExistsByGoogleID(ctx context.Context, googleID string) (bool, error)
	ExistsByAppleID(ctx context.Context, appleID string) (bool, error)
}
