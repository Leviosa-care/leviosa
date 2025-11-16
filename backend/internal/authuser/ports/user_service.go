package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
)

type UserService interface {
	CheckEmailAvailability(ctx context.Context, request *domain.CheckEmailAvailabilityRequest) (bool, error)
	CreatePendingUser(ctx context.Context, email string) (uuid.UUID, error)
	CompleteUser(ctx context.Context, userID uuid.UUID, request *domain.CompleteUserRequest) error
	GetPendingUsers(ctx context.Context) ([]*domain.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]*domain.UserResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.UserResponse, error)
	VerifyUserPassword(ctx context.Context, userID uuid.UUID, password string) error
	ChangePassword(ctx context.Context, userID uuid.UUID, request *domain.ChangePasswordRequest) error
	ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	ApproveUser(ctx context.Context, request *domain.ApproveUserRequest) error
	UpdateUserRole(ctx context.Context, request *domain.UpdateUserRoleRequest) error
	UpdateUser(ctx context.Context, userID uuid.UUID, request *domain.UpdateUserRequest) (*domain.UserResponse, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	GetOrCreateOAuthUser(ctx context.Context, provider, userID, email, firstName, lastName string) (*domain.UserResponse, bool, error)
	GetUserByGoogleID(ctx context.Context, googleID string) (*domain.UserResponse, error)
	GetUserByAppleID(ctx context.Context, appleID string) (*domain.UserResponse, error)
	UpdateLastLoginTime(ctx context.Context, userID uuid.UUID) error
}
