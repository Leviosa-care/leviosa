package webhookHandler

import "net/http"

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// Webhook endpoints (no auth required — verified by Stripe signature)
	router.HandleFunc("POST "+HandleStripeConnectWebhookEndpoint, h.HandleStripeConnectWebhook)
}
