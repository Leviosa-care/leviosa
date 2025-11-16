package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func TestClassifyServiceError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		// Client Errors (4xx)
		{"invalid value", errs.ErrInvalidValue, http.StatusBadRequest},
		{"invalid input", errs.ErrInvalidInput, http.StatusBadRequest},
		{"validation", errs.ErrValidation, http.StatusBadRequest},
		{"no fields for update", errs.ErrNoFieldsForUpdate, http.StatusBadRequest},

		{"unauthorized", errs.ErrUnauthorized, http.StatusUnauthorized},
		{"expired token", errs.ErrExpiredToken, http.StatusUnauthorized},

		{"forbidden", errs.ErrForbidden, http.StatusForbidden},
		{"permission denied", errs.ErrPermissionDenied, http.StatusForbidden},

		{"domain not found", errs.ErrDomainNotFound, http.StatusNotFound},
		{"repository not found", errs.ErrRepositoryNotFound, http.StatusNotFound},

		{"context canceled", context.Canceled, http.StatusRequestTimeout},
		{"query cancelled", errs.ErrQueryCancelled, http.StatusRequestTimeout},

		{"already exists", errs.ErrAlreadyExists, http.StatusConflict},
		{"conflict", errs.ErrConflict, http.StatusConflict},
		{"unique violation", errs.ErrUniqueViolation, http.StatusConflict},

		{"account locked", errs.ErrAccountLocked, http.StatusLocked},

		{"rate limit", errs.ErrRateLimit, http.StatusTooManyRequests},

		// Server Errors (5xx)
		{"internal", errs.ErrInternal, http.StatusInternalServerError},
		{"unexpected error", errs.ErrUnexpectedError, http.StatusInternalServerError},
		{"database", errs.ErrDatabase, http.StatusInternalServerError},
		{"db query", errs.ErrDBQuery, http.StatusInternalServerError},
		{"query failed", errs.ErrQueryFailed, http.StatusInternalServerError},
		{"not encrypted", errs.ErrNotEncrypted, http.StatusInternalServerError},
		{"not decrypted", errs.ErrNotDecrypted, http.StatusInternalServerError},
		{"domain not created", errs.ErrDomainNotCreated, http.StatusInternalServerError},
		{"domain not updated", errs.ErrDomainNotUpdated, http.StatusInternalServerError},
		{"domain not deleted", errs.ErrDomainNotDeleted, http.StatusInternalServerError},
		{"repository not created", errs.ErrRepositoryNotCreated, http.StatusInternalServerError},
		{"repository not updated", errs.ErrRepositoryNotUpdated, http.StatusInternalServerError},
		{"repository not deleted", errs.ErrRepositoryNotDeleted, http.StatusInternalServerError},

		{"external service", errs.ErrExternalService, http.StatusBadGateway},

		{"connection failure", errs.ErrConnectionFailure, http.StatusServiceUnavailable},
		{"too many connections", errs.ErrTooManyConnections, http.StatusServiceUnavailable},
		{"resource exhausted", errs.ErrResourceExhausted, http.StatusServiceUnavailable},
		{"transaction failure", errs.ErrTransactionFailure, http.StatusServiceUnavailable},
		{"deadlock", errs.ErrDeadlock, http.StatusServiceUnavailable},

		{"context deadline exceeded", context.DeadlineExceeded, http.StatusGatewayTimeout},

		// Default case
		{"unknown error", errors.New("some unknown error"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := classifyServiceError(tt.err)
			if status != tt.expectedStatus {
				t.Errorf("classifyServiceError() = %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestClassifyServiceError_Wrapped(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{"wrapped invalid value", fmt.Errorf("operation failed: %w", errs.ErrInvalidValue), http.StatusBadRequest},
		{"wrapped not found", fmt.Errorf("user not found: %w", errs.ErrRepositoryNotFound), http.StatusNotFound},
		{"wrapped connection failure", fmt.Errorf("database issue: %w", errs.ErrConnectionFailure), http.StatusServiceUnavailable},
		{"multiple wraps", fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", errs.ErrInvalidValue)), http.StatusBadRequest},
		{"wrapped context cancelled", fmt.Errorf("request cancelled: %w", context.Canceled), http.StatusRequestTimeout},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := classifyServiceError(tt.err)
			if status != tt.expectedStatus {
				t.Errorf("classifyServiceError() = %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestRespondWithServiceError_ClientError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		shouldLog      bool
	}{
		{
			name:           "bad request error",
			err:            errs.ErrInvalidValue,
			expectedStatus: http.StatusBadRequest,
			shouldLog:      false,
		},
		{
			name:           "not found error",
			err:            errs.ErrRepositoryNotFound,
			expectedStatus: http.StatusNotFound,
			shouldLog:      false,
		},
		{
			name:           "conflict error",
			err:            errs.ErrAlreadyExists,
			expectedStatus: http.StatusConflict,
			shouldLog:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture logs
			var logBuffer bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{}))

			w := httptest.NewRecorder()
			ctx := context.Background()

			RespondWithServiceError(w, logger, ctx, tt.err, "test operation")

			resp := w.Result()
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Check content type
			expectedContentType := "application/json"
			if ct := resp.Header.Get("Content-Type"); ct != expectedContentType {
				t.Errorf("Expected Content-Type %s, got %s", expectedContentType, ct)
			}

			// Check response body contains original error message
			var response map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response["error"] != tt.err.Error() {
				t.Errorf("Expected error message %q, got %q", tt.err.Error(), response["error"])
			}

			// Check that no logging occurred for client errors
			if logBuffer.Len() > 0 && tt.shouldLog == false {
				t.Errorf("Expected no logging for client error, but got: %s", logBuffer.String())
			}
		})
	}
}

func TestRespondWithServiceError_ServerError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		shouldLog      bool
	}{
		{
			name:           "internal error",
			err:            errs.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			shouldLog:      true,
		},
		{
			name:           "database connection error",
			err:            errs.ErrConnectionFailure,
			expectedStatus: http.StatusServiceUnavailable,
			shouldLog:      true,
		},
		{
			name:           "external service error",
			err:            errs.ErrExternalService,
			expectedStatus: http.StatusBadGateway,
			shouldLog:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture logs
			var logBuffer bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{}))

			w := httptest.NewRecorder()
			ctx := context.Background()

			RespondWithServiceError(w, logger, ctx, tt.err, "create user")

			resp := w.Result()
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Check content type
			expectedContentType := "application/json"
			if ct := resp.Header.Get("Content-Type"); ct != expectedContentType {
				t.Errorf("Expected Content-Type %s, got %s", expectedContentType, ct)
			}

			// Check response body contains generic error message (not internal details)
			var response map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			expectedMessage := "an internal server error occurred"
			if response["error"] != expectedMessage {
				t.Errorf("Expected error message %q, got %q", expectedMessage, response["error"])
			}

			// Check that logging occurred for server errors
			if tt.shouldLog && logBuffer.Len() == 0 {
				t.Errorf("Expected logging for server error, but no logs captured")
			}

			if tt.shouldLog {
				var logEntry map[string]interface{}
				if err := json.Unmarshal(logBuffer.Bytes(), &logEntry); err != nil {
					t.Fatalf("Failed to decode log entry: %v", err)
				}

				// Check that log contains operation information
				if msg, ok := logEntry["msg"].(string); !ok || !contains(msg, "create user") {
					t.Errorf("Expected log message to contain operation name, got: %v", msg)
				}

				// Check that log contains error
				if _, ok := logEntry["error"]; !ok {
					t.Errorf("Expected log entry to contain error field")
				}

				// Check that log contains status code
				if status, ok := logEntry["status_code"]; !ok || int(status.(float64)) != tt.expectedStatus {
					t.Errorf("Expected status code %d in log, got: %v", tt.expectedStatus, status)
				}
			}
		})
	}
}

func TestRespondWithServiceError_NilError(t *testing.T) {
	w := httptest.NewRecorder()
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))
	ctx := context.Background()

	RespondWithServiceError(w, logger, ctx, nil, "test operation")

	resp := w.Result()
	defer resp.Body.Close()

	// Should not write any response for nil error
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for nil error, got %d", resp.StatusCode)
	}

	// Response body should be empty
	body := make([]byte, 1024)
	n, _ := resp.Body.Read(body)
	if n > 0 {
		t.Errorf("Expected empty response body for nil error, got: %s", string(body[:n]))
	}
}

func TestRespondWithServiceError_ContextErrors(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		shouldLog      bool
	}{
		{
			name:           "context canceled",
			err:            context.Canceled,
			expectedStatus: http.StatusRequestTimeout,
			shouldLog:      false, // 408 is a client error
		},
		{
			name:           "context deadline exceeded",
			err:            context.DeadlineExceeded,
			expectedStatus: http.StatusGatewayTimeout,
			shouldLog:      true, // 504 is a server error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuffer bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{}))

			w := httptest.NewRecorder()
			ctx := context.Background()

			RespondWithServiceError(w, logger, ctx, tt.err, "timeout operation")

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Check logging behavior
			hasLogs := logBuffer.Len() > 0
			if hasLogs != tt.shouldLog {
				t.Errorf("Expected logging=%v, but got logs=%v", tt.shouldLog, hasLogs)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}