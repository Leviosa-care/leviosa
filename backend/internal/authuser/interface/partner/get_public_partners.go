package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetPublicPartners(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get public partners request",
		"operation", "get_public_partners",
		"method", r.Method,
		"path", r.URL.Path)

	partners, err := h.svc.GetPublicPartners(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get public partners")
		return
	}

	logger.InfoContext(ctx, "Handler: Get public partners completed",
		"operation", "get_public_partners",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"partner_count", len(partners))

	httpx.RespondWithJSON(w, partners, http.StatusOK)
}
