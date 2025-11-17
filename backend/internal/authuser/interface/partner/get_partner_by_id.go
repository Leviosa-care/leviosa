package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) GetPartnerByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract user ID from URL path
	partnerIDStr := r.PathValue("id")
	partnerID, err := uuid.Parse(partnerIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid partner ID format",
			"error", err,
			"operation", "get_partner_by_id",
			"method", r.Method,
			"path", r.URL.Path,
			"user_id", partnerIDStr)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get partner by ID request",
		"operation", "get_partner_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", partnerID,
		"user_agent", r.Header.Get("User-Agent"))

	// partner, err := h.svc.GetPartnerByUserID(ctx, userID)
	partner, err := h.svc.GetPartnerByID(ctx, partnerID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get partner by ID")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get partner by ID completed",
		"operation", "get_partner_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", partnerID,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, partner, http.StatusOK)
}
