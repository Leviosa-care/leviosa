package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAllPartners(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get all partners request",
		"operation", "get_all_partners",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	partners, err := h.svc.GetAllPartners(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get all partners")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get all partners completed",
		"operation", "get_all_partners",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"partner_count", len(partners))

	httpx.RespondWithJSON(w, partners, http.StatusOK)
}
