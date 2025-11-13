package aggregatorHandler

import (
	"context"
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
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "warn"
			errorContext = "validation failed"
		case errors.Is(err, errs.ErrDomainNotFound):
			logLevel = "warn"
			errorContext = "resource not found"
		case errors.Is(err, errs.ErrExpiredToken):
			logLevel = "info"
			errorContext = "token expired"
		case errors.Is(err, errs.ErrRateLimit):
			logLevel = "warn"
			errorContext = "rate limit exceeded"
		case errors.Is(err, errs.ErrValueMismatch):
			logLevel = "warn"
			errorContext = "authentication failed"
		case errors.Is(err, errs.ErrAlreadyConsumed):
			logLevel = "warn"
			errorContext = "resource already consumed"
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			logLevel = "error"
			errorContext = "infrastructure connection failure"
		case errors.Is(err, errs.ErrResourceExhausted):
			logLevel = "error"
			errorContext = "infrastructure resource exhaustion"
		case errors.Is(err, errs.ErrExternalService):
			logLevel = "error"
			errorContext = "external service failure"
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			logLevel = "warn"
			errorContext = "request cancelled"
		case errors.Is(err, context.DeadlineExceeded):
			logLevel = "warn"
			errorContext = "request timeout"
		case errors.Is(err, errs.ErrTransactionFailure):
			logLevel = "error"
			errorContext = "infrastructure transaction failure"
		default:
			logLevel = "error"
			errorContext = "unexpected error"
		}

		logFields := []any{
			"email", maskEmail(payload.Email),
			"operation", "validate_password_reset_otp",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"error", err,
		}

		var statusCode int
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			// Invalid OTP format or email validation
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrDomainNotFound):
			// OTP doesn't exist or already used
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrExpiredToken):
			// OTP has expired
			statusCode = http.StatusGone
		case errors.Is(err, errs.ErrRateLimit):
			// Maximum attempts exceeded
			statusCode = http.StatusTooManyRequests
		case errors.Is(err, errs.ErrValueMismatch):
			// Wrong OTP code
			statusCode = http.StatusUnauthorized
		case errors.Is(err, errs.ErrAlreadyConsumed):
			// OTP already consumed by concurrent request
			statusCode = http.StatusConflict
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			// Infrastructure connection issues
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrResourceExhausted):
			// Infrastructure resources exhausted
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrExternalService):
			// External service failure
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			// Query or request cancelled
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, context.DeadlineExceeded):
			// Request timed out
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, errs.ErrTransactionFailure):
			// Temporary infrastructure issues - client should retry
			statusCode = http.StatusServiceUnavailable
		default:
			// Any other error
			statusCode = http.StatusInternalServerError
		}

		logFields = append(logFields, "status_code", statusCode)

		switch logLevel {
		case "info":
			logger.InfoContext(ctx, "Handler: Password reset OTP validation request result", logFields...)
		case "warn":
			logger.WarnContext(ctx, "Handler: Password reset OTP validation request failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Password reset OTP validation request failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
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
