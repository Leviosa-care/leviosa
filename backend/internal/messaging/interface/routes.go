package messagingHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequirePartner := h.authmw.RequireMinimumRole(identity.Partner)
	RequireStandard := h.authmw.RequireMinimumRole(identity.Standard)

	// Thread endpoints
	router.HandleFunc("GET /threads", RequireStandard(mw.EnableCORS(h.ListThreads)))
	router.HandleFunc("POST /threads", RequirePartner(mw.EnableCORS(h.CreateThread)))
	router.HandleFunc("GET /threads/{id}/messages", RequireStandard(mw.EnableCORS(h.GetMessages)))
	router.HandleFunc("POST /threads/{id}/messages", RequireStandard(mw.EnableCORS(h.SendMessage)))
	router.HandleFunc("POST /threads/{id}/read", RequireStandard(mw.EnableCORS(h.MarkAsRead)))
	router.HandleFunc("GET /threads/unread-count", RequireStandard(mw.EnableCORS(h.GetUnreadCount)))

	// SSE endpoint — no CORS wrapper because SSE uses its own streaming response.
	// CORS headers are set directly in the handler.
	router.HandleFunc("GET /threads/{id}/events", RequireStandard(h.StreamThreadEvents))
}
