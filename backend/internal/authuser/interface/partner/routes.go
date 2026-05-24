package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin
	RequirePartner := h.authmw.RequireMinimumRole(identity.Partner)

	// Get partner by ID
	router.HandleFunc("GET "+GetPartnerByIDEndpoint, mw.EnableCORS(h.GetPartnerByID))

	// Get authenticated partner's own profile
	router.HandleFunc("GET "+GetPartnerMeEndpoint, RequirePartner(mw.EnableCORS(h.GetPartnerMe)))

	// Get all partners
	router.HandleFunc("GET "+GetAllPartnersEndpoint, RequireAdmin(mw.EnableCORS(h.GetAllPartners)))

	// Get partners by category
	router.HandleFunc("GET "+GetPartnersByCategoryEndpoint, mw.EnableCORS(h.GetAllPartnersByCategory))

	// Get partners by categories
	router.HandleFunc("GET "+GetPartnersByCategoriesEndpoint, mw.EnableCORS(h.GetAllPartnersByCategories))

	// Get partners by product
	router.HandleFunc("GET "+GetPartnersByProductEndpoint, mw.EnableCORS(h.GetAllPartnersByProduct))

	// Get partners by products
	router.HandleFunc("GET "+GetPartnersByProductsEndpoint, mw.EnableCORS(h.GetAllPartnersByProducts))

	// Delete partner (admin only)
	router.HandleFunc("DELETE "+DeletePartnerEndpoint, RequireAdmin(mw.EnableCORS(h.DeletePartner)))

	// Update partner profile (partner can update their own, admin can update any)
	router.HandleFunc("PUT "+UpdatePartnerEndpoint, RequirePartner(mw.EnableCORS(h.UpdatePartner)))

	// Update authenticated partner's own profile
	router.HandleFunc("PUT "+UpdatePartnerMeEndpoint, RequirePartner(mw.EnableCORS(h.UpdatePartnerMe)))

	// Verify partner credentials (admin only)
	router.HandleFunc("POST "+VerifyPartnerEndpoint, RequireAdmin(mw.EnableCORS(h.VerifyPartner)))
}
