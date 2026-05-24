package aggregatorHandler

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	provider := r.PathValue("provider")
	if provider == "" {
		logger.WarnContext(ctx, "Handler: Missing provider in OAuth callback request",
			"operation", "oauth_callback",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("provider parameter is required"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing OAuth callback request",
		"provider", provider,
		"operation", "oauth_callback",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// A state prefixed with "link." means this is an account-linking flow, not a
	// new sign-in.  Extract the user ID embedded in the state and complete the link.
	state := r.URL.Query().Get("state")
	if strings.HasPrefix(state, "link.") {
		h.handleLinkOAuthCallback(w, r, logger, provider, state)
		return
	}

	// Normal sign-in / sign-up flow.
	response, err := h.svc.OAuthCallback(ctx, w, r, provider)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "OAuth callback")
		return
	}

	logger.InfoContext(ctx, "Handler: OAuth callback request completed successfully",
		"provider", provider,
		"is_new_user", response.IsNewUser,
		"operation", "oauth_callback",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusCreated)

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

// handleLinkOAuthCallback completes an OAuth link flow initiated by LinkOAuth.
// State format: "link.{userID}.{base64random}"
func (h *handler) handleLinkOAuthCallback(w http.ResponseWriter, r *http.Request, logger *slog.Logger, provider, state string) {
	ctx := r.Context()

	parts := strings.SplitN(state, ".", 3)
	if len(parts) != 3 {
		logger.WarnContext(ctx, "Handler: Malformed link state in OAuth callback",
			"operation", "link_oauth_callback",
			"provider", provider)
		http.Redirect(w, r, "/staff/profile?error=link_failed", http.StatusFound)
		return
	}

	userID, err := uuid.Parse(parts[1])
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid user ID in link state",
			"operation", "link_oauth_callback",
			"provider", provider,
			"error", err)
		http.Redirect(w, r, "/staff/profile?error=link_failed", http.StatusFound)
		return
	}

	if err := h.svc.CompleteLinkOAuth(ctx, userID, provider, w, r); err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to complete OAuth link",
			"operation", "link_oauth_callback",
			"provider", provider,
			"user_id", userID,
			"error", err)
		http.Redirect(w, r, "/staff/profile?error=link_failed", http.StatusFound)
		return
	}

	logger.InfoContext(ctx, "Handler: OAuth link completed",
		"operation", "link_oauth_callback",
		"provider", provider,
		"user_id", userID)

	http.Redirect(w, r, "/staff/profile?linked="+provider, http.StatusFound)
}
