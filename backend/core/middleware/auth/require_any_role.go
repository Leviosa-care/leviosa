package auth

import (
	"net/http"
	"slices"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
	mw "github.com/Leviosa-care/core/middleware"
)

// RequireAnyRole validates access token and ensures user has one of the specified roles
func (m *SessionAuthMiddleware) RequireAnyRole(roles ...identity.Role) func(mw.Handler) mw.Handler {
	return func(next mw.Handler) mw.Handler {
		return m.RequireAccessToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger, err := ctxutil.GetLoggerFromContext(ctx)
			if err != nil {
				httpx.RespondWithError(w, err, http.StatusInternalServerError)
				return
			}

			sessionInfo, ok := session.SessionInfoFromContext(ctx)
			if !ok {
				logger.ErrorContext(ctx, "Auth middleware: Session info not found in context",
					"operation", "require_any_role",
					"method", r.Method,
					"path", r.URL.Path,
					"allowed_roles", roles)
				httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
				return
			}

			if !slices.Contains(roles, sessionInfo.Role) {
				logger.WarnContext(ctx, "Auth middleware: User role not in allowed roles",
					"operation", "require_any_role",
					"method", r.Method,
					"path", r.URL.Path,
					"user_id", sessionInfo.UserID,
					"user_role", sessionInfo.Role,
					"allowed_roles", roles)
				httpx.RespondWithError(w, errs.ErrForbidden, http.StatusForbidden)
				return
			}

			logger.InfoContext(ctx, "Auth middleware: Any role validation successful",
				"operation", "require_any_role",
				"method", r.Method,
				"path", r.URL.Path,
				"user_id", sessionInfo.UserID,
				"user_role", sessionInfo.Role,
				"allowed_roles", roles)

			next(w, r)
		}))
	}
}
