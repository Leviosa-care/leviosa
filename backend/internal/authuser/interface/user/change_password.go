package userHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get session info from context (set by middleware)
	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		logger.ErrorContext(ctx, "Handler: No session info in context",
			"operation", "change_password",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing change password request",
		"operation", "change_password",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", sessionInfo.UserID,
		"user_agent", r.Header.Get("User-Agent"))

	// Parse request body
	var request domain.ChangePasswordRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "change_password",
			"method", r.Method,
			"path", r.URL.Path,
			"user_id", sessionInfo.UserID)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to change password
	err = h.svc.ChangePassword(ctx, sessionInfo.UserID, &request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "change password")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Change password completed",
		"operation", "change_password",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"user_id", sessionInfo.UserID)

	// Return success response
	httpx.RespondWithJSON(w, map[string]string{"message": "Password changed successfully"}, http.StatusOK)
}
