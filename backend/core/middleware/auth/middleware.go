package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
)

// AuthMiddleware provides session-based authentication and role-based authorization
type AuthMiddleware interface {
	// RequireSession validates session and makes session available in context
	RequireSession(next http.Handler) http.Handler
	
	// RequireMinimumRole validates session and ensures user has at least the specified role
	RequireMinimumRole(minRole identity.Role) func(http.Handler) http.Handler
	
	// RequireAnyRole validates session and ensures user has one of the specified roles
	RequireAnyRole(roles ...identity.Role) func(http.Handler) http.Handler
	
	// RequireAdmin validates session and ensures user has admin role
	RequireAdmin(next http.Handler) http.Handler
}

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

// sessionContextKey is used to store session data in request context
type sessionContextKey struct{}

// SessionFromContext extracts session from request context
func SessionFromContext(ctx context.Context) (*Session, bool) {
	session, ok := ctx.Value(sessionContextKey{}).(*Session)
	return session, ok
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

// RequireMinimumRole validates session and ensures user has at least the specified role
func (m *SessionAuthMiddleware) RequireMinimumRole(minRole identity.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return m.RequireSession(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, ok := SessionFromContext(r.Context())
			if !ok {
				httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
				return
			}

			if !session.Role.IsAtLeast(minRole) {
				httpx.RespondWithError(w, errs.ErrPermissionDenied, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}

// RequireAnyRole validates session and ensures user has one of the specified roles
func (m *SessionAuthMiddleware) RequireAnyRole(roles ...identity.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return m.RequireSession(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, ok := SessionFromContext(r.Context())
			if !ok {
				httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
				return
			}

			hasAnyRole := false
			for _, role := range roles {
				if session.Role == role {
					hasAnyRole = true
					break
				}
			}

			if !hasAnyRole {
				httpx.RespondWithError(w, errs.ErrPermissionDenied, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}

// RequireAdmin validates session and ensures user has admin role
func (m *SessionAuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return m.RequireMinimumRole(identity.RoleAdmin)(next)
}