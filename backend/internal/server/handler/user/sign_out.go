package userHandler

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/session"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
)

func (a *AppInstance) Signout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get session value from cookie
	cookie, err := r.Cookie(sessionService.SessionName)
	if err != nil {
		logger.ErrorContext(ctx, "get session cookie for signout", "error", err)
		http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		return

	}
	sessionID := cookie.Value

	// remove session with sessionID
	if err := a.Svcs.Session.RemoveSession(ctx, sessionID); err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidValue):
			logger.WarnContext(ctx, "invalid value in session validation")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrQueryFailed):
			logger.WarnContext(ctx, "database removing session for user query failed")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrNotFound):
			logger.WarnContext(ctx, "user session not found")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		case errors.Is(err, domain.ErrUnexpectedType):
			logger.WarnContext(ctx, "unexpected error removing user session")
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		}
		return
	}
	// send cookie to inform client that cookie is no longer valid
	http.SetCookie(w, &http.Cookie{
		Name:    sessionService.SessionName,
		Value:   "",
		Expires: time.Now(),
	})
}
