package auth

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

// RequireMinimumRole validates access token and ensures user has at least the specified role
func (m *SessionAuthMiddleware) RequireMinimumRole(minRole identity.Role) func(mw.Handler) mw.Handler {
	return func(next mw.Handler) mw.Handler {
		return m.RequireAccessToken(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger, err := ctxutil.GetLoggerFromContext(ctx)
			if err != nil {
				httpx.RespondWithError(w, err, http.StatusInternalServerError)
				return
			}

			sessionInfo, ok := session.SessionInfoFromContext(ctx)
			if !ok {
				logger.ErrorContext(ctx, "Auth middleware: Session info not found in context",
					"operation", "require_minimum_role",
					"method", r.Method,
					"path", r.URL.Path,
					"required_role", minRole)
				httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
				return
			}

			if !sessionInfo.Role.IsAtLeast(minRole) {
				logger.WarnContext(ctx, "Auth middleware: Insufficient role for access",
					"operation", "require_minimum_role",
					"method", r.Method,
					"path", r.URL.Path,
					"user_id", sessionInfo.UserID,
					"user_role", sessionInfo.Role,
					"required_role", minRole)
				httpx.RespondWithError(w, errs.ErrForbidden, http.StatusForbidden)
				return
			}

			logger.InfoContext(ctx, "Auth middleware: Minimum role validation successful",
				"operation", "require_minimum_role",
				"method", r.Method,
				"path", r.URL.Path,
				"user_id", sessionInfo.UserID,
				"user_role", sessionInfo.Role,
				"required_role", minRole)

			next(w, r)
		})
	}
}
