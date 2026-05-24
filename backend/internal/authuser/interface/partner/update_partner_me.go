package partnerHandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) UpdatePartnerMe(w http.ResponseWriter, r *http.Request) {
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

	// Get session info from context (set by middleware)
	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		logger.ErrorContext(ctx, "Handler: No session info in context",
			"operation", "update_partner_me",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing update partner me request",
		"operation", "update_partner_me",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", sessionInfo.UserID)

	// Get the partner by user ID to find the partner ID
	partner, err := h.svc.GetPartnerByUserID(ctx, sessionInfo.UserID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "update partner me")
		return
	}

	// Parse request body
	var request domain.UpdatePartnerRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "update_partner_me",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid request body"), http.StatusBadRequest)
		return
	}

	// Call service to update partner using the partner's ID
	updatedPartner, err := h.svc.UpdatePartner(ctx, partner.ID, &request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "update partner me")
		return
	}

	logger.InfoContext(ctx, "Handler: Update partner me completed",
		"operation", "update_partner_me",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", sessionInfo.UserID,
		"partner_id", partner.ID,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, updatedPartner, http.StatusOK)
}
