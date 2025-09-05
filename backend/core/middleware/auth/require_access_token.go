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