package aggregatorHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) RefreshSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		logger.ErrorContext(ctx, "Handler: No session info in context",
			"operation", "refresh_session",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing session refresh request",
		"operation", "refresh_session",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	response, err := h.svc.RefreshSession(ctx, sessionInfo.ID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "refresh user session")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Session refresh request completed successfully",
		"operation", "refresh_session",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Set new token cookies
	cookies.SetTokenCookies(w, response.AccessToken, response.RefreshToken,
		response.AccessTokenExpiry, response.RefreshTokenExpiry)

	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "Session refreshed successfully",
		Status:  "success",
	}, http.StatusOK)
}
