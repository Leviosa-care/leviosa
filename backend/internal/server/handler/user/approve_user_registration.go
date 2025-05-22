package userHandler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/jsonio"
)

func (h *AppInstance) ApproveUserRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := ctxutil.ValidateRoleInContext(ctx, models.ADMINISTRATOR); err != nil {
		logger.ErrorContext(ctx, "get role from request", "error", err)
		http.Error(w, handler.NewForbiddenErr(err), http.StatusBadRequest)
		return
	}
	// TODO:
	// get the email_hash and the role for the user
	input, err := jsonio.DecodeValid[models.UserPendingResponse](ctx, r.Body)

	// TODO: make the error right here for the client and for the logs
	if err != nil {
		switch {
		case errors.Is(err, jsonio.ErrDecodeJSON):
			logger.WarnContext(ctx, "failed to decode user", "error", err)
			http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
			return
		default:
			logger.WarnContext(ctx, "validate user")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
			return
		}
	}

	user, err := h.Svcs.User.CreateUser(ctx, &input)
	if err != nil {
		switch {
		default:
			logger.WarnContext(ctx, "failed to create account", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
			return
		}
	}
	// send email to user to tell them that their account have been approved
	if err := h.Svcs.Mail.WelcomeUser(ctx, user); err != nil {
		logger.WarnContext(ctx, "failed to send welcome email to new added user", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return

	}

	logger.InfoContext(ctx, "user successfully approved")
	http.Error(w, "user approved", http.StatusCreated)
}
