package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type AvailabilityResponse struct {
	ID              uuid.UUID          `json:"id"`
	UserID          uuid.UUID          `json:"partner_id"`
	RoomID          uuid.UUID          `json:"room_id"`
	StartTime       time.Time          `json:"start_time"`
	EndTime         time.Time          `json:"end_time"`
	MaxCapacity     int                `json:"max_capacity"`
	CurrentBookings int                `json:"current_bookings"`
	Status          AvailabilityStatus `json:"status"`
	ServiceType     string             `json:"service_type,omitempty"`
	PriceCents      *int               `json:"price_cents,omitempty"`
	Notes           string             `json:"notes,omitempty"`
	RecurrenceRule  *RecurrencePattern `json:"recurrence_rule,omitempty"`
	ParentID        *uuid.UUID         `json:"parent_id,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

// Availability DTOs
type CreateAvailabilityRequest struct {
	UserID      uuid.UUID `json:"partner_id" validate:"required"`
	RoomID      uuid.UUID `json:"room_id" validate:"required"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
	MaxCapacity int       `json:"max_capacity" validate:"required,min=1,max=50"`
	ServiceType string    `json:"service_type,omitempty" validate:"max=255"`
	PriceCents  *int      `json:"price_cents,omitempty" validate:"omitempty,min=0"`
	Notes       string    `json:"notes,omitempty" validate:"max=1000"`
}

func (r *CreateAvailabilityRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate UUID fields
	if r.UserID == uuid.Nil {
		errs.Set("user_id", "user ID is required")
	}

	if r.RoomID == uuid.Nil {
		errs.Set("room_id", "room ID is required")
	}

	// Validate time constraints
	if r.StartTime.IsZero() {
		errs.Set("start_time", "start time is required")
	}

	if r.EndTime.IsZero() {
		errs.Set("end_time", "end time is required")
	}

	if !r.StartTime.IsZero() && !r.EndTime.IsZero() {
		if r.StartTime.Before(time.Now()) {
			errs.Set("start_time", "start time cannot be in the past")
		}

		if r.StartTime.After(r.EndTime) {
			errs.Set("start_time", "start time must be before end time")
		}

		if r.StartTime.Equal(r.EndTime) {
			errs.Set("start_time", "start time and end time cannot be the same")
		}

		// Check for reasonable duration (between 15 minutes and 12 hours)
		duration := r.EndTime.Sub(r.StartTime)
		if duration < 15*time.Minute {
			errs.Set("end_time", "availability duration must be at least 15 minutes")
		}
		if duration > 12*time.Hour {
			errs.Set("end_time", "availability duration cannot exceed 12 hours")
		}
	}

	// Validate capacity constraints
	if r.MaxCapacity <= 0 {
		errs.Set("max_capacity", "max capacity must be greater than 0")
	}

	if r.MaxCapacity > 50 {
		errs.Set("max_capacity", "max capacity cannot exceed 50")
	}

	// Validate optional fields
	if len(r.ServiceType) > 255 {
		errs.Set("service_type", "service type cannot exceed 255 characters")
	}

	if len(r.Notes) > 1000 {
		errs.Set("notes", "notes cannot exceed 1000 characters")
	}

	// Validate price if provided
	if r.PriceCents != nil {
		if *r.PriceCents < 0 {
			errs.Set("price_cents", "price cannot be negative")
		}
		if *r.PriceCents > 999999 { // $9,999.99 max
			errs.Set("price_cents", "price cannot exceed $9,999.99")
		}
	}

	return errs.AsError()
}

type CreateRecurringAvailabilityRequest struct {
	UserID      uuid.UUID         `json:"partner_id" validate:"required"`
	RoomID      uuid.UUID         `json:"room_id" validate:"required"`
	StartTime   time.Time         `json:"start_time" validate:"required"`
	EndTime     time.Time         `json:"end_time" validate:"required"`
	MaxCapacity int               `json:"max_capacity" validate:"required,min=1,max=50"`
	Pattern     RecurrencePattern `json:"pattern" validate:"required"`
	ServiceType string            `json:"service_type,omitempty" validate:"max=255"`
	PriceCents  *int              `json:"price_cents,omitempty" validate:"omitempty,min=0"`
	Notes       string            `json:"notes,omitempty" validate:"max=1000"`
}

func (r *CreateRecurringAvailabilityRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate UUID fields
	if r.UserID == uuid.Nil {
		errs.Set("user_id", "user ID is required")
	}

	if r.RoomID == uuid.Nil {
		errs.Set("room_id", "room ID is required")
	}

	// Validate time constraints
	if r.StartTime.IsZero() {
		errs.Set("start_time", "start time is required")
	}

	if r.EndTime.IsZero() {
		errs.Set("end_time", "end time is required")
	}

	if !r.StartTime.IsZero() && !r.EndTime.IsZero() {
		if r.StartTime.Before(time.Now()) {
			errs.Set("start_time", "start time cannot be in the past")
		}

		if r.StartTime.After(r.EndTime) {
			errs.Set("start_time", "start time must be before end time")
		}

		if r.StartTime.Equal(r.EndTime) {
			errs.Set("start_time", "start time and end time cannot be the same")
		}

		// Check for reasonable duration (between 15 minutes and 12 hours)
		duration := r.EndTime.Sub(r.StartTime)
		if duration < 15*time.Minute {
			errs.Set("end_time", "availability duration must be at least 15 minutes")
		}
		if duration > 12*time.Hour {
			errs.Set("end_time", "availability duration cannot exceed 12 hours")
		}
	}

	// Validate capacity constraints
	if r.MaxCapacity <= 0 {
		errs.Set("max_capacity", "max capacity must be greater than 0")
	}

	if r.MaxCapacity > 50 {
		errs.Set("max_capacity", "max capacity cannot exceed 50")
	}

	// Validate recurrence pattern
	if r.Pattern.Type == "" {
		errs.Set("pattern.type", "recurrence type is required")
	} else {
		// Validate recurrence type
		validTypes := []string{"daily", "weekly", "monthly"}
		isValidType := false
		for _, validType := range validTypes {
			if r.Pattern.Type == validType {
				isValidType = true
				break
			}
		}
		if !isValidType {
			errs.Set("pattern.type", "recurrence type must be one of: daily, weekly, monthly")
		}
	}

	// Validate recurrence interval
	if r.Pattern.Interval <= 0 {
		errs.Set("pattern.interval", "recurrence interval must be greater than 0")
	}
	if r.Pattern.Interval > 52 { // Reasonable upper limit for weeks
		errs.Set("pattern.interval", "recurrence interval cannot exceed 52")
	}

	// Validate recurrence end date if provided
	if r.Pattern.Until != nil {
		if r.Pattern.Until.Before(time.Now()) {
			errs.Set("pattern.until", "recurrence end date cannot be in the past")
		}
		if !r.StartTime.IsZero() && r.Pattern.Until.Before(r.StartTime) {
			errs.Set("pattern.until", "recurrence end date must be after start time")
		}
	}

	// Validate days of week for weekly patterns
	if r.Pattern.Type == "weekly" && len(r.Pattern.DaysOfWeek) == 0 {
		errs.Set("pattern.days_of_week", "days of week are required for weekly recurrence")
	}

	if len(r.Pattern.DaysOfWeek) > 0 {
		for _, day := range r.Pattern.DaysOfWeek {
			if day < 0 || day > 6 {
				errs.Set("pattern.days_of_week", "days of week must be between 0 (Sunday) and 6 (Saturday)")
				break
			}
		}
	}

	// Validate optional fields
	if len(r.ServiceType) > 255 {
		errs.Set("service_type", "service type cannot exceed 255 characters")
	}

	if len(r.Notes) > 1000 {
		errs.Set("notes", "notes cannot exceed 1000 characters")
	}

	// Validate price if provided
	if r.PriceCents != nil {
		if *r.PriceCents < 0 {
			errs.Set("price_cents", "price cannot be negative")
		}
		if *r.PriceCents > 999999 { // $9,999.99 max
			errs.Set("price_cents", "price cannot exceed $9,999.99")
		}
	}

	return errs.AsError()
}

type UpdateAvailabilityRequest struct {
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
	ServiceType string    `json:"service_type,omitempty" validate:"max=255"`
	PriceCents  *int      `json:"price_cents,omitempty" validate:"omitempty,min=0"`
	Notes       string    `json:"notes,omitempty" validate:"max=1000"`
}

func (r *UpdateAvailabilityRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	return errs.AsError()
}
