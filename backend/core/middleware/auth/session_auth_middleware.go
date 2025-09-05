package auth

import (
	"github.com/hengadev/encx"
)

// SessionAuthMiddleware implements AuthMiddleware using session repository
type SessionAuthMiddleware struct {
	sessionRepo SessionRepository
	crypto      encx.CryptoService
}

// NewSessionAuthMiddleware creates a new session-based auth middleware
func NewSessionAuthMiddleware(sessionRepo SessionRepository, crypto encx.CryptoService) AuthMiddleware {
	return &SessionAuthMiddleware{
		sessionRepo: sessionRepo,
		crypto:      crypto,
	}
}
