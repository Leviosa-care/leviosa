package auth

import (
	"net/http"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
	mw "github.com/Leviosa-care/core/middleware"
)

// RequireActiveSession validates access token and ensures session is in active state
// This middleware should be used for endpoints that require completed user registration
func (m *SessionAuthMiddleware) RequireActiveSession(next mw.Handler) mw.Handler {
	return m.RequireAccessToken(func(w http.ResponseWriter, r *http.Request) {
		// Get session info from context (already validated by RequireAccessToken)
		sessionInfo, ok := session.SessionInfoFromContext(r.Context())
		if !ok || sessionInfo == nil {
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// Verify session is in active state
		if sessionInfo.State != session.SessionActive {
			httpx.RespondWithError(w, errs.ErrForbidden, http.StatusForbidden)
			return
		}

		// Continue to next handler
		next(w, r)
	})
}
