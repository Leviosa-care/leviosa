package userHandler

import (
	"context"
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

type EmailRequest struct {
	Email models.Email `json:"email"`
}

func (e EmailRequest) Valid(ctx context.Context) error {
	return e.Email.Valid(ctx)
}

func (h *AppInstance) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("here is some Info message")
	logger.Debug("here is some Debug message")
	data, err := jsonio.DecodeValid[EmailRequest](ctx, r.Body)
	if err != nil {
		switch {
		case errors.Is(err, jsonio.ErrDecodeJSON):
			logger.WarnContext(ctx, "decode user", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusBadRequest)
		default:
			logger.WarnContext(ctx, "invalid sign up user", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusUnprocessableEntity)
		}
		return
	}
	email := data.Email.String()
	if err := h.Svcs.User.CheckUser(ctx, email); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			print("email not found\n")
			otp, err := h.Svcs.OTP.RequestOTP(ctx, email)
			if err != nil {
				print("error request OTP:", err.Error(), "\n")
				switch {
				case errors.Is(err, domain.ErrQueryFailed):
					logger.WarnContext(ctx, "context error, deadline or timeout while checking for user existence", "error", err)
					http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
				case errors.Is(err, domain.ErrUnmarshalJSON):
					logger.WarnContext(ctx, "error unmarshaling JSON", "error", err)
					http.Error(w, handler.NewInternalErr(err), http.StatusBadRequest)
				case errors.Is(err, domain.ErrNotDecrypted):
					logger.WarnContext(ctx, "decrypt OTP data", "error", err)
					http.Error(w, handler.NewInternalErr(err), http.StatusBadRequest)
				case errors.Is(err, domain.ErrRateLimit):
					logger.WarnContext(ctx, "rate limit reached", "error", err)
					http.Error(w, handler.NewInternalErr(err), http.StatusBadRequest)
				case errors.Is(err, domain.ErrNotCreated):
					logger.WarnContext(ctx, "OTP not created", "error", err)
					http.Error(w, handler.NewInternalErr(err), http.StatusBadRequest)
				case errors.Is(err, domain.ErrNotEncrypted):
					logger.WarnContext(ctx, "OTP not encrypted", "error", err)
					http.Error(w, handler.NewInternalErr(err), http.StatusBadRequest)
				case errors.Is(err, domain.ErrMarshalJSON):
					logger.WarnContext(ctx, "error marshal JSON", "error", err)
					http.Error(w, handler.NewInternalErr(err), http.StatusBadRequest)
				}
				return
			}
			print("otp found:", otp, "\n")
			if err := h.Svcs.Mail.SendOTP(ctx, email, otp); err != nil {
				logger.WarnContext(ctx, "failed to send mail with OTP to specified user", "error", err)
				http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
			}
			print("otp sent to mail \n")
			logger.InfoContext(ctx, "OTP generated and send for unverified user")
			w.WriteHeader(http.StatusOK)
			return
		case errors.Is(err, rp.ErrContext):
			logger.WarnContext(ctx, "context error, deadline or timeout while checking for user existence", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, "database checking for user existence query failed", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrUnexpectedType):
			logger.WarnContext(ctx, "unexpted errror checking for user existence", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusConflict)
}
