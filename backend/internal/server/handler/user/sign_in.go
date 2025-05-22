package userHandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/session"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/jsonio"
)

func (a *AppInstance) Signin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	input, err := a.decodeAndValidateUserSignIn(ctx, w, r, logger)
	if err != nil {
		return
	}

	if err = a.checkThrottling(ctx, w, input.Email, logger); err != nil {
		return
	}

	if err = a.validateUserUserSignIn(ctx, w, input, logger); err != nil {
		return
	}

	if err = a.createAndSetSession(ctx, w, input.Email, logger); err != nil {
		return
	}

	logger.InfoContext(ctx, "User successfully logged in", slog.String("key", "value"))
}

func (a *AppInstance) decodeAndValidateUserSignIn(ctx context.Context, w http.ResponseWriter, r *http.Request, logger *slog.Logger) (*models.UserSignIn, error) {
	input, err := jsonio.DecodeValid[models.UserSignIn](ctx, r.Body)
	if err != nil {
		switch {
		case errors.Is(err, jsonio.ErrDecodeJSON):
			logger.WarnContext(ctx, err.Error())
			http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
			return nil, err
		default:
			logger.WarnContext(ctx, "invalid user credentials", "error", err)
			http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
			return nil, err
		}
	}
	return &input, nil
}

func (h *AppInstance) checkThrottling(ctx context.Context, w http.ResponseWriter, email string, logger *slog.Logger) error {
	isLocked, err := h.Repos.Throttler.IsLocked(ctx, email)
	if isLocked != nil {
		logger.WarnContext(ctx, "user locked, too many attempts")
		http.Error(w, handler.NewGetErr("throttler locking status", err), http.StatusTooManyRequests)
		return fmt.Errorf("user locked")
	}
	switch {
	case errors.Is(err, rp.ErrNotFound):
		logger.DebugContext(ctx, "first sign in attempt", slog.String("user", email))
	case errors.Is(err, rp.ErrDatabase):
		logger.ErrorContext(ctx, "failed database")
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return err
	default:
		logger.WarnContext(ctx, "unhandled error type", slog.String("error", err.Error()))
		http.Error(w, "unexpected error occurred", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (a *AppInstance) validateUserUserSignIn(ctx context.Context, w http.ResponseWriter, input *models.UserSignIn, logger *slog.Logger) error {
	err := a.Svcs.User.ValidateCredentials(ctx, input)
	switch {
	case errors.Is(err, domain.ErrNotFound):
		logger.WarnContext(ctx, "invalid email", "error", err)
		http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		return err
	case errors.Is(err, domain.ErrQueryFailed):
		logger.ErrorContext(ctx, "database error", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return err
	case errors.Is(err, domain.ErrUnexpectedType):
		logger.WarnContext(ctx, "unhandled error type", slog.String("error", err.Error()))
		http.Error(w, "unexpected error occurred", http.StatusInternalServerError)
		return err
	case err != nil:
		attemptErr := a.Svcs.Throttler.RegisterAttempt(ctx, input.Email)
		logger.WarnContext(ctx, "invalid password", "error", errors.Join(err, attemptErr))
		http.Error(w, handler.NewBadRequestErr(errors.Join(err, attemptErr)), http.StatusBadRequest)
		return err
	}
	return nil
}

func (a *AppInstance) createAndSetSession(ctx context.Context, w http.ResponseWriter, email string, logger *slog.Logger) error {
	userID, role, err := a.Svcs.User.GetUserSessionData(ctx, email)
	switch {
	case errors.Is(err, domain.ErrNotFound):
		logger.WarnContext(ctx, "invalid email", "error", err)
		http.Error(w, handler.NewGetErr("user session data", err), http.StatusBadRequest)
		return err
	case errors.Is(err, domain.ErrQueryFailed):
		logger.ErrorContext(ctx, "database error", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return err
	case errors.Is(err, domain.ErrUnexpectedType):
		logger.WarnContext(ctx, "unhandled error type", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("unexpected error occurred: %s", err.Error()), http.StatusInternalServerError)
		return err
	case err != nil:
		attemptErr := a.Svcs.Throttler.RegisterAttempt(ctx, email)
		logger.WarnContext(ctx, "get userID and role with email", "error", errors.Join(err, attemptErr))
		http.Error(w, handler.NewBadRequestErr(errors.Join(err, attemptErr)), http.StatusBadRequest)
		return err
	}
	sessionID, err := a.Svcs.Session.CreateSession(ctx, userID, role)
	switch {
	case errors.Is(err, domain.ErrQueryFailed):
		logger.ErrorContext(ctx, "database error", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return err
	case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
		logger.WarnContext(ctx, err.Error())
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return err
	case errors.Is(err, domain.ErrUnexpectedType):
	}
	err = a.Svcs.Throttler.Reset(ctx, email)
	switch {
	case errors.Is(err, domain.ErrUnmarshalJSON):
		logger.WarnContext(ctx, err.Error())
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return err
	case errors.Is(err, domain.ErrNotFound):
		logger.WarnContext(ctx, err.Error())
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return err
	case errors.Is(err, domain.ErrQueryFailed):
		logger.ErrorContext(ctx, "database error", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
	case errors.Is(err, domain.ErrAccountLocked):
		logger.WarnContext(ctx, "user locker, too many attempts")
		http.Error(w, handler.NewGetErr("throttler locking status", err), http.StatusTooManyRequests)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionService.SessionName,
		Value:    sessionID,
		Expires:  time.Now().Add(sessionService.SessionDuration),
		HttpOnly: true,
	})
	return nil
}
