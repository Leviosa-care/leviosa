package imageHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	UploadImage(w http.ResponseWriter, r *http.Request)
	RemoveImage(w http.ResponseWriter, r *http.Request)
	SetActiveImage(w http.ResponseWriter, r *http.Request)
	// TODO: add the other routes
}

type handler struct {
	svc    ports.ImageCommandService
	authmw auth.AuthMiddleware
}

func New(service ports.ImageCommandService, authmw auth.AuthMiddleware) Handler {
	return &handler{
		svc:    service,
		authmw: authmw,
	}
}
