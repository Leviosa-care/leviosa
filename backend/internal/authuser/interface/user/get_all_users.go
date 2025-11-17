package userHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get all users request",
		"operation", "get_all_users",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	users, err := h.svc.GetAllUsers(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get all users")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get all users completed",
		"operation", "get_all_users",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"user_count", len(users))

	httpx.RespondWithJSON(w, users, http.StatusOK)
}
