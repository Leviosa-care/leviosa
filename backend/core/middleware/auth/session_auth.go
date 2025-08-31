package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
)

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

// RequireSession validates session token and makes session available in context
func (m *SessionAuthMiddleware) RequireSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusUnauthorized)
			return
		}

		token := parts[1]
		if token == "" {
			httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusUnauthorized)
			return
		}

		// Find session by token hash
		sessionData, err := m.sessionRepo.FindSessionByTokenHash(r.Context(), token)
		if err != nil {
			if errors.Is(err, errs.ErrRepositoryNotFound) {
				httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
				return
			}
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		// Decode session
		session, err := DecodeSession(sessionData)
		if err != nil {
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		// Check session state
		if session.State != SessionActive {
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// Add session to context
		ctx := context.WithValue(r.Context(), sessionContextKey{}, session)
		r = r.WithContext(ctx)

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

