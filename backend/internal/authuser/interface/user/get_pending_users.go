package userHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetPendingUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get pending users request",
		"operation", "get_pending_users",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	users, err := h.svc.GetPendingUsers(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get pending users")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get pending users completed",
		"operation", "get_pending_users",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"user_count", len(users))

	httpx.RespondWithJSON(w, users, http.StatusOK)
}
