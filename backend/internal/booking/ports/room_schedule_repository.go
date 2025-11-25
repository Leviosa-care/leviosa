package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// RoomScheduleRepository manages room availability schedules
type RoomScheduleRepository interface {
	// GetRoomHoursForDate returns the operating hours for a room on a specific date
	// Returns the highest priority schedule that applies:
	// - Specific dates take precedence over recurring patterns
	// - Among recurring patterns, higher priority wins
	// Returns ErrRepositoryNotFound if no schedule exists for the date
	GetRoomHoursForDate(ctx context.Context, roomID uuid.UUID, date time.Time) (*domain.RoomAvailabilitySchedule, error)

	// Create adds a new room availability schedule
	Create(ctx context.Context, schedule *domain.RoomAvailabilitySchedule) error

	// GetByID retrieves a schedule by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.RoomAvailabilitySchedule, error)

	// Update modifies an existing schedule
	Update(ctx context.Context, schedule *domain.RoomAvailabilitySchedule) error

	// Delete removes a schedule (soft delete via is_active flag)
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByRoom returns all schedules for a specific room
	// Includes both recurring patterns and specific dates
	// Ordered by priority DESC, specific_date DESC NULLS LAST
	ListByRoom(ctx context.Context, roomID uuid.UUID) ([]*domain.RoomAvailabilitySchedule, error)
}
