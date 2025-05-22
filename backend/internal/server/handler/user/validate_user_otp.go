package userHandler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/jsonio"
)

func (h *AppInstance) ValidateUserOTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("get otp and email from payload")
	userOTP, err := jsonio.DecodeValid[models.UserOTP](r.Context(), r.Body)
	if err != nil {
		switch {
		case errors.Is(err, jsonio.ErrDecodeJSON):
			logger.WarnContext(ctx, "decode user", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		default:
			logger.WarnContext(ctx, "invalid sign up user", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return
	}
	fmt.Printf("the userOTP that I get: %#+v\n", userOTP)

	// validate otp
	if err = h.Svcs.OTP.ValidateOTP(ctx, userOTP.Email, userOTP.OTP); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			logger.WarnContext(ctx, "OTP validation query failed due to database error")
			http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		case errors.Is(err, rp.ErrContext):
			logger.WarnContext(ctx, "context error, deadline or timeout while adding pending user")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, "database compare OTPs failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrInvalidValue):
			logger.WarnContext(ctx, "invalid OTP provided")
			http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		case errors.Is(err, domain.ErrUnexpectedType):
			logger.WarnContext(ctx, "unexpected error while validation OTP sent")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return
	}

	// add to pending user (using the hash email)
	if err = h.Svcs.User.CreatePendingUser(ctx, userOTP.Email); err != nil {
		switch {
		case errors.Is(err, rp.ErrContext):
			logger.WarnContext(ctx, "context error, deadline or timeout while adding pending user")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrNotFound):
			logger.WarnContext(ctx, "user not found in unverified_users table")
			http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		case errors.Is(err, domain.ErrUnexpectedType):
			logger.WarnContext(ctx, "unexpected error while adding user to pending_users table")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrNotCreated):
			logger.WarnContext(ctx, "database creating user to pending_users table failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, "database adding user to pending_users table failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return
	}

	logger.InfoContext(ctx, "pending user successfully approved")
	http.Error(w, "pending user successfully created", http.StatusCreated)
}
