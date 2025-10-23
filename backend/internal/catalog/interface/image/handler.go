package imageHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	UploadImage(w http.ResponseWriter, r *http.Request)
	RemoveImage(w http.ResponseWriter, r *http.Request)
	SetActiveImage(w http.ResponseWriter, r *http.Request)
	// TODO: add the other routes
}

type handler struct {
	svc ports.ImageCommandService
}

func New(service ports.ImageCommandService) Handler {
	return &handler{
		svc: service,
	}
}
