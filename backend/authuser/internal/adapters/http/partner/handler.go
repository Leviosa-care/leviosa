package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/authuser/internal/ports"
	"github.com/Leviosa-care/core/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	CreatePartner(w http.ResponseWriter, r *http.Request)
	GetPartnerByID(w http.ResponseWriter, r *http.Request)
	GetPartnerByUserID(w http.ResponseWriter, r *http.Request)
	GetAllPartners(w http.ResponseWriter, r *http.Request)
	UpdatePartner(w http.ResponseWriter, r *http.Request)
	DeletePartner(w http.ResponseWriter, r *http.Request)
	VerifyPartner(w http.ResponseWriter, r *http.Request)
	AddPartnerSpecialization(w http.ResponseWriter, r *http.Request)
	RemovePartnerSpecialization(w http.ResponseWriter, r *http.Request)
	GetPartnerSpecializations(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.PartnerService
	authmw auth.AuthMiddleware
}

func New(svc ports.PartnerService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}