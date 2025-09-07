package aggregatorHandler

import (
	"fmt"
	"net/http"

	"github.com/Leviosa-care/core/auth/cookies"
	"github.com/Leviosa-care/core/contracts/identity"
	mw "github.com/Leviosa-care/core/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireMinVisitor := h.authmw.RequireMinimumRole(identity.Visitor)
	RequireRefreshToken := h.authmw.RequireRefreshToken

	// Sends an OTP to the provided email address for verification.
	router.HandleFunc("POST /auth/email", mw.EnableCORS(h.CheckEmailSendOTP))

	// Validates an OTP and creates a pending user entry.
	router.HandleFunc("POST /auth/otp", mw.EnableCORS(h.ValidateOTPCreatePendingUser))

	// Completes user registration (e.g., creates Stripe user, sets initial state).
	router.HandleFunc("POST /auth/complete", RequireMinVisitor(mw.EnableCORS(h.CompleteUser)))

	// Refreshes the user session and issues new tokens.
	router.HandleFunc(
		fmt.Sprintf("POST %s", cookies.RefreshEndpoint),
		RequireRefreshToken(mw.EnableCORS(h.RefreshSession)),
	)

	// TODO:
	// ==============================
	// Suggested additional handlers:
	// ==============================

	// Activates a pending user, sets their state to active, and assigns the specified role (admin only).
	// router.HandleFunc("PATCH /admin/auth/users/{id}/activate", mw.EnableCORS(h.ValidateUserRegistration))

	// Logs in a user with email + password (if you support password-based login).
	// router.HandleFunc("POST /auth/login", mw.EnableCORS(h.LoginUser))

	// Logs out the currently authenticated user (e.g., invalidates refresh token).
	// router.HandleFunc("POST /auth/logout", mw.EnableCORS(h.LogoutUser))

	// Initiates password reset flow by sending a reset link/OTP to the user.
	// router.HandleFunc("POST /auth/password/reset/request", mw.EnableCORS(h.RequestPasswordReset))

	// Confirms password reset with a valid token/OTP and updates the password.
	// router.HandleFunc("POST /auth/password/reset/confirm", mw.EnableCORS(h.ConfirmPasswordReset))

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
