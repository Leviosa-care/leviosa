package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	GetPartnerByID(w http.ResponseWriter, r *http.Request)
	GetPartnerMe(w http.ResponseWriter, r *http.Request)
	GetAllPartners(w http.ResponseWriter, r *http.Request)
	GetAllPartnersByCategory(w http.ResponseWriter, r *http.Request)
	UpdatePartner(w http.ResponseWriter, r *http.Request)
	DeletePartner(w http.ResponseWriter, r *http.Request)
	VerifyPartner(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.PartnerService
	authmw auth.AuthMiddleware
}

func New(svc ports.PartnerService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}
