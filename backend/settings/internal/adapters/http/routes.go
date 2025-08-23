package http

import (
	"net/http"

	"github.com/Leviosa-care/core/httpx"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// Public company settings (GET only)
	router.HandleFunc("GET /settings/name", httpx.EnableCORS(h.GetCompanyName))
	router.HandleFunc("GET /settings/email", httpx.EnableCORS(h.GetCompanyEmail))
	router.HandleFunc("GET /settings/address", httpx.EnableCORS(h.GetCompanyAddress))
	router.HandleFunc("GET /settings/instagram", httpx.EnableCORS(h.GetCompanyInstagram))
	router.HandleFunc("GET /settings/logo", httpx.EnableCORS(h.GetCompanyLogo))

	// Admin-only company settings (GET sensitive + all POST)
	router.HandleFunc("GET /admin/settings/phone", httpx.EnableCORS(h.GetCompanyPhone))
	router.HandleFunc("POST /admin/settings/name", httpx.EnableCORS(h.SetCompanyName))
	router.HandleFunc("POST /admin/settings/email", httpx.EnableCORS(h.SetCompanyEmail))
	router.HandleFunc("POST /admin/settings/phone", httpx.EnableCORS(h.SetCompanyPhone))
	router.HandleFunc("POST /admin/settings/address", httpx.EnableCORS(h.SetCompanyAddress))
	router.HandleFunc("POST /admin/settings/instagram", httpx.EnableCORS(h.SetCompanyInstagram))
	router.HandleFunc("POST /admin/settings/logo", httpx.EnableCORS(h.SetCompanyLogo))

	// Admin-only OTP settings (all access)
	router.HandleFunc("GET /admin/settings/otp/duration", httpx.EnableCORS(h.GetOTPDuration))
	router.HandleFunc("POST /admin/settings/otp/duration", httpx.EnableCORS(h.SetOTPDuration))
	router.HandleFunc("GET /admin/settings/otp/length", httpx.EnableCORS(h.GetOTPLength))
	router.HandleFunc("POST /admin/settings/otp/length", httpx.EnableCORS(h.SetOTPLength))
	router.HandleFunc("GET /admin/settings/otp/max-attempts", httpx.EnableCORS(h.GetOTPMaxAttempts))
	router.HandleFunc("POST /admin/settings/otp/max-attempts", httpx.EnableCORS(h.SetOTPMaxAttempts))

	// Bulk settings endpoint
	router.HandleFunc("GET /settings/bulk", httpx.EnableCORS(h.BulkSettingsHandler))
}
