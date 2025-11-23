package availabilityHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// RequireAdmin := h.authmw.RequireAdmin
	RequirePartner := h.authmw.RequireMinimumRole(identity.Partner)
	RequireStandard := h.authmw.RequireMinimumRole(identity.Standard)

	// Availability management endpoints
	router.HandleFunc("POST /availabilities", RequirePartner(mw.EnableCORS(h.CreateAvailability)))
	router.HandleFunc("POST /availabilities/recurring", RequirePartner(mw.EnableCORS(h.CreateRecurringAvailability)))
	router.HandleFunc("GET /availabilities/{id}", RequireStandard(mw.EnableCORS(h.GetAvailability)))
	router.HandleFunc("GET /partners/{partnerId}/availabilities", RequireStandard(mw.EnableCORS(h.GetPartnerAvailabilities)))
	// router.HandleFunc("GET /availabilities", RequireStandard(mw.EnableCORS(h.GetAvailableSlots)))
	// router.HandleFunc("PUT /availabilities/{id}", RequirePartner(mw.EnableCORS(h.UpdateAvailability)))
	// router.HandleFunc("POST /availabilities/{id}/cancel", RequirePartner(mw.EnableCORS(h.CancelAvailability)))
	// router.HandleFunc("POST /availabilities/{id}/block", RequireAdmin(mw.EnableCORS(h.BlockAvailability)))
	// router.HandleFunc("GET /partners/{partnerId}/availabilities/conflict", RequirePartner(mw.EnableCORS(h.CheckAvailabilityConflict)))
}
