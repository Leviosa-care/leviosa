package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) VerifyPartner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get session info from context (set by middleware)
	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		logger.ErrorContext(ctx, "Handler: No session info in context",
			"operation", "verify_partner",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// Extract partner ID from URL path
	partnerIDStr := r.PathValue("id")
	partnerID, err := uuid.Parse(partnerIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid partner ID format",
			"error", err,
			"operation", "verify_partner",
			"method", r.Method,
			"path", r.URL.Path,
			"partner_id", partnerIDStr)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing verify partner request",
		"operation", "verify_partner",
		"method", r.Method,
		"path", r.URL.Path,
		"partner_id", partnerID,
		"verified_by_user_id", sessionInfo.UserID,
		"user_agent", r.Header.Get("User-Agent"))

	// Call service to verify partner
	partner, err := h.svc.VerifyPartner(ctx, partnerID, sessionInfo.UserID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "verify partner")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Verify partner completed",
		"operation", "verify_partner",
		"method", r.Method,
		"path", r.URL.Path,
		"partner_id", partnerID,
		"verified_by_user_id", sessionInfo.UserID,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, partner, http.StatusOK)
}
