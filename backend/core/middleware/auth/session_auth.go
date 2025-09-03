package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
	mw "github.com/Leviosa-care/core/middleware"
	"github.com/google/uuid"
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

// RequireAccessToken validates access token from cookies and makes session available in context
func (m *SessionAuthMiddleware) RequireAccessToken(next mw.Handler) mw.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract access token from cookies
		cookie, err := r.Cookie(AccessTokenCookieName)
		if err != nil {
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		accessToken := cookie.Value
		if accessToken == "" {
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// Find session by access token using two-step lookup
		sessionID, sessionData, err := m.sessionRepo.FindSessionByAccessToken(r.Context(), accessToken)
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

		// Check session state - access tokens work for both pending and active sessions
		if session.State != SessionActive && session.State != SessionPending {
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		parsedSessionID, err := uuid.Parse(sessionID)
		if err != nil {
			httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid session ID format"), http.StatusInternalServerError)
			return
		}

		// Create lightweight SessionInfo for context
		sessionInfo := &SessionInfo{
			ID:     parsedSessionID,
			UserID: session.UserID,
			Role:   session.Role,
			State:  session.State,
		}

		// Add session info to context
		ctx := context.WithValue(r.Context(), sessionContextKey{}, sessionInfo)
		r = r.WithContext(ctx)

		// Continue to next handler
		next(w, r)
	}
}

// RequireRefreshToken validates refresh token from cookies for token refresh operations only
func (m *SessionAuthMiddleware) RequireRefreshToken(next mw.Handler) mw.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// Security: Only allow refresh tokens on /auth/refresh endpoint
		if r.URL.Path != RefreshEndpoint {
			httpx.RespondWithError(w, errs.ErrForbidden, http.StatusForbidden)
			return
		}

		// Extract refresh token from cookies
		cookie, err := r.Cookie(RefreshTokenCookieName)
		if err != nil {
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		refreshToken := cookie.Value
		if refreshToken == "" {
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// Find session by refresh token using two-step lookup
		sessionID, sessionData, err := m.sessionRepo.FindSessionByRefreshToken(r.Context(), refreshToken)
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

		// Check session state - refresh tokens work for both pending and active sessions
		if session.State != SessionActive && session.State != SessionPending {
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		parsedSessionID, err := uuid.Parse(sessionID)
		if err != nil {
			httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid session ID format"), http.StatusInternalServerError)
			return
		}

		// Create lightweight SessionInfo for context
		sessionInfo := &SessionInfo{
			ID:     parsedSessionID,
			UserID: session.UserID,
			Role:   session.Role,
			State:  session.State,
		}

		// Add session info to context
		ctx := context.WithValue(r.Context(), sessionContextKey{}, sessionInfo)
		r = r.WithContext(ctx)

		// Continue to next handler
		next(w, r)
	}
}

// RequireRefreshToken validates refresh token from cookies for token refresh operations only
// func (m *SessionAuthMiddleware) RequireRefreshToken(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Security: Only allow refresh tokens on /auth/refresh endpoint
// 		if r.URL.Path != RefreshEndpoint {
// 			httpx.RespondWithError(w, errs.ErrForbidden, http.StatusForbidden)
// 			return
// 		}
//
// 		// Extract refresh token from cookies
// 		cookie, err := r.Cookie(RefreshTokenCookieName)
// 		if err != nil {
// 			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
// 			return
// 		}
//
// 		refreshToken := cookie.Value
// 		if refreshToken == "" {
// 			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
// 			return
// 		}
//
// 		// Find session by refresh token using two-step lookup
// 		sessionID, sessionData, err := m.sessionRepo.FindSessionByRefreshToken(r.Context(), refreshToken)
// 		if err != nil {
// 			if errors.Is(err, errs.ErrRepositoryNotFound) {
// 				httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
// 				return
// 			}
// 			httpx.RespondWithError(w, err, http.StatusInternalServerError)
// 			return
// 		}
//
// 		// Decode session
// 		session, err := DecodeSession(sessionData)
// 		if err != nil {
// 			httpx.RespondWithError(w, err, http.StatusInternalServerError)
// 			return
// 		}
//
// 		// Check session state - refresh tokens work for both pending and active sessions
// 		if session.State != SessionActive && session.State != SessionPending {
// 			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
// 			return
// 		}
//
// 		parsedSessionID, err := uuid.Parse(sessionID)
// 		if err != nil {
// 			httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid session ID format"), http.StatusInternalServerError)
// 			return
// 		}
//
// 		// Create lightweight SessionInfo for context
// 		sessionInfo := &SessionInfo{
// 			ID:     parsedSessionID,
// 			UserID: session.UserID,
// 			Role:   session.Role,
// 			State:  session.State,
// 		}
//
// 		// Add session info to context
// 		ctx := context.WithValue(r.Context(), sessionContextKey{}, sessionInfo)
// 		r = r.WithContext(ctx)
//
// 		// Continue to next handler
// 		next.ServeHTTP(w, r)
// 	})
// }
