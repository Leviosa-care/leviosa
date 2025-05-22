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
	"github.com/hengadev/leviosa/pkg/serverutil"
)

// TODO: add the fact that creating an event, should also create a vote and a table with the style votes_month_year.
func (a *AppInstance) CreateEvent(w http.ResponseWriter, r *http.Request) {
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
	event, err := serverutil.Decode[eventModels.Event](r.Body)
	if err != nil {
		logger.WarnContext(ctx, "failed to decode event in create event handler", "error", err)
		http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		return
	}
	eventID, err := a.Svcs.Event.CreateEvent(ctx, &event)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFormat):
			logger.WarnContext(ctx, "failed to format event for database operation")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrNotCreated):
			logger.WarnContext(ctx, "failed to create event")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, "database creating event failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, rp.ErrContext):
			logger.WarnContext(ctx, "context error, deadline or timeout while creating event")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrUnexpectedType):
			logger.WarnContext(ctx, "unexpected error while creating event")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return
	}
	if err := serverutil.Encode(w, http.StatusCreated, eventID); err != nil {
		logger.WarnContext(ctx, "failed to send the event", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return
	}
}
