package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin
	RequirePartner := h.authmw.RequireMinimumRole(identity.Partner)

	// NOTE: Partner registration has moved to /auth/complete/partner (self-service during user registration)
	// This admin-only endpoint for creating partners is disabled until we implement a use case for it
	// router.HandleFunc("POST "+CreatePartnerEndpoint, RequireAdmin(mw.EnableCORS(h.CreatePartner)))

	// Get partner by ID (admin only)
	router.HandleFunc("GET "+GetPartnerByIDEndpoint, RequireAdmin(mw.EnableCORS(h.GetPartnerByID)))

	// Get partner by user ID (admin only)
	router.HandleFunc("GET "+GetPartnerByUserIDEndpoint, RequireAdmin(mw.EnableCORS(h.GetPartnerByUserID)))

	// Get all partners (admin only)
	router.HandleFunc("GET "+GetAllPartnersEndpoint, RequireAdmin(mw.EnableCORS(h.GetAllPartners)))

	// Update partner profile (partner can update their own, admin can update any)
	router.HandleFunc("PUT "+UpdatePartnerEndpoint, RequirePartner(mw.EnableCORS(h.UpdatePartner)))

	// Delete partner (admin only)
	router.HandleFunc("DELETE "+DeletePartnerEndpoint, RequireAdmin(mw.EnableCORS(h.DeletePartner)))

	// Verify partner credentials (admin only)
	router.HandleFunc("POST "+VerifyPartnerEndpoint, RequireAdmin(mw.EnableCORS(h.VerifyPartner)))

	// Catalog validation endpoints
	// Validate products exist in catalog (admin only)
	router.HandleFunc("POST "+ValidatePartnerProductsEndpoint, RequireAdmin(mw.EnableCORS(h.ValidatePartnerProducts)))
}
