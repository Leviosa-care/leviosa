package vote

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/server/handler"
	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/jsonio"
)

// Function that get all the votes for a user.
func (a *AppInstance) GetVotesByUserID(w http.ResponseWriter, r *http.Request) {
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

	month := r.PathValue("month")
	year := r.PathValue("year")
	if month == "" || year == "" {
		err := fmt.Errorf("month or year are not well formatted")
		logger.ErrorContext(ctx, "failed to parse month or year", "error", err)
		http.Error(w, handler.NewBadRequestErr(err), http.StatusBadRequest)
		return
	}
	// get votes
	votes, err := a.Svcs.Vote.GetVotesByUserID(ctx, month, year, userID)
	if err != nil {
		logger.ErrorContext(ctx, "failed to get the votes from database", "error", err)
		http.Error(w, fmt.Sprintf("Failed to get the data from the database - %s", err), http.StatusInternalServerError)
		return
	}
	// encode result to user
	if err := jsonio.Encode(w, http.StatusOK, votes); err != nil {
		logger.ErrorContext(ctx, "failed to encode votes found for user with provided ID", "error", err)
		http.Error(w, handler.NewInternalErr(err), http.StatusInternalServerError)
		return
	}
}
