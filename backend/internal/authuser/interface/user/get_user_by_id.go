package userHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract user ID from path parameter
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		logger.WarnContext(ctx, "Handler: Missing user ID in path",
			"operation", "get_user_by_id",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Parse user ID as UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid user ID format",
			"operation", "get_user_by_id",
			"method", r.Method,
			"path", r.URL.Path,
			"user_id", userIDStr,
			"error", err)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get user by ID request",
		"operation", "get_user_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", userID,
		"user_agent", r.Header.Get("User-Agent"))

	user, err := h.svc.GetUserByID(ctx, userID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get user by ID")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get user by ID completed",
		"operation", "get_user_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", userID,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, user, http.StatusOK)
}
