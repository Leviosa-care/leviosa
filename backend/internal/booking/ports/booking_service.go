package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

type EarningsSummary = domain.EarningsSummary
type DashboardStats = domain.DashboardStats
type AnalyticsSummaryResponse = domain.AnalyticsSummaryResponse

// BookingService defines the interface for booking business logic
type BookingService interface {
	// CreateBooking creates a new booking with product and time slot information
	CreateBooking(ctx context.Context, availabilityID, clientID, productID uuid.UUID, slotStartTime time.Time, clientNotes string) (*domain.Booking, error)

	// GetBooking retrieves a booking by ID
	GetBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error)

	// UpdateBookingNotes updates client or partner notes
	UpdateBookingNotes(ctx context.Context, id uuid.UUID, clientNotes, partnerNotes string) (*domain.Booking, error)

	// CancelBooking cancels a booking with reason
	CancelBooking(ctx context.Context, id uuid.UUID, reason string) (*domain.Booking, error)

	// CompleteBooking marks a booking as completed
	CompleteBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error)

	// MarkNoShow marks a booking as no-show
	MarkNoShow(ctx context.Context, id uuid.UUID) (*domain.Booking, error)

	// ProcessPayment handles payment processing
	ProcessPayment(ctx context.Context, id uuid.UUID, paymentIntentID string) (*domain.Booking, error)

	// RefundBooking processes a refund
	RefundBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error)

	// GetClientBookings retrieves bookings for a client
	GetClientBookings(ctx context.Context, clientID uuid.UUID, filter BookingFilter) ([]*domain.Booking, error)

	// GetPartnerBookings retrieves bookings for a partner
	GetPartnerBookings(ctx context.Context, partnerID uuid.UUID, filter BookingFilter) ([]*domain.Booking, error)

	// GetUpcomingBookings retrieves upcoming confirmed bookings
	GetUpcomingBookings(ctx context.Context, filter BookingFilter) ([]*domain.Booking, error)

	// HandlePaymentWebhook processes a Stripe payment webhook event and updates the booking status
	HandlePaymentWebhook(ctx context.Context, event *WebhookEvent) error

	// GetPartnerEarnings retrieves earnings summary and transactions for a partner
	GetPartnerEarnings(ctx context.Context, partnerID uuid.UUID) (*EarningsSummary, error)

	// GetDashboardStats retrieves aggregated dashboard statistics for admins
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)

	// GetAdminBookings retrieves a paginated list of all bookings with enriched
	// fields (client name, partner name, product name, room name) for the admin UI.
	GetAdminBookings(ctx context.Context, filter AdminBookingsFilter) (*domain.AdminBookingsListResponse, error)

	// GetAnalyticsSummary computes aggregated analytics for the admin dashboard:
	// current-month KPIs, monthly revenue time-series, and top products.
	GetAnalyticsSummary(ctx context.Context, months int) (*AnalyticsSummaryResponse, error)

	// GetFinancialSummary computes gross revenue, refunds, net revenue, and returns
	// all paid/refunded transactions for the given date range.
	GetFinancialSummary(ctx context.Context, from, to time.Time) (*domain.FinancialSummaryResponse, error)
}

// AdminBookingsFilter holds query parameters for the admin bookings list endpoint.
type AdminBookingsFilter struct {
	Status    *domain.BookingStatus
	PartnerID *uuid.UUID
	From      *time.Time
	To        *time.Time
	Page      int
	Limit     int
}
