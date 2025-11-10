package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// AvailabilityService defines the interface for availability business logic
type AvailabilityService interface {
	// CreateAvailability creates a new availability slot
	CreateAvailability(ctx context.Context, partnerID, roomID uuid.UUID, startTime, endTime time.Time, maxCapacity int) (*domain.Availability, error)

	// CreateRecurringAvailability creates a recurring availability slot
	CreateRecurringAvailability(ctx context.Context, partnerID, roomID uuid.UUID, startTime, endTime time.Time, maxCapacity int, pattern domain.RecurrencePattern) (*domain.Availability, error)

	// GetAvailability retrieves an availability by ID
	GetAvailability(ctx context.Context, id uuid.UUID) (*domain.Availability, error)

	// UpdateAvailability updates availability details
	UpdateAvailability(ctx context.Context, id uuid.UUID, startTime, endTime time.Time, serviceType string, priceCents *int, notes string) (*domain.Availability, error)

	// CancelAvailability cancels an availability slot
	CancelAvailability(ctx context.Context, id uuid.UUID) error

	// BlockAvailability blocks an availability slot
	BlockAvailability(ctx context.Context, id uuid.UUID) error

	// GetPartnerAvailabilities retrieves availabilities for a partner
	GetPartnerAvailabilities(ctx context.Context, partnerID uuid.UUID, filter AvailabilityFilter) ([]*domain.Availability, error)

	// GetAvailableSlots retrieves available slots for booking
	GetAvailableSlots(ctx context.Context, filter AvailabilityFilter) ([]*domain.Availability, error)

	// CheckAvailabilityConflict checks for scheduling conflicts
	CheckAvailabilityConflict(ctx context.Context, partnerID uuid.UUID, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error)
}
