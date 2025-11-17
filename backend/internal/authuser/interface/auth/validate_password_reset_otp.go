package aggregatorHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

// Password reset token cookie configuration
const (
	PasswordResetTokenCookieName = "leviosa_password_reset_token"
	PasswordResetEndpoint        = "/auth/password/reset/confirm"
)

// setPasswordResetTokenCookie sets the password reset token cookie with appropriate security settings
func setPasswordResetTokenCookie(w http.ResponseWriter, token string, expiry time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     PasswordResetTokenCookieName,
		Value:    token,
		Path:     PasswordResetEndpoint, // Restrict to reset confirmation endpoint only
		HttpOnly: true,
		Secure:   true, // Only sent over HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  expiry,
	})
}

// getPasswordResetTokenFromCookies retrieves the password reset token from cookies
func getPasswordResetTokenFromCookies(r *http.Request) (string, error) {
	cookie, err := r.Cookie(PasswordResetTokenCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", fmt.Errorf("password reset token cookie not found: %w", err)
		}
		return "", fmt.Errorf("failed to read password reset token cookie: %w", err)
	}
	return cookie.Value, nil
}

// clearPasswordResetTokenCookie clears the password reset token cookie
func clearPasswordResetTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     PasswordResetTokenCookieName,
		Value:    "",
		Path:     PasswordResetEndpoint,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Delete the cookie immediately
	})
}

func (h *handler) ValidatePasswordResetOTP(w http.ResponseWriter, r *http.Request) {
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

	var payload domain.ValidatePasswordResetOTPRequest

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

	logger.InfoContext(ctx, "Handler: Processing password reset OTP validation request",
		"email", maskEmail(payload.Email),
		"operation", "validate_password_reset_otp",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	response, err := h.svc.ValidatePasswordResetOTP(ctx, &payload)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "validate password and reset OTP")
		return
	}

	// Log successful operation with enhanced context
	logger.InfoContext(ctx, "Handler: Password reset OTP validation request completed successfully",
		"email", maskEmail(payload.Email),
		"operation", "validate_password_reset_otp",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Set password reset token cookie
	setPasswordResetTokenCookie(w, response.Token, response.ExpiresAt)

	httpx.RespondWithJSON(w, struct {
		Message   string    `json:"message"`
		Status    string    `json:"status"`
		ExpiresAt time.Time `json:"expires_at"`
	}{
		Message:   "Password reset OTP validated successfully",
		Status:    "validated",
		ExpiresAt: response.ExpiresAt,
	}, http.StatusOK)
}
