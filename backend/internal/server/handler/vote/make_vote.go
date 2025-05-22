package vote

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/domain/vote"
	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/serverutil"
)

// Function that create or update vote
func (a *AppInstance) MakeVote() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// get votes from request
		votes, err := serverutil.Decode[[]*vote.Vote](r.Body)
		if err != nil {
			logger.ErrorContext(ctx, "failed to decode vote", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
			return
		}
		// add userID field to votes.
		for _, vote := range votes {
			vote.UserID = userID
		}
		// create vote
		if err := a.Svcs.Vote.CreateVote(ctx, votes); err != nil {
			logger.ErrorContext(ctx, "failed to create vote", "error", err)
			http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
			return

		}
	})
}
