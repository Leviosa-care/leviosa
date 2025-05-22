package userHandler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/jsonio"
)

func (a *AppInstance) GetUser(w http.ResponseWriter, r *http.Request) {
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
	user, err := a.Svcs.User.FindUserByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			logger.WarnContext(ctx, handler.NewNotFoundErr(err))
			http.Error(w, handler.NewNotFoundErr(err), http.StatusNotFound)
		case errors.Is(err, rp.ErrContext), errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, handler.NewInternalErr(err))
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrInvalidValue):
			logger.WarnContext(ctx, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			logger.WarnContext(ctx, "find user by ID:", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return
	}
	if err := jsonio.Encode(w, http.StatusOK, user); err != nil {
		logger.WarnContext(ctx, "get found user:", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return
	}
}
