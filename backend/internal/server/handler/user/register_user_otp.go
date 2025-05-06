package userHandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/contextutil"
	"github.com/hengadev/leviosa/pkg/serverutil"
)

func (h *AppInstance) RegisterUserOTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := contextutil.GetLoggerFromContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := serverutil.DecodeValid[models.UserSignUp](ctx, r.Body)
	if err != nil {
		switch {
		case errors.Is(err, serverutil.ErrDecodeJSON):
			logger.WarnContext(ctx, "decode user", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		default:
			logger.WarnContext(ctx, "invalid sign up user", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return
	}
	// TODO: I need to check the sent password against a list of leaked password
	if err := h.Svcs.User.CheckUser(ctx, user.Email); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			emailHash, err := h.createUser(ctx, w, logger, &user)
			if err != nil {
				return
			}
			if err := h.generateAndSendOTP(ctx, w, logger, emailHash, user.Email, user.FirstName); err != nil {
				return
			}
		case errors.Is(err, rp.ErrContext):
			logger.WarnContext(ctx, "context error, deadline or timeout while checking for user existence")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, "database checking for user existence query failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrUnexpectedType):
			logger.WarnContext(ctx, "unexpted errror checking for user existence")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return
	}

	logger.InfoContext(ctx, "unverified user successfully approved")
	w.WriteHeader(http.StatusCreated)
}

func (h *AppInstance) createUser(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, user *models.UserSignUp) (string, error) {
	emailHash, err := h.Svcs.User.CreateUnverifiedUser(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, "database adding unverified user query failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrNotEncrypted):
			logger.WarnContext(ctx, "fail to encrypt unverified user")
			http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		case errors.Is(err, domain.ErrNotCreated):
			logger.WarnContext(ctx, "database adding unverified user query failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, rp.ErrContext):
			logger.WarnContext(ctx, "context error, deadline or timeout while adding unverified user")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrUnexpectedType):
			logger.WarnContext(ctx, "unexpected errror adding unverified user")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return "", err
	}
	return emailHash, nil
}

// TODO: do the error handling on that
func (h *AppInstance) generateAndSendOTP(
	ctx context.Context,
	w http.ResponseWriter,
	logger *slog.Logger,
	emailHash string,
	userEmail string,
	firstname string,
) error {
	// generate OTP
	otp, err := h.Svcs.OTP.CreateOTP(ctx, emailHash)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, "database generate OTP failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrMarshalJSON):
			logger.WarnContext(ctx, "marshal JSON OTP failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, rp.ErrContext):
			logger.WarnContext(ctx, "context error, deadline or timeout while adding unverified user")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrRateLimit):
			logger.WarnContext(ctx, "too many requests")
			http.Error(w, handler.NewBadRequestErr(err), http.StatusTooManyRequests)
		default:
			logger.WarnContext(ctx, "failed to generate OTP")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return err
	}
	fmt.Printf("the OTP that I need to send back: %#+v\n", otp)
	// send email with OTP for user
	if err := h.Svcs.Mail.SendOTP(ctx, userEmail, firstname, otp); err != nil {
		logger.WarnContext(ctx, "failed to send mail with OTP to specified user")
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return err
	}
	return nil
}
