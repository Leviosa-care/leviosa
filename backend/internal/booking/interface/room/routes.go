package roomHandler

import (
	"net/http"

	// "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin
	// RequirePartner := h.authmw.RequireMinimumRole(identity.Partner)

	// Room management endpoints
	router.HandleFunc("POST /rooms", RequireAdmin(mw.EnableCORS(h.CreateRoom)))
	router.HandleFunc("GET /rooms/{id}", mw.EnableCORS(h.GetRoom))
	router.HandleFunc("GET /rooms", mw.EnableCORS(h.GetAllRooms))
	router.HandleFunc("PUT /rooms/{id}", RequireAdmin(mw.EnableCORS(h.UpdateRoom)))
	router.HandleFunc("GET /buildings/{buildingId}/rooms", mw.EnableCORS(h.GetRoomsByBuilding))
	// router.HandleFunc("POST /rooms/{id}/activate", RequireAdmin(mw.EnableCORS(h.ActivateRoom)))
	// router.HandleFunc("POST /rooms/{id}/deactivate", RequireAdmin(mw.EnableCORS(h.DeactivateRoom)))
}
