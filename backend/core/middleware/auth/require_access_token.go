package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Leviosa-care/core/auth/cookies"
	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
	mw "github.com/Leviosa-care/core/middleware"
	"github.com/google/uuid"
)

// RequireAccessToken validates access token from cookies and makes session available in context
func (m *SessionAuthMiddleware) RequireAccessToken(next mw.Handler) mw.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger, err := ctxutil.GetLoggerFromContext(ctx)
		if err != nil {
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		// Extract access token from cookies
		cookie, err := r.Cookie(cookies.AccessTokenCookieName)
		if err != nil {
			logger.WarnContext(ctx, "Auth middleware: Missing access token cookie",
				"operation", "require_access_token",
				"method", r.Method,
				"path", r.URL.Path,
				"error", err)
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		accessToken := cookie.Value
		if accessToken == "" {
			logger.WarnContext(ctx, "Auth middleware: Empty access token",
				"operation", "require_access_token",
				"method", r.Method,
				"path", r.URL.Path)
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		accessTokenHash := m.crypto.HashBasic(ctx, []byte(accessToken))

		// Find session by access token using two-step lookup
		sessionID, sessionData, err := m.sessionRepo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		if err != nil {
			if errors.Is(err, errs.ErrRepositoryNotFound) {
				logger.WarnContext(ctx, "Auth middleware: Session not found for access token",
					"operation", "require_access_token",
					"method", r.Method,
					"path", r.URL.Path)
				httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
				return
			}
			logger.ErrorContext(ctx, "Auth middleware: Failed to find session by access token",
				"operation", "require_access_token",
				"method", r.Method,
				"path", r.URL.Path,
				"error", err)
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		// Decode session
		sessionStruct, err := session.DecodeSession(sessionData)
		if err != nil {
			logger.ErrorContext(ctx, "Auth middleware: Failed to decode session",
				"operation", "require_access_token",
				"method", r.Method,
				"path", r.URL.Path,
				"session_id", sessionID,
				"error", err)
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		m.crypto.DecryptStruct(ctx, sessionStruct)

		logger.InfoContext(ctx, "Auth middleware: Session retrieved and decrypted",
			"operation", "require_access_token",
			"method", r.Method,
			"path", r.URL.Path,
			"session_id", sessionID,
			"user_id", sessionStruct.UserID,
			"session_state", sessionStruct.State,
			"user_role", sessionStruct.Role)

		// Check session expiration
		if time.Now().After(sessionStruct.ExpiresAt) {
			logger.WarnContext(ctx, "Auth middleware: Session has expired",
				"operation", "require_access_token",
				"method", r.Method,
				"path", r.URL.Path,
				"session_id", sessionID,
				"user_id", sessionStruct.UserID,
				"expires_at", sessionStruct.ExpiresAt,
				"current_time", time.Now())
			httpx.RespondWithError(w, errs.ErrExpiredToken, http.StatusUnauthorized)
			return
		}

		// Check session state - access tokens work for both pending and active sessions
		if sessionStruct.State != session.SessionActive && sessionStruct.State != session.SessionPending {
			logger.WarnContext(ctx, "Auth middleware: Invalid session state",
				"operation", "require_access_token",
				"method", r.Method,
				"path", r.URL.Path,
				"session_id", sessionID,
				"user_id", sessionStruct.UserID,
				"session_state", sessionStruct.State,
				"expected_states", []string{string(session.SessionActive), string(session.SessionPending)})
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		parsedSessionID, err := uuid.Parse(sessionID)
		if err != nil {
			logger.ErrorContext(ctx, "Auth middleware: Invalid session ID format",
				"operation", "require_access_token",
				"method", r.Method,
				"path", r.URL.Path,
				"session_id", sessionID,
				"error", err)
			httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid session ID format"), http.StatusInternalServerError)
			return
		}

		// Create lightweight SessionInfo for context
		sessionInfo := &session.SessionInfo{
			ID:     parsedSessionID,
			UserID: sessionStruct.UserID,
			Role:   sessionStruct.Role,
			State:  sessionStruct.State,
		}

		// Add session info to context
		ctx = context.WithValue(ctx, session.GetSessionContextKey(), sessionInfo)
		r = r.WithContext(ctx)

		logger.InfoContext(ctx, "Auth middleware: Access token validation successful",
			"operation", "require_access_token",
			"method", r.Method,
			"path", r.URL.Path,
			"session_id", sessionID,
			"user_id", sessionStruct.UserID,
			"user_role", sessionStruct.Role)

		// Continue to next handler
		next(w, r)
	}
}
