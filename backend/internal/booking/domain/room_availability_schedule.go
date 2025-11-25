package domain

import (
	"time"

	"github.com/google/uuid"
)

// RoomAvailabilitySchedule represents the operating hours for a room
// Supports both recurring patterns (e.g., "every Saturday") and
// specific date exceptions (e.g., "Christmas Day")
type RoomAvailabilitySchedule struct {
	ID uuid.UUID
	RoomID uuid.UUID

	// Recurring pattern: 0=Sunday, 1=Monday, ..., 6=Saturday
	// nil for specific date exceptions
	DayOfWeek *int

	// Specific date override (nil for recurring patterns)
	// Used for one-time exceptions like holidays or special hours
	SpecificDate *time.Time

	// Operating hours for this pattern/date (TIME fields)
	OpenTime  time.Time
	CloseTime time.Time

	// Priority for conflict resolution (higher number = higher priority)
	// Specific dates should have higher priority than recurring patterns
	Priority int

	// Administrative fields
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AppliesTo checks if this schedule applies to a given date
func (s *RoomAvailabilitySchedule) AppliesTo(date time.Time) bool {
	if s.SpecificDate != nil {
		// Specific date schedule - compare dates
		return s.SpecificDate.Truncate(24 * time.Hour).Equal(date.Truncate(24 * time.Hour))
	}

	if s.DayOfWeek != nil {
		// Recurring pattern - compare day of week
		return int(date.Weekday()) == *s.DayOfWeek
	}

	return false
}

// IsRecurring returns true if this is a recurring pattern (not a specific date)
func (s *RoomAvailabilitySchedule) IsRecurring() bool {
	return s.DayOfWeek != nil && s.SpecificDate == nil
}

// IsSpecificDate returns true if this is a specific date override
func (s *RoomAvailabilitySchedule) IsSpecificDate() bool {
	return s.SpecificDate != nil && s.DayOfWeek == nil
}
