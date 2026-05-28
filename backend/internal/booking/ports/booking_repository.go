package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// BookingRepository defines the interface for booking data persistence
type BookingRepository interface {
	// Create stores a new booking
	Create(ctx context.Context, booking *domain.BookingEncx) error

	// GetByID retrieves a booking by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BookingEncx, error)

	// Update modifies an existing booking
	Update(ctx context.Context, booking *domain.BookingEncx) error

	// Delete removes a booking (hard delete for GDPR compliance)
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves bookings with optional filtering
	List(ctx context.Context, filter BookingFilter) ([]*domain.BookingEncx, error)

	// GetByClientID retrieves all bookings for a specific client
	GetByClientID(ctx context.Context, clientID uuid.UUID, filter BookingFilter) ([]*domain.BookingEncx, error)

	// GetByPartnerID retrieves all bookings for a specific partner
	GetByPartnerID(ctx context.Context, partnerID uuid.UUID, filter BookingFilter) ([]*domain.BookingEncx, error)

	// GetByAvailabilityID retrieves booking for a specific availability (should be max 1)
	// DEPRECATED: Use GetBookingsByAvailability for slot-based booking system
	GetByAvailabilityID(ctx context.Context, availabilityID uuid.UUID) (*domain.BookingEncx, error)

	// GetBookingsByAvailability retrieves all bookings for a specific availability.
	// Used for slot overlap detection in the slot-based booking system.
	// Returns bookings sorted by slot_start_time ASC.
	GetBookingsByAvailability(ctx context.Context, availabilityID uuid.UUID) ([]*domain.BookingEncx, error)

	// GetUpcoming retrieves upcoming bookings (status = confirmed, start time in future)
	GetUpcoming(ctx context.Context, filter BookingFilter) ([]*domain.BookingEncx, error)

	// GetByPaymentIntentID retrieves booking by Stripe payment intent ID
	GetByPaymentIntentID(ctx context.Context, paymentIntentID string) (*domain.BookingEncx, error)

	// GetGuestBookings retrieves all bookings that have no client_id (guest bookings).
	GetGuestBookings(ctx context.Context) ([]*domain.BookingEncx, error)

	// SetBookingClientID atomically sets client_id on a single booking.
	// Only succeeds if the booking's client_id is currently NULL.
	// Returns true if the row was updated, false if it was already owned.
	SetBookingClientID(ctx context.Context, bookingID, clientID uuid.UUID) (bool, error)
}

// BookingFilter defines filtering options for booking queries
type BookingFilter struct {
	// Entity filters
	ClientID       *uuid.UUID
	PartnerID      *uuid.UUID
	RoomID         *uuid.UUID
	AvailabilityID *uuid.UUID

	// Status filters
	Status        []domain.BookingStatus
	PaymentStatus []domain.PaymentStatus

	// Time-based filters
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	TimeRange     *TimeRange // Filter by availability time range

	// Price filters
	MinPrice *int // In cents
	MaxPrice *int // In cents

	// Pagination
	Limit  int
	Offset int

	// Sorting
	OrderBy        string // "created_at", "total_price_cents"
	OrderDirection string // "asc", "desc"
}
