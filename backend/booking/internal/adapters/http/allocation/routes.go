package allocationHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin
	RequirePartner := h.authmw.RequireMinimumRole(identity.Partner)

	// Room allocation management endpoints
	router.HandleFunc("POST /allocations/shared", RequireAdmin(mw.EnableCORS(h.CreateSharedAllocation)))
	router.HandleFunc("POST /allocations/dedicated", RequireAdmin(mw.EnableCORS(h.CreateDedicatedAllocation)))
	router.HandleFunc("GET /allocations/{id}", RequirePartner(mw.EnableCORS(h.GetAllocation)))
	router.HandleFunc("GET /partners/{partnerId}/allocations", RequirePartner(mw.EnableCORS(h.GetPartnerAllocations)))
	router.HandleFunc("GET /rooms/{roomId}/allocations", RequireAdmin(mw.EnableCORS(h.GetRoomAllocations)))
	router.HandleFunc("PUT /allocations/{id}/period", RequireAdmin(mw.EnableCORS(h.UpdateDedicatedPeriod)))
	router.HandleFunc("POST /allocations/{id}/deactivate", RequireAdmin(mw.EnableCORS(h.DeactivateAllocation)))
	router.HandleFunc("GET /partners/{partnerId}/rooms/{roomId}/access", RequirePartner(mw.EnableCORS(h.CheckPartnerRoomAccess)))
}