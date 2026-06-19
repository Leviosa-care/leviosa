package metricsHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

// RegisterRoutes registers metrics HTTP routes
func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	RequirePartner := h.authmw.RequireMinimumRole(identity.Partner)

	// Metrics endpoints
	router.HandleFunc("GET /rooms/{room_id}/metrics", RequirePartner(mw.EnableCORS(h.GetRoomMetrics)))
	router.HandleFunc("GET /partners/metrics/{partner_id}", RequirePartner(mw.EnableCORS(h.GetPartnerMetrics)))
}
