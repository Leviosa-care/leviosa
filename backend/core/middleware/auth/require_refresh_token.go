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

// RequireRefreshToken validates refresh token from cookies for token refresh operations only
func (m *SessionAuthMiddleware) RequireRefreshToken(next mw.Handler) mw.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger, err := ctxutil.GetLoggerFromContext(ctx)
		if err != nil {
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		// Security: Only allow refresh tokens on /auth/refresh endpoint
		if r.URL.Path != cookies.RefreshEndpoint {
			logger.WarnContext(ctx, "Auth middleware: Refresh token attempted on wrong endpoint",
				"operation", "require_refresh_token",
				"method", r.Method,
				"path", r.URL.Path,
				"expected_path", cookies.RefreshEndpoint)
			httpx.RespondWithError(w, errs.ErrForbidden, http.StatusForbidden)
			return
		}

		// Extract refresh token from cookies
		cookie, err := r.Cookie(cookies.RefreshTokenCookieName)
		if err != nil {
			logger.WarnContext(ctx, "Auth middleware: Missing refresh token cookie",
				"operation", "require_refresh_token",
				"method", r.Method,
				"path", r.URL.Path,
				"error", err)
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		refreshToken := cookie.Value
		if refreshToken == "" {
			logger.WarnContext(ctx, "Auth middleware: Empty refresh token",
				"operation", "require_refresh_token",
				"method", r.Method,
				"path", r.URL.Path)
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		refreshTokenHash := m.crypto.HashBasic(ctx, []byte(refreshToken))

		// Find session by refresh token using two-step lookup
		sessionID, sessionData, err := m.sessionRepo.FindSessionByRefreshTokenHash(ctx, refreshTokenHash)
		if err != nil {
			if errors.Is(err, errs.ErrRepositoryNotFound) {
				logger.WarnContext(ctx, "Auth middleware: Session not found for refresh token",
					"operation", "require_refresh_token",
					"method", r.Method,
					"path", r.URL.Path)
				httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
				return
			}
			logger.ErrorContext(ctx, "Auth middleware: Failed to find session by refresh token",
				"operation", "require_refresh_token",
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
				"operation", "require_refresh_token",
				"method", r.Method,
				"path", r.URL.Path,
				"session_id", sessionID,
				"error", err)
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		err = m.crypto.DecryptStruct(ctx, sessionStruct)
		if err != nil {
			logger.ErrorContext(ctx, "Auth middleware: Failed to decrypt session",
				"operation", "require_refresh_token",
				"method", r.Method,
				"path", r.URL.Path,
				"session_id", sessionID,
				"error", err)
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		logger.InfoContext(ctx, "Auth middleware: Session retrieved and decrypted for refresh",
			"operation", "require_refresh_token",
			"method", r.Method,
			"path", r.URL.Path,
			"session_id", sessionID,
			"user_id", sessionStruct.UserID,
			"session_state", sessionStruct.State,
			"user_role", sessionStruct.Role)

		// Check session expiration
		if time.Now().After(sessionStruct.ExpiresAt) {
			logger.WarnContext(ctx, "Auth middleware: Session has expired for refresh token",
				"operation", "require_refresh_token",
				"method", r.Method,
				"path", r.URL.Path,
				"session_id", sessionID,
				"user_id", sessionStruct.UserID,
				"expires_at", sessionStruct.ExpiresAt,
				"current_time", time.Now())
			httpx.RespondWithError(w, errs.ErrExpiredToken, http.StatusUnauthorized)
			return
		}

		// Check session state - refresh tokens work for both pending and active sessions
		if sessionStruct.State != session.SessionActive && sessionStruct.State != session.SessionPending {
			logger.WarnContext(ctx, "Auth middleware: Invalid session state for refresh token",
				"operation", "require_refresh_token",
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
				"operation", "require_refresh_token",
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

		logger.InfoContext(ctx, "Auth middleware: Refresh token validation successful",
			"operation", "require_refresh_token",
			"method", r.Method,
			"path", r.URL.Path,
			"session_id", sessionID,
			"user_id", sessionStruct.UserID,
			"user_role", sessionStruct.Role)

		// Continue to next handler
		next(w, r)
	}
}
