package aggregatorHandler

import (
	"net/http"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get provider from URL path
	provider := r.PathValue("provider")
	if provider == "" {
		logger.WarnContext(ctx, "Handler: Missing provider in OAuth callback request",
			"operation", "oauth_callback",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("provider parameter is required"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing OAuth callback request",
		"provider", provider,
		"operation", "oauth_callback",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Call service
	response, err := h.svc.OAuthCallback(ctx, w, r, provider)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "OAuth callback")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: OAuth callback request completed successfully",
		"provider", provider,
		"is_new_user", response.IsNewUser,
		"operation", "oauth_callback",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusCreated)

	// Set dual token cookies
	cookies.SetTokenCookies(w, response.AccessToken, response.RefreshToken,
		time.Unix(response.AccessTokenExpiry, 0), time.Unix(response.RefreshTokenExpiry, 0))

	httpx.RespondWithJSON(w, struct {
		Message   string `json:"message"`
		Status    string `json:"status"`
		IsNewUser bool   `json:"is_new_user"`
	}{
		Message:   "OAuth login successful",
		Status:    "created",
		IsNewUser: response.IsNewUser,
	}, http.StatusCreated)
}
