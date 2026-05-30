package webhookHandler

import (
	"net/http"

	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// Webhook endpoints (no auth required — verified by Stripe signature)
	router.HandleFunc("POST "+HandleStripeConnectWebhookEndpoint, mw.EnableCORS(h.HandleStripeConnectWebhook))
}
