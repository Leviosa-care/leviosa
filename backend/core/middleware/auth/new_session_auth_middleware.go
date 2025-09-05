package auth

// SessionAuthMiddleware implements AuthMiddleware using session repository
type SessionAuthMiddleware struct {
	sessionRepo SessionRepository
}

// NewSessionAuthMiddleware creates a new session-based auth middleware
func NewSessionAuthMiddleware(sessionRepo SessionRepository) AuthMiddleware {
	return &SessionAuthMiddleware{
		sessionRepo: sessionRepo,
	}
}