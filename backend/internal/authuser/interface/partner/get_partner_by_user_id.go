package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetPartnerMe(w http.ResponseWriter, r *http.Request) {
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
			"operation", "get_partner_me",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get partner me request",
		"operation", "get_partner_me",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", sessionInfo.UserID,
		"user_agent", r.Header.Get("User-Agent"))

	partner, err := h.svc.GetPartnerByUserID(ctx, sessionInfo.UserID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get partner me")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get partner me completed",
		"operation", "get_partner_me",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", sessionInfo.UserID,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, partner, http.StatusOK)
}
