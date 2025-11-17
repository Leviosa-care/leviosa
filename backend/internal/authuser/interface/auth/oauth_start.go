package aggregatorHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) OAuthStart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get provider from URL path
	provider := r.PathValue("provider")
	if provider == "" {
		logger.WarnContext(ctx, "Handler: Missing provider in OAuth start request",
			"operation", "oauth_start",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("provider parameter is required"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing OAuth start request",
		"provider", provider,
		"operation", "oauth_start",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Create request
	request := &domain.OAuthStartRequest{
		Provider: provider,
	}

	// Call service
	response, err := h.svc.OAuthStart(ctx, request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "OAuth start")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: OAuth start request completed successfully",
		"provider", provider,
		"operation", "oauth_start",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusFound,
		"redirect_url", response.AuthorizationURL)

	// Redirect to OAuth provider's authorization URL
	http.Redirect(w, r, response.AuthorizationURL, http.StatusFound)
}
