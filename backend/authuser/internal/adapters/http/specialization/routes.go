package specializationHandler

import (
	"net/http"

	"github.com/Leviosa-care/core/contracts/identity"
	mw "github.com/Leviosa-care/core/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin
	RequireStandard := h.authmw.RequireMinimumRole(identity.Standard)

	// Admin specialization management
	// Create new specialization (admin only)
	router.HandleFunc("POST /admin/specializations", RequireAdmin(mw.EnableCORS(h.CreateSpecialization)))

	// Get specialization by ID (admin only for detailed view)
	router.HandleFunc("GET /admin/specializations/{id}", RequireAdmin(mw.EnableCORS(h.GetSpecializationByID)))

	// Update specialization (admin only)
	router.HandleFunc("PUT /admin/specializations/{id}", RequireAdmin(mw.EnableCORS(h.UpdateSpecialization)))

	// Delete specialization (admin only)
	router.HandleFunc("DELETE /admin/specializations/{id}", RequireAdmin(mw.EnableCORS(h.DeleteSpecialization)))

	// Public specialization access
	// Get all active specializations (any authenticated user can view for selection)
	router.HandleFunc("GET /specializations", RequireStandard(mw.EnableCORS(h.GetAllSpecializations)))
}

