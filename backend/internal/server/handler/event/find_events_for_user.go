package eventHandler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/jsonio"
)

func (a *AppInstance) FindEventsForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := ctxutil.ValidateRoleInContext(ctx, models.STANDARD); err != nil {
		logger.WarnContext(ctx, "get role from request", "error", err)
		http.Error(w, handler.NewForbiddenErr(err), http.StatusBadRequest)
		return
	}
	userID, ok := ctx.Value(ctxutil.UserIDKey).(string)
	if !ok {
		logger.ErrorContext(ctx, "user ID not found in context")
		http.Error(w, errors.New("failed to get user ID from context").Error(), http.StatusInternalServerError)
		return
	}

	userEvents, err := a.Repos.Event.GetEventForUser(ctx, userID)
	if err != nil {
		switch {
		default:
			logger.ErrorContext(ctx, "failed to get the events for the user", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return
	}
	if err := jsonio.Encode(w, http.StatusOK, userEvents); err != nil {
		logger.ErrorContext(ctx, "failed to send the user", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return
	}
}
