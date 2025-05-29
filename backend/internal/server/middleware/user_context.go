package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"slices"

	"github.com/hengadev/leviosa/pkg/ctxutil"
)

// get the role for the user in question
func SetUserContext(sessionGetter sessionGetterFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			exceptPath := []string{
				"/healthz",
				"/hello",
				"/auth/email",
				"/auth/register",
				"/auth/validate-otp",
				"/auth/approve-user",
				"/oauth/google/user",
				"/upload-image",
			}
			if slices.Contains(exceptPath, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			logger, ok := ctx.Value(ctxutil.LoggerKey).(*slog.Logger)
			if !ok {
				slog.ErrorContext(ctx, "logger not found in context")
				http.Error(w, "logger not found in context", http.StatusInternalServerError)
				return
			}

			// get session ID from request
			sessionID, err := getSessionIDFromRequest(r)
			if err != nil {
				logger.ErrorContext(ctx, "get sessionID from request header", "error", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// get session object from database
			session, err := sessionGetter(ctx, sessionID)
			if err != nil {
				logger.ErrorContext(ctx, "get session from database", "error", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// validate session
			if err := session.Valid(ctx); err != nil {
				logger.ErrorContext(ctx, "invalid session", "error", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// set the role in the context using the session
			ctx = context.WithValue(r.Context(), ctxutil.RoleKey, session.Role.String())
			ctx = context.WithValue(r.Context(), ctxutil.UserIDKey, session.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
