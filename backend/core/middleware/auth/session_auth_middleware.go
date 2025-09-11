package auth

import (
	"github.com/Leviosa-care/core/auth/session"
	"github.com/hashicorp/vault/api"
	"github.com/hengadev/encx"
)

// SessionAuthMiddleware implements AuthMiddleware using session repository
type SessionAuthMiddleware struct {
	sessionRepo session.SessionRepository
	crypto      encx.CryptoService
	vaultClient *api.Client
}

// NewSessionAuthMiddleware creates a new session-based auth middleware
func NewSessionAuthMiddleware(sessionRepo session.SessionRepository, crypto encx.CryptoService, vaultClient *api.Client) AuthMiddleware {
	return &SessionAuthMiddleware{
		sessionRepo: sessionRepo,
		crypto:      crypto,
		vaultClient: vaultClient,
	}
}
