package messagingHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
	"github.com/Leviosa-care/leviosa/backend/internal/messaging/infrastructure/sse"
	"github.com/Leviosa-care/leviosa/backend/internal/messaging/ports"
	"github.com/google/uuid"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	ListThreads(w http.ResponseWriter, r *http.Request)
	CreateThread(w http.ResponseWriter, r *http.Request)
	GetMessages(w http.ResponseWriter, r *http.Request)
	SendMessage(w http.ResponseWriter, r *http.Request)
	MarkAsRead(w http.ResponseWriter, r *http.Request)
	GetUnreadCount(w http.ResponseWriter, r *http.Request)
	StreamThreadEvents(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.MessagingService
	authmw auth.AuthMiddleware
	broker Broker
}

// Broker subscribes clients and publishes SSE events.
type Broker interface {
	Subscribe(threadID uuid.UUID) sse.Subscriber
	Unsubscribe(threadID uuid.UUID, ch sse.Subscriber)
}

func New(svc ports.MessagingService, authmw auth.AuthMiddleware, broker Broker) Handler {
	return &handler{svc: svc, authmw: authmw, broker: broker}
}
