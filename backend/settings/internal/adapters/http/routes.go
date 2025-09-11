package http

import (
	"net/http"

	"github.com/Leviosa-care/core/contracts/identity"
	mw "github.com/Leviosa-care/core/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	h.RegisterExternalRoutes(router)
	h.RegisterInternalRoutes(router)
}

func (h *handler) RegisterExternalRoutes(router *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	// Public company settings (GET only)
	router.HandleFunc("GET /settings/name", mw.EnableCORS(h.GetCompanyName))
	router.HandleFunc("GET /settings/email", mw.EnableCORS(h.GetCompanyEmail))
	router.HandleFunc("GET /settings/address", mw.EnableCORS(h.GetCompanyAddress))
	router.HandleFunc("GET /settings/instagram", mw.EnableCORS(h.GetCompanyInstagram))
	router.HandleFunc("GET /settings/logo", mw.EnableCORS(h.GetCompanyLogo))

	// Admin-only company settings (GET sensitive + all POST)
	router.HandleFunc("GET /admin/settings/phone", RequireAdministrator(mw.EnableCORS(h.GetCompanyPhone)))
	router.HandleFunc("POST /admin/settings/name", RequireAdministrator(mw.EnableCORS(h.SetCompanyName)))
	router.HandleFunc("POST /admin/settings/email", RequireAdministrator(mw.EnableCORS(h.SetCompanyEmail)))
	router.HandleFunc("POST /admin/settings/phone", RequireAdministrator(mw.EnableCORS(h.SetCompanyPhone)))
	router.HandleFunc("POST /admin/settings/address", RequireAdministrator(mw.EnableCORS(h.SetCompanyAddress)))
	router.HandleFunc("POST /admin/settings/instagram", RequireAdministrator(mw.EnableCORS(h.SetCompanyInstagram)))
	router.HandleFunc("POST /admin/settings/logo", RequireAdministrator(mw.EnableCORS(h.SetCompanyLogo)))

	// Admin-only OTP settings (all access)
	router.HandleFunc("GET /admin/settings/otp/duration", RequireAdministrator(mw.EnableCORS(h.GetOTPDuration)))
	router.HandleFunc("POST /admin/settings/otp/duration", RequireAdministrator(mw.EnableCORS(h.SetOTPDuration)))
	router.HandleFunc("GET /admin/settings/otp/length", RequireAdministrator(mw.EnableCORS(h.GetOTPLength)))
	router.HandleFunc("POST /admin/settings/otp/length", RequireAdministrator(mw.EnableCORS(h.SetOTPLength)))
	router.HandleFunc("GET /admin/settings/otp/max-attempts", RequireAdministrator(mw.EnableCORS(h.GetOTPMaxAttempts)))
	router.HandleFunc("POST /admin/settings/otp/max-attempts", RequireAdministrator(mw.EnableCORS(h.SetOTPMaxAttempts)))

	// Admin-only token duration settings (all access)
	router.HandleFunc("GET /admin/settings/tokens/access-duration", RequireAdministrator(mw.EnableCORS(h.GetAccessTokenDuration)))
	router.HandleFunc("POST /admin/settings/tokens/access-duration", RequireAdministrator(mw.EnableCORS(h.SetAccessTokenDuration)))
	router.HandleFunc("GET /admin/settings/tokens/refresh-duration", RequireAdministrator(mw.EnableCORS(h.GetRefreshTokenDuration)))
	router.HandleFunc("POST /admin/settings/tokens/refresh-duration", RequireAdministrator(mw.EnableCORS(h.SetRefreshTokenDuration)))

	// Bulk settings endpoint
	router.HandleFunc("GET /settings/bulk", RequireAdministrator(mw.EnableCORS(h.BulkSettingsHandler)))
}

// Internal service-to-service endpoints (protected by service authentication)
// These endpoints allow other microservices to retrieve settings they need
func (h *handler) RegisterInternalRoutes(router *http.ServeMux) {
	RequireService := h.authmw.RequireServiceAuth

	// Company settings accessible to services
	router.HandleFunc("GET /internal/settings/name", RequireService(mw.EnableCORS(h.GetCompanyName)))
	router.HandleFunc("GET /internal/settings/email", RequireService(mw.EnableCORS(h.GetCompanyEmail)))
	router.HandleFunc("GET /internal/settings/address", RequireService(mw.EnableCORS(h.GetCompanyAddress)))
	router.HandleFunc("GET /internal/settings/instagram", RequireService(mw.EnableCORS(h.GetCompanyInstagram)))
	router.HandleFunc("GET /internal/settings/logo", RequireService(mw.EnableCORS(h.GetCompanyLogo)))
	router.HandleFunc("GET /internal/settings/phone", RequireService(mw.EnableCORS(h.GetCompanyPhone)))

	// OTP settings for authentication service
	router.HandleFunc("GET /internal/settings/otp/duration", RequireService(mw.EnableCORS(h.GetOTPDuration)))
	router.HandleFunc("GET /internal/settings/otp/length", RequireService(mw.EnableCORS(h.GetOTPLength)))
	router.HandleFunc("GET /internal/settings/otp/max-attempts", RequireService(mw.EnableCORS(h.GetOTPMaxAttempts)))

	// Token duration settings for authentication service
	router.HandleFunc("GET /internal/settings/tokens/access-duration", RequireService(mw.EnableCORS(h.GetAccessTokenDuration)))
	router.HandleFunc("GET /internal/settings/tokens/refresh-duration", RequireService(mw.EnableCORS(h.GetRefreshTokenDuration)))

	// Bulk settings for services
	router.HandleFunc("GET /internal/settings/bulk", RequireService(mw.EnableCORS(h.BulkSettingsHandler)))
}
