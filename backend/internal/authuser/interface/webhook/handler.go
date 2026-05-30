package webhookHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

// Handler defines the webhook handler interface.
type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	HandleStripeConnectWebhook(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.PartnerService
	stripe ports.StripeService
	authmw auth.AuthMiddleware
}

// New creates a new webhook handler.
func New(svc ports.PartnerService, stripe ports.StripeService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, stripe: stripe, authmw: authmw}
}
