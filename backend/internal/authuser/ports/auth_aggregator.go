package ports

import (
	"context"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
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
	CompletePartner(ctx context.Context, sessionInfo *session.SessionInfo, request *domain.CompletePartnerRequest) error
	RefreshSession(ctx context.Context, sessionID uuid.UUID) (*domain.RefreshSessionResponse, error)
	DeleteUserByAdmin(ctx context.Context, userID uuid.UUID) error
	DeleteOwnAccount(ctx context.Context, sessionInfo *session.SessionInfo) error
	SignIn(ctx context.Context, request *domain.SignInRequest) (*domain.CreateSessionResponse, error)
	SignOut(ctx context.Context, sessionInfo *session.SessionInfo) error
	RequestPasswordReset(ctx context.Context, request *domain.RequestPasswordResetRequest) error
	ValidatePasswordResetOTP(ctx context.Context, request *domain.ValidatePasswordResetOTPRequest) (*domain.ValidatePasswordResetOTPResponse, error)
	ConfirmPasswordReset(ctx context.Context, request *domain.ConfirmPasswordResetRequest) error
	OAuthStart(ctx context.Context, request *domain.OAuthStartRequest) (*domain.OAuthStartResponse, error)
	OAuthCallback(ctx context.Context, w http.ResponseWriter, r *http.Request, provider string) (*domain.OAuthCallbackResponse, error)
	// TODO: not sure about that one
	// ValidateSession(ctx context.Context) error
}
