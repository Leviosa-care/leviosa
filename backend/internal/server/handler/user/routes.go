package userHandler

import (
	"net/http"
)

func (h *AppInstance) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /users/me", h.GetUser)
	router.HandleFunc("PUT /users/me", h.UpdateUser)
	router.HandleFunc("DELETE /users/me", h.DeleteUser)
	router.HandleFunc("POST /users/exists", h.CheckUserExists)
}
