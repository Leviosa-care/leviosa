package auth

import (
	"context"
	"errors"
	"net/http"

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
		if r.URL.Path != RefreshEndpoint {
			logger.WarnContext(ctx, "Auth middleware: Refresh token attempted on wrong endpoint",
				"operation", "require_refresh_token",
				"method", r.Method,
				"path", r.URL.Path,
				"expected_path", RefreshEndpoint)
			httpx.RespondWithError(w, errs.ErrForbidden, http.StatusForbidden)
			return
		}

		// Extract refresh token from cookies
		cookie, err := r.Cookie(RefreshTokenCookieName)
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

		// Find session by refresh token using two-step lookup
		sessionID, sessionData, err := m.sessionRepo.FindSessionByRefreshToken(ctx, refreshToken)
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
		session, err := DecodeSession(sessionData)
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

		err = m.crypto.DecryptStruct(ctx, session)
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
			"user_id", session.UserID,
			"session_state", session.State,
			"user_role", session.Role)

		// Check session state - refresh tokens work for both pending and active sessions
		if session.State != SessionActive && session.State != SessionPending {
			logger.WarnContext(ctx, "Auth middleware: Invalid session state for refresh token",
				"operation", "require_refresh_token",
				"method", r.Method,
				"path", r.URL.Path,
				"session_id", sessionID,
				"user_id", session.UserID,
				"session_state", session.State,
				"expected_states", []string{string(SessionActive), string(SessionPending)})
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
		sessionInfo := &SessionInfo{
			ID:     parsedSessionID,
			UserID: session.UserID,
			Role:   session.Role,
			State:  session.State,
		}

		// Add session info to context
		ctx = context.WithValue(ctx, sessionContextKey{}, sessionInfo)
		r = r.WithContext(ctx)

		logger.InfoContext(ctx, "Auth middleware: Refresh token validation successful",
			"operation", "require_refresh_token",
			"method", r.Method,
			"path", r.URL.Path,
			"session_id", sessionID,
			"user_id", session.UserID,
			"user_role", session.Role)

		// Continue to next handler
		next(w, r)
	}
}

