package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/core/contracts/identity"
	mw "github.com/Leviosa-care/core/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin
	RequireStandard := h.authmw.RequireMinimumRole(identity.Standard)
	RequirePartner := h.authmw.RequireMinimumRole(identity.Partner)

	// Partner registration endpoint - allows creating partner users
	router.HandleFunc("POST "+CreatePartnerEndpoint, RequireAdmin(mw.EnableCORS(h.CreatePartner)))

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

	// Partner specialization management
	// Add specialization to partner (admin only)
	router.HandleFunc("POST "+AddPartnerSpecializationEndpoint, RequireAdmin(mw.EnableCORS(h.AddPartnerSpecialization)))

	// Remove specialization from partner (admin only)
	router.HandleFunc("DELETE "+RemovePartnerSpecializationEndpoint, RequireAdmin(mw.EnableCORS(h.RemovePartnerSpecialization)))

	// Get partner specializations (any authenticated user can view)
	router.HandleFunc("GET "+GetPartnerSpecializationsEndpoint, RequireStandard(mw.EnableCORS(h.GetPartnerSpecializations)))
}

