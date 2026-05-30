package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

type onboardingLinkRequest struct {
	ReturnURL  string `json:"return_url"`
	RefreshURL string `json:"refresh_url"`
}

type onboardingLinkResponse struct {
	URL string `json:"url"`
}

func (h *handler) GetOnboardingLink(w http.ResponseWriter, r *http.Request) {
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
			"operation", "get_onboarding_link",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	req, err := httpx.Decode[onboardingLinkRequest](r.Body)
	if err != nil || req.ReturnURL == "" || req.RefreshURL == "" {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("return_url and refresh_url are required"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get onboarding link request",
		"operation", "get_onboarding_link",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", sessionInfo.UserID)

	url, err := h.svc.GetOnboardingLink(ctx, sessionInfo.UserID, req.ReturnURL, req.RefreshURL)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get onboarding link")
		return
	}

	logger.InfoContext(ctx, "Handler: Get onboarding link completed",
		"operation", "get_onboarding_link",
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", sessionInfo.UserID,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, onboardingLinkResponse{URL: url}, http.StatusOK)
}
