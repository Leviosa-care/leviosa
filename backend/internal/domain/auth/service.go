package auth

import "context"

type Service interface {
	Login(ctx context.Context, email, password string) (AuthResult, error)
	LoginWithOTP(ctx context.Context, phone, otp string) (AuthResult, error)
	SendOTP(ctx context.Context, phone string) error
	RefreshSession(ctx context.Context, refreshToken string) (AuthResult, error)
	Logout(ctx context.Context, sessionID string) error
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	UserID       string
	SessionID    string
}

type service struct {
}

func New() Service {
	return &service{}
}
