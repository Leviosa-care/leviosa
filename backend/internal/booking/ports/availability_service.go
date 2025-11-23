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
	CreateAvailability(ctx context.Context, request *domain.CreateAvailabilityRequest) (*domain.Availability, error)

	// CreateRecurringAvailability creates a recurring availability slot
	CreateRecurringAvailability(ctx context.Context, request *domain.CreateRecurringAvailabilityRequest) (*domain.Availability, error)

	// GetAvailability retrieves an availability by ID
	GetAvailability(ctx context.Context, id uuid.UUID) (*domain.Availability, error)

	// UpdateAvailability updates availability details
	UpdateAvailability(ctx context.Context, id uuid.UUID, request *domain.UpdateAvailabilityRequest) (*domain.Availability, error)

	// CancelAvailability cancels an availability slot
	CancelAvailability(ctx context.Context, id uuid.UUID) error

	// BlockAvailability blocks an availability slot
	BlockAvailability(ctx context.Context, id uuid.UUID) error

	// GetPartnerAvailabilities retrieves availabilities for a Partner
	GetPartnerAvailabilities(ctx context.Context, userID uuid.UUID, filter AvailabilityFilter) ([]*domain.Availability, error)

	// GetAvailableSlots retrieves available slots for booking
	GetAvailableSlots(ctx context.Context, filter AvailabilityFilter) ([]*domain.Availability, error)

	// CheckAvailabilityConflict checks for scheduling conflicts
	CheckAvailabilityConflict(ctx context.Context, userID uuid.UUID, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error)
}
