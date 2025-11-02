package categoryHandler

import (
	"encoding/json"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) Healthz(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing health check",
		"operation", "health check",
		"method", r.Method,
		"path", r.URL.Path)

	message := struct {
		Message string `json:"message"`
	}{
		Message: "Hello world",
	}

	logger.InfoContext(ctx, "Handler: health check completed",
		"operation", "health check",
		"status_code", http.StatusOK)

	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
