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
		cookie, err := r.Cookie(AccessTokenCookieName)
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
		session, err := DecodeSession(sessionData)
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

		m.crypto.DecryptStruct(ctx, session)

		logger.InfoContext(ctx, "Auth middleware: Session retrieved and decrypted",
			"operation", "require_access_token",
			"method", r.Method,
			"path", r.URL.Path,
			"session_id", sessionID,
			"user_id", session.UserID,
			"session_state", session.State,
			"user_role", session.Role)

		// Check session state - access tokens work for both pending and active sessions
		if session.State != SessionActive && session.State != SessionPending {
			logger.WarnContext(ctx, "Auth middleware: Invalid session state",
				"operation", "require_access_token",
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
				"operation", "require_access_token",
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

		logger.InfoContext(ctx, "Auth middleware: Access token validation successful",
			"operation", "require_access_token",
			"method", r.Method,
			"path", r.URL.Path,
			"session_id", sessionID,
			"user_id", session.UserID,
			"user_role", session.Role)

		// Continue to next handler
		next(w, r)
	}
}
