package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// BookingService defines the interface for booking business logic
type BookingService interface {
	// CreateBooking creates a new booking
	CreateBooking(ctx context.Context, availabilityID, clientID uuid.UUID, clientNotes string) (*domain.Booking, error)

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
}
