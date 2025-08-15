package auth

import "context"

// These are the things that I need to implement for the package to work
// otp.Service
type OTPService interface {
	GenerateOTP(ctx context.Context, phone string) (string, error)
	ValidateOTP(ctx context.Context, phone, otp string) (bool, error)
}

// session.Service
type SessionService interface {
	CreateSession(ctx context.Context, userID string) (Session, error)
	RefreshSession(ctx context.Context, refreshToken string) (Session, error)
	InvalidateSession(ctx context.Context, sessionID string) error
}

type Session struct {
	ID        string
	UserID    string
	ExpiresAt int64
}
