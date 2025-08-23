package http

import (
	"net/http"

	"github.com/Leviosa-care/settings/internal/ports"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	// Company settings
	GetCompanyName(w http.ResponseWriter, r *http.Request)
	SetCompanyName(w http.ResponseWriter, r *http.Request)
	GetCompanyEmail(w http.ResponseWriter, r *http.Request)
	SetCompanyEmail(w http.ResponseWriter, r *http.Request)
	GetCompanyPhone(w http.ResponseWriter, r *http.Request)
	SetCompanyPhone(w http.ResponseWriter, r *http.Request)
	GetCompanyAddress(w http.ResponseWriter, r *http.Request)
	SetCompanyAddress(w http.ResponseWriter, r *http.Request)
	GetCompanyInstagram(w http.ResponseWriter, r *http.Request)
	SetCompanyInstagram(w http.ResponseWriter, r *http.Request)
	GetCompanyLogo(w http.ResponseWriter, r *http.Request)
	SetCompanyLogo(w http.ResponseWriter, r *http.Request)
	// OTP settings
	GetOTPDuration(w http.ResponseWriter, r *http.Request)
	SetOTPDuration(w http.ResponseWriter, r *http.Request)
	GetOTPLength(w http.ResponseWriter, r *http.Request)
	SetOTPLength(w http.ResponseWriter, r *http.Request)
	GetOTPMaxAttempts(w http.ResponseWriter, r *http.Request)
	SetOTPMaxAttempts(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc ports.SettingsService
}

func New(svc ports.SettingsService) *handler {
	return &handler{svc}
}
