package domain

import (
	"time"

	"github.com/google/uuid"
)

// AvailabilityStatus defines the current status of an availability slot
type AvailabilityStatus string

const (
	AvailabilityStatusAvailable AvailabilityStatus = "available"
	AvailabilityStatusBooked    AvailabilityStatus = "booked"
	AvailabilityStatusCancelled AvailabilityStatus = "cancelled"
	AvailabilityStatusBlocked   AvailabilityStatus = "blocked"
)

// RecurrencePattern defines how an availability repeats
type RecurrencePattern struct {
	Type       string     `json:"type"`                   // "daily", "weekly", "monthly"
	Interval   int        `json:"interval"`               // Every N periods
	Until      *time.Time `json:"until,omitempty"`        // End date for recurrence
	DaysOfWeek []int      `json:"days_of_week,omitempty"` // For weekly: 0=Sunday, 1=Monday, etc.
}

// Availability represents a time slot a partner offers for services
type Availability struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	RoomID uuid.UUID `json:"room_id"`

	// Time slot definition
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	// Service offering details
	ServiceType string `json:"service_type,omitempty" encx:"encrypt"`
	PriceCents  *int   `json:"price_cents,omitempty"`
	MaxCapacity int    `json:"max_capacity"`

	// Availability metadata
	Notes             string             `json:"notes,omitempty" encx:"encrypt"`
	IsRecurring       bool               `json:"is_recurring"`
	RecurrencePattern *RecurrencePattern `json:"recurrence_pattern,omitempty"`

	// Status tracking
	Status AvailabilityStatus `json:"status"`

	// Administrative fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewAvailability creates a new availability slot
func NewAvailability(partnerID, roomID uuid.UUID, startTime, endTime time.Time, maxCapacity int) (*Availability, error) {
	if partnerID == uuid.Nil {
		return nil, ErrInvalidPartnerID
	}
	if roomID == uuid.Nil {
		return nil, ErrInvalidRoomID
	}
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return nil, ErrInvalidTimeSlot
	}
	if startTime.Before(time.Now()) {
		return nil, ErrCannotCreatePastAvailability
	}
	if maxCapacity <= 0 {
		return nil, ErrInvalidAvailabilityCapacity
	}

	return &Availability{
		ID:          uuid.New(),
		UserID:      partnerID,
		RoomID:      roomID,
		StartTime:   startTime,
		EndTime:     endTime,
		MaxCapacity: maxCapacity,
		Status:      AvailabilityStatusAvailable,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// NewRecurringAvailability creates a new recurring availability slot
func NewRecurringAvailability(partnerID, roomID uuid.UUID, startTime, endTime time.Time, maxCapacity int, pattern RecurrencePattern) (*Availability, error) {
	availability, err := NewAvailability(partnerID, roomID, startTime, endTime, maxCapacity)
	if err != nil {
		return nil, err
	}

	if err := validateRecurrencePattern(pattern); err != nil {
		return nil, err
	}

	availability.IsRecurring = true
	availability.RecurrencePattern = &pattern
	return availability, nil
}

// UpdateTimeSlot updates the availability's time slot
func (a *Availability) UpdateTimeSlot(startTime, endTime time.Time) error {
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return ErrInvalidTimeSlot
	}
	if startTime.Before(time.Now()) {
		return ErrCannotUpdateToPastTime
	}
	if a.Status == AvailabilityStatusBooked {
		return ErrCannotUpdateBookedAvailability
	}

	a.StartTime = startTime
	a.EndTime = endTime
	a.UpdatedAt = time.Now()
	return nil
}

// SetServiceDetails updates service-specific information
func (a *Availability) SetServiceDetails(serviceType string, priceCents *int, notes string) {
	a.ServiceType = serviceType
	a.PriceCents = priceCents
	a.Notes = notes
	a.UpdatedAt = time.Now()
}

// SetRecurrencePattern sets or updates the recurrence pattern
func (a *Availability) SetRecurrencePattern(pattern *RecurrencePattern) error {
	if pattern != nil {
		if err := validateRecurrencePattern(*pattern); err != nil {
			return err
		}
		a.IsRecurring = true
		a.RecurrencePattern = pattern
	} else {
		a.IsRecurring = false
		a.RecurrencePattern = nil
	}
	a.UpdatedAt = time.Now()
	return nil
}

// MarkAsBooked marks the availability as booked
func (a *Availability) MarkAsBooked() error {
	if a.Status != AvailabilityStatusAvailable {
		return ErrAvailabilityNotAvailable
	}
	a.Status = AvailabilityStatusBooked
	a.UpdatedAt = time.Now()
	return nil
}

// MarkAsAvailable marks the availability as available (e.g., after cancellation)
func (a *Availability) MarkAsAvailable() {
	a.Status = AvailabilityStatusAvailable
	a.UpdatedAt = time.Now()
}

// Cancel cancels the availability
func (a *Availability) Cancel() {
	a.Status = AvailabilityStatusCancelled
	a.UpdatedAt = time.Now()
}

// Block blocks the availability (admin-initiated administrative action)
func (a *Availability) Block() {
	a.Status = AvailabilityStatusBlocked
	a.UpdatedAt = time.Now()
}

// IsAvailable checks if the slot is available for booking
func (a *Availability) IsAvailable() bool {
	return a.Status == AvailabilityStatusAvailable && a.StartTime.After(time.Now())
}

// Duration returns the duration of the availability slot
func (a *Availability) Duration() time.Duration {
	return a.EndTime.Sub(a.StartTime)
}

// OverlapsWith checks if this availability overlaps with another time period
func (a *Availability) OverlapsWith(startTime, endTime time.Time) bool {
	return a.StartTime.Before(endTime) && a.EndTime.After(startTime)
}

// validateRecurrencePattern validates a recurrence pattern
func validateRecurrencePattern(pattern RecurrencePattern) error {
	switch pattern.Type {
	case "daily", "weekly", "monthly":
		// Valid types
	default:
		return ErrInvalidRecurrenceType
	}

	if pattern.Interval <= 0 {
		return ErrInvalidRecurrenceInterval
	}

	if pattern.Type == "weekly" && len(pattern.DaysOfWeek) == 0 {
		return ErrMissingWeeklyDays
	}

	for _, day := range pattern.DaysOfWeek {
		if day < 0 || day > 6 {
			return ErrInvalidWeekDay
		}
	}

	return nil
}
