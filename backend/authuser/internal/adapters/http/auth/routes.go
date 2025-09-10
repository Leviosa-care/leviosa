package aggregatorHandler

import (
	"fmt"
	"net/http"

	"github.com/Leviosa-care/core/auth/cookies"
	"github.com/Leviosa-care/core/contracts/identity"
	mw "github.com/Leviosa-care/core/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireVisitor := h.authmw.RequireMinimumRole(identity.Visitor)
	RequireStandard := h.authmw.RequireMinimumRole(identity.Standard)
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)
	RequireRefreshToken := h.authmw.RequireRefreshToken

	// Sends an OTP to the provided email address for verification.
	router.HandleFunc("POST /auth/email", mw.EnableCORS(h.CheckEmailSendOTP))

	// Validates an OTP and creates a pending user entry.
	router.HandleFunc("POST /auth/otp", mw.EnableCORS(h.ValidateOTPCreatePendingUser))

	// Completes user registration (e.g., creates Stripe user, sets initial state).
	router.HandleFunc("POST /auth/complete", RequireVisitor(mw.EnableCORS(h.CompleteUser)))

	// Refreshes the user session and issues new tokens.
	router.HandleFunc(
		fmt.Sprintf("POST %s", cookies.RefreshEndpoint),
		RequireRefreshToken(mw.EnableCORS(h.RefreshSession)),
	)

	// Deletes any user account (admin only).
	router.HandleFunc("DELETE /admin/auth/users/{id}", RequireAdministrator(mw.EnableCORS(h.DeleteUserByAdmin)))

	// Deletes the current user's own account.
	router.HandleFunc("DELETE /auth/me", RequireStandard(mw.EnableCORS(h.DeleteOwnAccount)))

	// Logs in a user with email + password (if you support password-based login).
	router.HandleFunc("POST /auth/login", mw.EnableCORS(h.SignIn))

	// Logs out the currently authenticated user (e.g., invalidates refresh token).
	router.HandleFunc("POST /auth/logout", RequireStandard(mw.EnableCORS(h.SignOut)))

	// Initiates password reset flow by sending a reset OTP to the user.
	router.HandleFunc("POST /auth/password/reset/request", mw.EnableCORS(h.RequestPasswordReset))

	// Validates password reset OTP and issues a reset confirmation token.
	router.HandleFunc("POST /auth/password/reset/validate", mw.EnableCORS(h.ValidatePasswordResetOTP))

	// Confirms password reset with a valid token/OTP and updates the password.
	router.HandleFunc("POST /auth/password/reset/confirm", mw.EnableCORS(h.ConfirmPasswordReset))

	// TODO:
	// ==============================
	// Suggested additional handlers:
	// ==============================

	// ==============================
	// OAuth: Generic provider
	// ==============================

	// Starts the OAuth flow by redirecting to the provider's consent/authorization screen.
	// Example: GET /auth/oauth/google → start Google login
	//          GET /auth/oauth/apple → start Apple login
	// router.HandleFunc("GET /auth/oauth/{provider}", mw.EnableCORS(h.OAuthStart))

	// Handles the provider callback, exchanges code for tokens, and creates/logs in the user.
	// Example: GET /auth/oauth/google/callback → Google callback
	//          GET /auth/oauth/apple/callback → Apple callback
	// router.HandleFunc("GET /auth/oauth/{provider}/callback", mw.EnableCORS(h.OAuthCallback))

	// ==============================
	// Optional: Account linking
	// ==============================

	// Links a provider account to the currently authenticated user.
	// Example: POST /users/me/oauth/google/link
	// router.HandleFunc("POST /users/me/oauth/{provider}/link", mw.EnableCORS(h.LinkOAuthAccount))

	// Unlinks a provider account from the currently authenticated user.
	// Example: DELETE /users/me/oauth/apple/unlink
	// router.HandleFunc("DELETE /users/me/oauth/{provider}/unlink", mw.EnableCORS(h.UnlinkOAuthAccount))
}
