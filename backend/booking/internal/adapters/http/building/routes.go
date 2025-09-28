package buildingHandler

import (
	"net/http"

	"github.com/Leviosa-care/core/contracts/identity"
	mw "github.com/Leviosa-care/core/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin
	RequirePartner := h.authmw.RequireMinimumRole(identity.Partner)

	// Building management endpoints
	router.HandleFunc("POST /buildings", RequireAdmin(mw.EnableCORS(h.CreateBuilding)))
	router.HandleFunc("GET /buildings/{id}", RequirePartner(mw.EnableCORS(h.GetBuilding)))
	router.HandleFunc("GET /buildings", RequirePartner(mw.EnableCORS(h.GetAllBuildings)))
	router.HandleFunc("PUT /buildings/{id}", RequireAdmin(mw.EnableCORS(h.UpdateBuilding)))
	router.HandleFunc("PUT /buildings/{id}/contact", RequireAdmin(mw.EnableCORS(h.UpdateBuildingContactInfo)))
	router.HandleFunc("POST /buildings/{id}/activate", RequireAdmin(mw.EnableCORS(h.ActivateBuilding)))
	router.HandleFunc("POST /buildings/{id}/deactivate", RequireAdmin(mw.EnableCORS(h.DeactivateBuilding)))
}