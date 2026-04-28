package aggregatorHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) SignOut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get session info from context (injected by auth middleware)
	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		logger.ErrorContext(ctx, "Handler: No session info in context",
			"operation", "sign_out",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("session information from context required"), http.StatusUnauthorized)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing sign-out request",
		"operation", "sign_out",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Call the aggregator service to sign out the user
	err = h.svc.SignOut(ctx, sessionInfo)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "sign out")
		return
	}

	// Clear the session cookies on the client
	cookies.ClearTokenCookies(w)

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Sign-out request completed successfully",
		"operation", "sign_out",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "Successfully signed out",
		Status:  "signed_out",
	}, http.StatusOK)
}
