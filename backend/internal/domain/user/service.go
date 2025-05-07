package userService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/user/models"

	"github.com/hengadev/encx"
)

type Service interface {
	CheckUser(ctx context.Context, email string) error
	CreateOAuthPendingUser(ctx context.Context, user *models.User, provider models.ProviderType) error
	CreateOAuthAccount(ctx context.Context, userCandidate *models.OAuthUser) (*models.User, error)
	CreatePendingUser(ctx context.Context, email string) error
	CreateUnverifiedUser(ctx context.Context, userSignUp *models.UserSignUp) (string, error)
	CreateUser(ctx context.Context, userResponse *models.UserPendingResponse) (*models.User, error)
	FindUserByID(ctx context.Context, userID string) (*models.User, error)
	GetAllPendingUsers(ctx context.Context) ([]*models.UserPending, error)
	GetUserSessionData(ctx context.Context, email string) (string, models.Role, error)
	DeleteUser(ctx context.Context, userID string) error
	UpdateAccount(ctx context.Context, user *models.User) error
	ValidateCredentials(ctx context.Context, user *models.UserSignIn) error
	CheckOAuthUser(ctx context.Context, email string, provider models.ProviderType) error
}

type service struct {
	repo   ReadWriter
	crypto encx.CryptoService
}

func New(repo ReadWriter, crypto encx.CryptoService) Service {
	return &service{
		repo,
		crypto,
	}
}
