package roomHandler

import (
	"net/http"

	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	CreateRoom(w http.ResponseWriter, r *http.Request)
	GetRoom(w http.ResponseWriter, r *http.Request)
	GetRoomsByBuilding(w http.ResponseWriter, r *http.Request)
	GetAllRooms(w http.ResponseWriter, r *http.Request)
	UpdateRoom(w http.ResponseWriter, r *http.Request)
	UpdateRoomPricing(w http.ResponseWriter, r *http.Request)
	ActivateRoom(w http.ResponseWriter, r *http.Request)
	DeactivateRoom(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.RoomService
	authmw auth.AuthMiddleware
}

func New(svc ports.RoomService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}