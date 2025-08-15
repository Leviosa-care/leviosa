package healthHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/hengadev/leviosa/pkg/ctxutil"
)

func (h *handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	services := make(map[string]string)
	isHealthy := true
	// TODO: do the checking on the databases
	if err := h.Repos.SQLiteDB.PingContext(ctx); err != nil {
		logger.InfoContext(ctx, "health check: SQLite database unhealthy")
		services["database"] = "unhealthy"
		isHealthy = false
	} else {
		logger.InfoContext(ctx, "health check: SQLite database healthy")
		services["database"] = "healthy"
	}
	if _, err := h.Repos.RedisClient.Ping(ctx).Result(); err != nil {
		logger.InfoContext(ctx, "health check: redis database unhealthy")
		services["redis"] = "unhealthy"
		isHealthy = false
	} else {
		logger.InfoContext(ctx, "health check: redis database healthy")
		services["redis"] = "healthy"
	}
	status := "healthy"
	statusCode := http.StatusOK

	if !isHealthy {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}
	response := HealthResponse{
		Status:   status,
		Time:     time.Now().UTC().Format(time.RFC3339),
		Services: services,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
