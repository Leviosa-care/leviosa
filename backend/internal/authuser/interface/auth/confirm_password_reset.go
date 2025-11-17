package aggregatorHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
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

	var payload domain.ConfirmPasswordResetRequest

	// Check for password reset token in cookie first
	resetTokenFromCookie, err := getPasswordResetTokenFromCookies(r)
	if err == nil && resetTokenFromCookie != "" {
		payload.Token = resetTokenFromCookie
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "decode_request",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// If token wasn't provided in request body, ensure it's from cookie
	if payload.Token == "" {
		payload.Token = resetTokenFromCookie
	}

	logger.InfoContext(ctx, "Handler: Processing password reset confirmation request",
		"operation", "confirm_password_reset",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	if err := h.svc.ConfirmPasswordReset(ctx, &payload); err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "confirm password reset")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Password reset confirmation request completed successfully",
		"operation", "confirm_password_reset",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Clear the password reset token cookie (single-use)
	clearPasswordResetTokenCookie(w)

	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "Password reset completed successfully",
		Status:  "completed",
	}, http.StatusOK)
}
