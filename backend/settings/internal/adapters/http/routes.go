package http

import (
	"net/http"

	"github.com/Leviosa-care/core/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// Public company settings (GET only)
	router.HandleFunc("GET /settings/name", middleware.EnableCORS(h.GetCompanyName))
	router.HandleFunc("GET /settings/email", middleware.EnableCORS(h.GetCompanyEmail))
	router.HandleFunc("GET /settings/address", middleware.EnableCORS(h.GetCompanyAddress))
	router.HandleFunc("GET /settings/instagram", middleware.EnableCORS(h.GetCompanyInstagram))
	router.HandleFunc("GET /settings/logo", middleware.EnableCORS(h.GetCompanyLogo))

	// Admin-only company settings (GET sensitive + all POST)
	router.HandleFunc("GET /admin/settings/phone", middleware.EnableCORS(h.GetCompanyPhone))
	router.HandleFunc("POST /admin/settings/name", middleware.EnableCORS(h.SetCompanyName))
	router.HandleFunc("POST /admin/settings/email", middleware.EnableCORS(h.SetCompanyEmail))
	router.HandleFunc("POST /admin/settings/phone", middleware.EnableCORS(h.SetCompanyPhone))
	router.HandleFunc("POST /admin/settings/address", middleware.EnableCORS(h.SetCompanyAddress))
	router.HandleFunc("POST /admin/settings/instagram", middleware.EnableCORS(h.SetCompanyInstagram))
	router.HandleFunc("POST /admin/settings/logo", middleware.EnableCORS(h.SetCompanyLogo))

	// Admin-only OTP settings (all access)
	router.HandleFunc("GET /admin/settings/otp/duration", middleware.EnableCORS(h.GetOTPDuration))
	router.HandleFunc("POST /admin/settings/otp/duration", middleware.EnableCORS(h.SetOTPDuration))
	router.HandleFunc("GET /admin/settings/otp/length", middleware.EnableCORS(h.GetOTPLength))
	router.HandleFunc("POST /admin/settings/otp/length", middleware.EnableCORS(h.SetOTPLength))
	router.HandleFunc("GET /admin/settings/otp/max-attempts", middleware.EnableCORS(h.GetOTPMaxAttempts))
	router.HandleFunc("POST /admin/settings/otp/max-attempts", middleware.EnableCORS(h.SetOTPMaxAttempts))

	// Bulk settings endpoint
	router.HandleFunc("GET /settings/bulk", middleware.EnableCORS(h.BulkSettingsHandler))
}
