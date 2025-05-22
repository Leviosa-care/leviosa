package eventHandler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/domain"
	eventModels "github.com/hengadev/leviosa/internal/domain/event/models"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/jsonio"
)

func (a *AppInstance) ModifyEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := ctxutil.ValidateRoleInContext(ctx, models.ADMINISTRATOR); err != nil {
		logger.WarnContext(ctx, "get role from request", "error", err)
		http.Error(w, handler.NewForbiddenErr(err), http.StatusBadRequest)
		return
	}
	event, err := jsonio.Decode[eventModels.Event](r.Body)
	if err != nil {
		logger.WarnContext(ctx, "failed to decode the event", "error", err)
		http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		return
	}
	err = a.Svcs.Event.ModifyEvent(ctx, &event)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrContext):
			logger.WarnContext(ctx, "context error, deadline or timeout while updating event")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, "database updating event failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrNotUpdated):
			logger.WarnContext(ctx, "failed to update event")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrUnexpectedType):
			logger.WarnContext(ctx, "unexpected error while updating event")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrInvalidValue):
			logger.WarnContext(ctx, "invalid event error while updating event")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return
	}
	if err = jsonio.Encode(w, http.StatusInternalServerError, event.ID); err != nil {
		logger.WarnContext(ctx, "failed to send event ID", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return
	}
}
