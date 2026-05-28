package aggregatorHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

// guestClaimVerifyRequest combines the original claim data with the OTP code.
type guestClaimVerifyRequest struct {
	Email     string `json:"email"`
	Code      string `json:"code"`
	Password  string `json:"password"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Phone     string `json:"phone"`
}

func (h *handler) GuestClaimVerify(w http.ResponseWriter, r *http.Request) {
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

	var payload guestClaimVerifyRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "decode_request",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing guest claim verify request",
		"email", maskEmail(payload.Email),
		"operation", "guest_claim_verify",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Build the domain DTOs from the combined payload
	verifyReq := &domain.GuestClaimVerifyRequest{
		Email: payload.Email,
		Code:  payload.Code,
	}

	claimData := &domain.GuestClaimRequest{
		Email:     payload.Email,
		Phone:     payload.Phone,
		Password:  payload.Password,
		LastName:  payload.LastName,
		FirstName: payload.FirstName,
	}

	session, err := h.svc.GuestClaimVerify(ctx, verifyReq, claimData)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "guest claim verify")
		return
	}

	logger.InfoContext(ctx, "Handler: Guest claim verify completed successfully",
		"email", maskEmail(payload.Email),
		"operation", "guest_claim_verify",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusCreated)

	// Set dual token cookies
	cookies.SetTokenCookies(w, session.AccessToken, session.RefreshToken,
		session.AccessTokenExpiry, session.RefreshTokenExpiry)

	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "Guest account created successfully",
		Status:  "created",
	}, http.StatusCreated)
}
