package specializationHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	CreateSpecialization(w http.ResponseWriter, r *http.Request)
	GetSpecializationByID(w http.ResponseWriter, r *http.Request)
	GetAllSpecializations(w http.ResponseWriter, r *http.Request)
	UpdateSpecialization(w http.ResponseWriter, r *http.Request)
	DeleteSpecialization(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.SpecializationService
	authmw auth.AuthMiddleware
}

func New(svc ports.SpecializationService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}