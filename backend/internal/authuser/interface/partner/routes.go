package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin
	RequirePartner := h.authmw.RequireMinimumRole(identity.Partner)

	// NOTE: done

	// Get partner by ID
	router.HandleFunc("GET "+GetPartnerByIDEndpoint, mw.EnableCORS(h.GetPartnerByID))

	// Get authenticated partner's own profile
	router.HandleFunc("GET "+GetPartnerMeEndpoint, RequirePartner(mw.EnableCORS(h.GetPartnerMe)))

	// Get all partners
	router.HandleFunc("GET "+GetAllPartnersEndpoint, mw.EnableCORS(h.GetAllPartners))

	// Get partners by category
	router.HandleFunc("GET "+GetPartnersByCategoryEndpoint, mw.EnableCORS(h.GetAllPartnersByCategory))

	// Get partners by categories
	router.HandleFunc("GET "+GetPartnersByCategoriesEndpoint, mw.EnableCORS(h.GetAllPartnersByCategories))

	// Delete partner (admin only)
	router.HandleFunc("DELETE "+DeletePartnerEndpoint, RequireAdmin(mw.EnableCORS(h.DeletePartner)))

	// TODO:

	// Update partner profile (partner can update their own, admin can update any)
	router.HandleFunc("PUT "+UpdatePartnerEndpoint, RequirePartner(mw.EnableCORS(h.UpdatePartner)))

	// // Delete partner's own profile
	// router.HandleFunc("DELETE "+DeletePartnerEndpoint, RequirePartner(mw.EnableCORS(h.DeletePartnerMe)))

	// Verify partner credentials (admin only)
	router.HandleFunc("POST "+VerifyPartnerEndpoint, RequireAdmin(mw.EnableCORS(h.VerifyPartner)))
}
