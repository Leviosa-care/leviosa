package ports

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
)

type UserService interface {
	CheckEmailAvailability(ctx context.Context, request *domain.CheckEmailAvailabilityRequest) (bool, error)
	CreatePendingUser(ctx context.Context, email string) (uuid.UUID, error)
	CompleteUser(ctx context.Context, userID uuid.UUID, request *domain.CompleteUserRequest) error
	GetPendingUsers(ctx context.Context) ([]*domain.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]*domain.UserResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error)
	GetUserByEmailHash(ctx context.Context, emailHash string) (*domain.UserResponse, error)
	ApproveUser(ctx context.Context, request *domain.ApproveUserRequest) error
}
