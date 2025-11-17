package aggregatorHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) DeleteOwnAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract logger from context
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get session info from context (user ID comes from here)
	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("session information from context required"), http.StatusUnauthorized)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing delete own account request",
		"operation", "delete_own_account",
		"user_id", sessionInfo.UserID.String(),
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Call the aggregator service to delete the user's own account
	err = h.svc.DeleteOwnAccount(ctx, sessionInfo)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "delete user's own account")
		return
	}

	// Log successful completion
	logger.InfoContext(ctx, "Handler: Delete own account completed",
		"operation", "delete_own_account",
		"user_id", sessionInfo.UserID.String(),
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with success message
	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
	}{
		Message: "Account deleted successfully",
	}, http.StatusOK)
}
