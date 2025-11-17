package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"

	"github.com/google/uuid"
)

func (h *handler) DeletePartner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract partner ID from URL path
	partnerIDStr := r.PathValue("id")
	partnerID, err := uuid.Parse(partnerIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid partner ID format",
			"error", err,
			"operation", "delete_partner",
			"method", r.Method,
			"path", r.URL.Path,
			"partner_id", partnerIDStr)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing delete partner request",
		"operation", "delete_partner",
		"method", r.Method,
		"path", r.URL.Path,
		"partner_id", partnerID,
		"user_agent", r.Header.Get("User-Agent"))

	// Call service to delete partner
	err = h.svc.DeletePartner(ctx, partnerID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "delete partner")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Delete partner completed",
		"operation", "delete_partner",
		"method", r.Method,
		"path", r.URL.Path,
		"partner_id", partnerID,
		"status_code", http.StatusNoContent)

	w.WriteHeader(http.StatusNoContent)
}
