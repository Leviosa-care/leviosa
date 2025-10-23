package allocationHandler

import (
	"net/http"

	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	CreateSharedAllocation(w http.ResponseWriter, r *http.Request)
	CreateDedicatedAllocation(w http.ResponseWriter, r *http.Request)
	GetAllocation(w http.ResponseWriter, r *http.Request)
	GetPartnerAllocations(w http.ResponseWriter, r *http.Request)
	GetRoomAllocations(w http.ResponseWriter, r *http.Request)
	UpdateDedicatedPeriod(w http.ResponseWriter, r *http.Request)
	DeactivateAllocation(w http.ResponseWriter, r *http.Request)
	CheckPartnerRoomAccess(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.RoomAllocationService
	authmw auth.AuthMiddleware
}

func New(svc ports.RoomAllocationService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}