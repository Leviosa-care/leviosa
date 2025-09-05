package auth

import (
	"net/http"

	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
	mw "github.com/Leviosa-care/core/middleware"
)

// RequireMinimumRole validates access token and ensures user has at least the specified role
func (m *SessionAuthMiddleware) RequireMinimumRole(minRole identity.Role) func(mw.Handler) mw.Handler {
	return func(next mw.Handler) mw.Handler {
		return m.RequireAccessToken(func(w http.ResponseWriter, r *http.Request) {
			sessionInfo, ok := SessionInfoFromContext(r.Context())
			if !ok {
				httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
				return
			}

			if !sessionInfo.Role.IsAtLeast(minRole) {
				httpx.RespondWithError(w, errs.ErrForbidden, http.StatusForbidden)
				return
			}

			next(w, r)
		})
	}
}