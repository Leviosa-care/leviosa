package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/domainutil"
	"github.com/hengadev/leviosa/pkg/envmode"

	"github.com/google/uuid"
)

func AttachLogger(env envmode.Mode, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			skipLogging := []string{
				"/healthz",
				// "/hello",
			}
			if slices.Contains(skipLogging, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()

			requestID := uuid.NewString()

			IPHeader := os.Getenv("CLIENT_IP_HEADER")
			// I just make a fake IP for now, I know my function to work:
			var IP string
			switch env {
			case envmode.Dev:
				IP = "127.0.0.1"
			case envmode.Staging, envmode.Prod:
				IP = r.Header.Get(IPHeader)
			}

			if IP == "" {
				logger.ErrorContext(ctx, "client IP not found with required header")
				http.Error(w, "Cannot determine Client IP", http.StatusBadRequest)
				return
			}

			loggingSalt := os.Getenv("LOGGING_SALT")
			if loggingSalt == "" {
				logger.ErrorContext(ctx, "LOGGING_SALT not found in environment variables")
				http.Error(w, "Missing environment variable: LOGGING_SALT", http.StatusInternalServerError)
				return
			}

			hashedIP := domainutil.HashWithSalt(IP, loggingSalt)

			requestLogger := logger.With(
				"method", r.Method,
				"URL", r.URL.String(),
				"IP", hashedIP,
				"requestID", requestID,
			)

			ctx = context.WithValue(r.Context(), ctxutil.LoggerKey, requestLogger)

			requestLogger.InfoContext(ctx, "Request started")
			start := time.Now()
			next.ServeHTTP(w, r.WithContext(ctx))
			duration := time.Since(start)
			requestLogger.InfoContext(ctx, "Request completed", "duration", duration)
		})
	}
}
