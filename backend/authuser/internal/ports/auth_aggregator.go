package ports

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/auth/session"
	"github.com/google/uuid"
)

// signup flow for the frontend pages
// 1. auth/email (email)
// 2. auth/otp (verify the OTP sent)
// 3. auth/general (firstname, lastname, age, gender etc...)
// 4. auth/address (self explanatory)
// 5. auth/password (self explanatory)
// 6. auth/pending (redirect user to a page that allows them to visit the website as a pending user because admin need to approve user)

type AuthAggregatorService interface {
	CheckEmailSendOTP(ctx context.Context, request *domain.CheckEmailAvailabilityRequest) error
	ValidateOTPCreatePendingUser(ctx context.Context, request *domain.ValidateOTPRequest) (*domain.CreateSessionResponse, error)
	CompleteUser(ctx context.Context, sessionInfo *session.SessionInfo, request *domain.CompleteUserRequest) error
	RefreshSession(ctx context.Context, sessionID uuid.UUID) (*domain.RefreshSessionResponse, error)
	DeleteUserByAdmin(ctx context.Context, userID uuid.UUID) error
	DeleteOwnAccount(ctx context.Context, sessionInfo *session.SessionInfo) error
	SignIn(ctx context.Context, request *domain.SignInRequest) (*domain.CreateSessionResponse, error)
	SignOut(ctx context.Context, sessionInfo *session.SessionInfo) error
	// TODO: not sure about that one
	// ValidateSession(ctx context.Context) error
}
