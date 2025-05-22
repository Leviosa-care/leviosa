package userHandler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
)

func (a *AppInstance) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, ok := ctx.Value(ctxutil.UserIDKey).(string)
	if !ok {
		logger.ErrorContext(ctx, "user ID not found in context")
		http.Error(w, errors.New("failed to get user ID from context").Error(), http.StatusInternalServerError)
		return
	}

	// get the sessionID from the cookie
	sessionID := r.Cookies()[0].Value
	// - call the user repo to delete the user account from the user repository
	err = a.Svcs.User.DeleteUser(ctx, userID)
	if err != nil {
		logger.ErrorContext(ctx, "delete user:", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return
	}
	// - call the service repo to delete all user's session
	err = a.Svcs.Session.RemoveSession(ctx, sessionID)
	if err != nil {
		logger.ErrorContext(ctx, "delete user session:", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return
	}
}
