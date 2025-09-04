package userHandler

import (
	"net/http"

	"github.com/Leviosa-care/authuser/internal/ports"
	"github.com/Leviosa-care/core/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	GetPendingUsers(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	ApproveUser(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.UserService
	authmw auth.AuthMiddleware
}

func New(svc ports.UserService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}
