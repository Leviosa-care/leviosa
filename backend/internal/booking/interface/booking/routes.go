package bookingHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	mw "github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdmin := h.authmw.RequireAdmin
	RequirePartner := h.authmw.RequireMinimumRole(identity.Partner)
	RequireStandard := h.authmw.RequireMinimumRole(identity.Standard)

	// Admin endpoints
	router.HandleFunc("GET /admin/dashboard/stats", RequireAdmin(mw.EnableCORS(h.GetDashboardStats)))
	router.HandleFunc("GET /admin/bookings", RequireAdmin(mw.EnableCORS(h.GetAdminBookings)))
	router.HandleFunc("GET /admin/analytics/summary", RequireAdmin(mw.EnableCORS(h.GetAnalyticsSummary)))
	router.HandleFunc("GET /admin/bookings/financial-summary", RequireAdmin(mw.EnableCORS(h.GetFinancialSummary)))

	// Booking management endpoints
	router.HandleFunc("POST /bookings", RequireStandard(mw.EnableCORS(h.CreateBooking)))
	router.HandleFunc("GET /bookings/{id}", RequireStandard(mw.EnableCORS(h.GetBooking)))
	router.HandleFunc("GET /clients/{clientId}/bookings", RequireStandard(mw.EnableCORS(h.GetClientBookings)))
	router.HandleFunc("GET /partners/{partnerId}/bookings", RequirePartner(mw.EnableCORS(h.GetPartnerBookings)))
	router.HandleFunc("GET /partners/{partnerId}/earnings", RequirePartner(mw.EnableCORS(h.GetPartnerEarnings)))
	router.HandleFunc("GET /bookings", RequirePartner(mw.EnableCORS(h.GetUpcomingBookings)))
	router.HandleFunc("PUT /bookings/{id}/notes", RequireStandard(mw.EnableCORS(h.UpdateBookingNotes)))
	router.HandleFunc("POST /bookings/{id}/cancel", RequireStandard(mw.EnableCORS(h.CancelBooking)))
	router.HandleFunc("POST /bookings/{id}/complete", RequirePartner(mw.EnableCORS(h.CompleteBooking)))
	router.HandleFunc("POST /bookings/{id}/no-show", RequirePartner(mw.EnableCORS(h.MarkNoShow)))
	router.HandleFunc("POST /bookings/{id}/payment", RequireStandard(mw.EnableCORS(h.ProcessPayment)))
	router.HandleFunc("POST /bookings/{id}/refund", RequireAdmin(mw.EnableCORS(h.RefundBooking)))

	// Webhook endpoints (no auth required - verified by signature)
	router.HandleFunc("POST /webhooks/stripe", h.HandleStripeWebhook)
}
