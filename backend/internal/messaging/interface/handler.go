package messagingHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/messaging/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	ListThreads(w http.ResponseWriter, r *http.Request)
	CreateThread(w http.ResponseWriter, r *http.Request)
	GetMessages(w http.ResponseWriter, r *http.Request)
	SendMessage(w http.ResponseWriter, r *http.Request)
	MarkAsRead(w http.ResponseWriter, r *http.Request)
	GetUnreadCount(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.MessagingService
	authmw auth.AuthMiddleware
}

func New(svc ports.MessagingService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}
