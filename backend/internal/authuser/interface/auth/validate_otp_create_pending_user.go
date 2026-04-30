package aggregatorHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) ValidateOTPCreatePendingUser(w http.ResponseWriter, r *http.Request) {
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

	var payload domain.ValidateOTPRequest

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

	logger.InfoContext(ctx, "Handler: Processing OTP verification request",
		"email", maskEmail(payload.Email),
		"operation", "validate_otp_create_pending_user",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	session, err := h.svc.ValidateOTPCreatePendingUser(ctx, &payload)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "validate OTP and create pending user")
		return
	}

	// Log successful operation with enhanced context
	logger.InfoContext(ctx, "Handler: OTP verification request completed successfully",
		"email", maskEmail(payload.Email),
		"operation", "validate_otp_create_pending_user",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusCreated)

	// Set dual token cookies
	cookies.SetTokenCookies(w, session.AccessToken, session.RefreshToken,
		session.AccessTokenExpiry, session.RefreshTokenExpiry)

	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "Pending user created successfully",
		Status:  "created",
	}, http.StatusCreated)
}
