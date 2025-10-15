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
	router.HandleFunc("GET "+GetCompanyNameEndpoint, mw.EnableCORS(h.GetCompanyName))
	router.HandleFunc("GET "+GetCompanyEmailEndpoint, mw.EnableCORS(h.GetCompanyEmail))
	router.HandleFunc("GET "+GetCompanyAddressEndpoint, mw.EnableCORS(h.GetCompanyAddress))
	router.HandleFunc("GET "+GetCompanyInstagramEndpoint, mw.EnableCORS(h.GetCompanyInstagram))
	router.HandleFunc("GET "+GetCompanyLogoEndpoint, mw.EnableCORS(h.GetCompanyLogo))

	// Admin-only company settings (GET sensitive + all POST)
	router.HandleFunc("GET "+AdminGetCompanyPhoneEndpoint, RequireAdministrator(mw.EnableCORS(h.GetCompanyPhone)))
	router.HandleFunc("POST "+SetCompanyNameEndpoint, RequireAdministrator(mw.EnableCORS(h.SetCompanyName)))
	router.HandleFunc("POST "+SetCompanyEmailEndpoint, RequireAdministrator(mw.EnableCORS(h.SetCompanyEmail)))
	router.HandleFunc("POST "+SetCompanyPhoneEndpoint, RequireAdministrator(mw.EnableCORS(h.SetCompanyPhone)))
	router.HandleFunc("POST "+SetCompanyAddressEndpoint, RequireAdministrator(mw.EnableCORS(h.SetCompanyAddress)))
	router.HandleFunc("POST "+SetCompanyInstagramEndpoint, RequireAdministrator(mw.EnableCORS(h.SetCompanyInstagram)))
	router.HandleFunc("POST "+SetCompanyLogoEndpoint, RequireAdministrator(mw.EnableCORS(h.SetCompanyLogo)))

	// Admin-only OTP settings (all access)
	router.HandleFunc("GET "+AdminGetOTPDurationEndpoint, RequireAdministrator(mw.EnableCORS(h.GetOTPDuration)))
	router.HandleFunc("POST "+AdminSetOTPDurationEndpoint, RequireAdministrator(mw.EnableCORS(h.SetOTPDuration)))
	router.HandleFunc("GET "+AdminGetOTPLengthEndpoint, RequireAdministrator(mw.EnableCORS(h.GetOTPLength)))
	router.HandleFunc("POST "+AdminSetOTPLengthEndpoint, RequireAdministrator(mw.EnableCORS(h.SetOTPLength)))
	router.HandleFunc("GET "+AdminGetOTPMaxAttemptsEndpoint, RequireAdministrator(mw.EnableCORS(h.GetOTPMaxAttempts)))
	router.HandleFunc("POST "+AdminSetOTPMaxAttemptsEndpoint, RequireAdministrator(mw.EnableCORS(h.SetOTPMaxAttempts)))

	// Admin-only token duration settings (all access)
	router.HandleFunc("GET "+AdminGetAccessTokenDurationEndpoint, RequireAdministrator(mw.EnableCORS(h.GetAccessTokenDuration)))
	router.HandleFunc("POST "+AdminSetAccessTokenDurationEndpoint, RequireAdministrator(mw.EnableCORS(h.SetAccessTokenDuration)))
	router.HandleFunc("GET "+AdminGetRefreshTokenDurationEndpoint, RequireAdministrator(mw.EnableCORS(h.GetRefreshTokenDuration)))
	router.HandleFunc("POST "+AdminSetRefreshTokenDurationEndpoint, RequireAdministrator(mw.EnableCORS(h.SetRefreshTokenDuration)))

	// Admin-only Bulk settings endpoint
	router.HandleFunc("GET "+AdminBulkEndpoint, RequireAdministrator(mw.EnableCORS(h.BulkSettingsHandler)))
}

// Internal service-to-service endpoints (protected by service authentication)
// These endpoints allow other microservices to retrieve settings they need
func (h *handler) RegisterInternalRoutes(router *http.ServeMux) {
	RequireService := h.authmw.RequireServiceAuth

	// Company settings accessible to services
	router.HandleFunc("GET "+InternalGetCompanyNameEndpoint, RequireService(mw.EnableCORS(h.GetCompanyName)))
	router.HandleFunc("GET "+InternalGetCompanyEmailEndpoint, RequireService(mw.EnableCORS(h.GetCompanyEmail)))
	router.HandleFunc("GET "+InternalGetCompanyAddressEndpoint, RequireService(mw.EnableCORS(h.GetCompanyAddress)))
	router.HandleFunc("GET "+InternalGetCompanyInstagramEndpoint, RequireService(mw.EnableCORS(h.GetCompanyInstagram)))
	router.HandleFunc("GET "+InternalGetCompanyLogoEndpoint, RequireService(mw.EnableCORS(h.GetCompanyLogo)))
	router.HandleFunc("GET "+InternalGetCompanyPhoneEndpoint, RequireService(mw.EnableCORS(h.GetCompanyPhone)))

	// OTP settings for authentication service
	router.HandleFunc("GET "+InternalGetOTPDurationEndpoint, RequireService(mw.EnableCORS(h.GetOTPDuration)))
	router.HandleFunc("GET "+InternalGetOTPLengthEndpoint, RequireService(mw.EnableCORS(h.GetOTPLength)))
	router.HandleFunc("GET "+InternalGetOTPMaxAttemptsEndpoint, RequireService(mw.EnableCORS(h.GetOTPMaxAttempts)))

	// Token duration settings for authentication service
	router.HandleFunc("GET "+InternalGetAccessTokenDurationEndpoint, RequireService(mw.EnableCORS(h.GetAccessTokenDuration)))
	router.HandleFunc("GET "+InternalGetRefreshTokenDurationEndpoint, RequireService(mw.EnableCORS(h.GetRefreshTokenDuration)))

	// Bulk settings for services
	router.HandleFunc("GET "+InternalBulkEndpoint, RequireService(mw.EnableCORS(h.BulkSettingsHandler)))
}
