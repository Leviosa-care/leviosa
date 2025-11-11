package buildingHandler

import (
	"net/http"

	// "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin
	// RequirePartner := h.authmw.RequireMinimumRole(identity.Partner)

	// Building management endpoints
	router.HandleFunc("POST /buildings", RequireAdmin(mw.EnableCORS(h.CreateBuilding)))
	router.HandleFunc("GET /buildings/count", mw.EnableCORS(h.GetBuildingCount))
	router.HandleFunc("GET /buildings/{id}", mw.EnableCORS(h.GetBuildingByID))
	router.HandleFunc("GET /buildings", mw.EnableCORS(h.GetAllBuildings))
	router.HandleFunc("PUT /buildings/{id}", RequireAdmin(mw.EnableCORS(h.UpdateBuilding)))
	// router.HandleFunc("PUT /buildings/{id}/contact", RequireAdmin(mw.EnableCORS(h.UpdateBuildingContactInfo)))
	// router.HandleFunc("POST /buildings/{id}/activate", RequireAdmin(mw.EnableCORS(h.ActivateBuilding)))
	// router.HandleFunc("POST /buildings/{id}/deactivate", RequireAdmin(mw.EnableCORS(h.DeactivateBuilding)))
}
