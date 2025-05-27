package userHandler

import (
	"context"
	"errors"
	"io"
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

func (a *AppInstance) HandleOAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var provider models.ProviderType
	inputProvider := r.PathValue("provider")
	if err = provider.Set(inputProvider); err != nil {
		logger.ErrorContext(ctx, "invalid provider", "error", err)
		http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		return
	}

	user, err := a.decodeAndValidUser(ctx, w, r.Body, logger, provider)
	if err != nil {
		return
	}
	if err := a.handleUser(ctx, w, user, logger, provider); err != nil {
		return
	}
	if err := a.handleSession(ctx, w, user.Email, logger); err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *AppInstance) decodeAndValidUser(ctx context.Context, w http.ResponseWriter, body io.ReadCloser, logger *slog.Logger, provider models.ProviderType) (*models.User, error) {
	// put that into a decodeAndValid function
	var oauthUser models.OAuthUser
	var user *models.User
	var err error
	switch provider {
	case models.Google:
		oauthUser, err = jsonio.DecodeValid[models.GoogleUser](ctx, body)
		user = oauthUser.ToUser()
	case models.Apple:
		oauthUser, err = jsonio.DecodeValid[models.AppleUser](ctx, body)
		user = oauthUser.ToUser()
	}
	if err != nil {
		switch {
		case errors.Is(err, jsonio.ErrDecodeJSON):
			logger.WarnContext(ctx, err.Error())
			http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		case errors.Is(err, jsonio.ErrValidStruct):
			logger.WarnContext(ctx, "invalid struct")
			http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		default:
			logger.WarnContext(ctx, "invalid decode valid")
			http.Error(w, handler.NewInternalErr(err), http.StatusBadRequest)
		}
		return nil, err
	}
	return user, nil
}

func (a *AppInstance) handleUser(ctx context.Context, w http.ResponseWriter, user *models.User, logger *slog.Logger, provider models.ProviderType) error {
	if err := a.Svcs.User.CheckOAuthUser(ctx, user.Email, provider); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			if err := a.Svcs.User.CreateOAuthPendingUser(ctx, user, provider); err != nil {
				switch {
				case errors.Is(err, domain.ErrInvalidValue):
					logger.WarnContext(ctx, "invalid value")
					http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
				case errors.Is(err, rp.ErrContext):
					logger.WarnContext(ctx, "context error, deadline or timeout while checking for user existence")
					http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
				case errors.Is(err, domain.ErrNotCreated):
					logger.WarnContext(ctx, "failed to create oauth user")
					http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
				case errors.Is(err, domain.ErrUnexpectedType):
					logger.WarnContext(ctx, "unexpected errror adding oauth user")
					http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
				}
				return err
			}
			// TODO: handle that part
			// if err := a.Svcs.Mail.SendOTP(ctx, user); err != nil {
			// 	logger.WarnContext(ctx, "sending mail to welcome new oauth pending user")
			// 	http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
			// 	return err
			// }
		case errors.Is(err, rp.ErrContext):
			logger.WarnContext(ctx, "context error, deadline or timeout while checking for user existence")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, "database checking for oauth user existence query failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrInvalidValue):
			logger.WarnContext(ctx, "check for invalid email for oauth user existence query failed")
			http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		case errors.Is(err, domain.ErrUnexpectedType):
			logger.WarnContext(ctx, "unexpected errror checking for oauth user existence in pending_users table")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return err
	}
	return nil
}

func (a *AppInstance) handleSession(ctx context.Context, w http.ResponseWriter, email string, logger *slog.Logger) error {
	userID, role, err := a.Svcs.User.GetUserSessionData(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidValue):
			logger.WarnContext(ctx, "invalid value in getting oauth user session data")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrNotFound):
			logger.WarnContext(ctx, "user session data (userID, role) not found in database")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, rp.ErrContext):
			logger.WarnContext(ctx, "context error, deadline or timeout while checking for user existence")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, "database getting user session data query failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrUnexpectedType):
			logger.WarnContext(ctx, "unexpected error getting oauth user session data")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
	}

	sessionID, err := a.Svcs.Session.CreateSession(ctx, userID, role)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidValue):
			logger.WarnContext(ctx, "invalid value in session validation")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrMarshalJSON):
			logger.WarnContext(ctx, "marshal session data for oauth user")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, "database creating session for oauth user query failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrUnexpectedType):
			logger.WarnContext(ctx, "unexpected error creating session for oauth user")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionService.SessionName,
		Value:    sessionID,
		Expires:  time.Now().Add(sessionService.SessionDuration),
		HttpOnly: true,
	})
	return nil
}
