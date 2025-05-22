package userHandler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/jsonio"
)

func (a *AppInstance) CheckUserExists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	email, err := jsonio.Decode[models.Email](r.Body)
	if err != nil {
		switch {
		case errors.Is(err, jsonio.ErrDecodeJSON):
			logger.WarnContext(ctx, err.Error())
			http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		default:
			logger.WarnContext(ctx, "invalid email", "error", err)
			http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		}
	}

	type Response struct {
		Exists bool `json:"exists"`
	}

	err = a.Svcs.User.CheckUser(ctx, email.String())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			if err := jsonio.Encode(w, http.StatusOK, Response{Exists: false}); err != nil {
				logger.ErrorContext(ctx, "failed to encode pendings users list", "error", err)
				http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
			}
		case errors.Is(err, rp.ErrContext):
			logger.WarnContext(ctx, "context error", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, "database query check user email failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrUnexpectedType):
			logger.WarnContext(ctx, "unexpected error while trying to check user")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
	}
	if err := jsonio.Encode(w, http.StatusOK, Response{Exists: true}); err != nil {
		logger.ErrorContext(ctx, "failed to encode pendings users list", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
	}
}
