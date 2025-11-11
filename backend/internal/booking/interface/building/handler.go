package buildingHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	CreateBuilding(w http.ResponseWriter, r *http.Request)
	GetBuildingByID(w http.ResponseWriter, r *http.Request)
	GetAllBuildings(w http.ResponseWriter, r *http.Request)
	UpdateBuilding(w http.ResponseWriter, r *http.Request)
	GetBuildingCount(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.BuildingService
	authmw auth.AuthMiddleware
}

func New(svc ports.BuildingService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}
