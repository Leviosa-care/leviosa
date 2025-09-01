package http

import (
	"net/http"

	mw "github.com/Leviosa-care/core/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// Public company settings (GET only)
	router.HandleFunc("GET /settings/name", mw.EnableCORS(h.GetCompanyName))
	router.HandleFunc("GET /settings/email", mw.EnableCORS(h.GetCompanyEmail))
	router.HandleFunc("GET /settings/address", mw.EnableCORS(h.GetCompanyAddress))
	router.HandleFunc("GET /settings/instagram", mw.EnableCORS(h.GetCompanyInstagram))
	router.HandleFunc("GET /settings/logo", mw.EnableCORS(h.GetCompanyLogo))

	// Admin-only company settings (GET sensitive + all POST)
	router.HandleFunc("GET /admin/settings/phone", mw.EnableCORS(h.GetCompanyPhone))
	router.HandleFunc("POST /admin/settings/name", mw.EnableCORS(h.SetCompanyName))
	router.HandleFunc("POST /admin/settings/email", mw.EnableCORS(h.SetCompanyEmail))
	router.HandleFunc("POST /admin/settings/phone", mw.EnableCORS(h.SetCompanyPhone))
	router.HandleFunc("POST /admin/settings/address", mw.EnableCORS(h.SetCompanyAddress))
	router.HandleFunc("POST /admin/settings/instagram", mw.EnableCORS(h.SetCompanyInstagram))
	router.HandleFunc("POST /admin/settings/logo", mw.EnableCORS(h.SetCompanyLogo))

	// Admin-only OTP settings (all access)
	router.HandleFunc("GET /admin/settings/otp/duration", mw.EnableCORS(h.GetOTPDuration))
	router.HandleFunc("POST /admin/settings/otp/duration", mw.EnableCORS(h.SetOTPDuration))
	router.HandleFunc("GET /admin/settings/otp/length", mw.EnableCORS(h.GetOTPLength))
	router.HandleFunc("POST /admin/settings/otp/length", mw.EnableCORS(h.SetOTPLength))
	router.HandleFunc("GET /admin/settings/otp/max-attempts", mw.EnableCORS(h.GetOTPMaxAttempts))
	router.HandleFunc("POST /admin/settings/otp/max-attempts", mw.EnableCORS(h.SetOTPMaxAttempts))

	// Admin-only token duration settings (all access)
	router.HandleFunc("GET /admin/settings/tokens/access-duration", mw.EnableCORS(h.GetAccessTokenDuration))
	router.HandleFunc("POST /admin/settings/tokens/access-duration", mw.EnableCORS(h.SetAccessTokenDuration))
	router.HandleFunc("GET /admin/settings/tokens/refresh-duration", mw.EnableCORS(h.GetRefreshTokenDuration))
	router.HandleFunc("POST /admin/settings/tokens/refresh-duration", mw.EnableCORS(h.SetRefreshTokenDuration))

	// Bulk settings endpoint
	router.HandleFunc("GET /settings/bulk", mw.EnableCORS(h.BulkSettingsHandler))
}
