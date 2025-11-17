package aggregatorHandler

import (
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"

	"github.com/google/uuid"
)

func (h *handler) DeleteUserByAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract logger from context
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract user ID from URL path
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("user ID is required"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing delete user by admin request",
		"operation", "delete_user_by_admin",
		"user_id", userIDStr,
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid user ID format: %v", err)), http.StatusBadRequest)
		return
	}

	// Call the aggregator service to delete the user
	err = h.svc.DeleteUserByAdmin(ctx, userID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "delete user by admin")
		return
	}

	// Log successful completion
	logger.InfoContext(ctx, "Handler: Delete user by admin completed",
		"operation", "delete_user_by_admin",
		"user_id", userID.String(),
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with success message
	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
	}{
		Message: "User deleted successfully",
	}, http.StatusOK)
}
