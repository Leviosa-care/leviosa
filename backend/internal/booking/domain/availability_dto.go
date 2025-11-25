package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type AvailabilityResponse struct {
	ID                uuid.UUID          `json:"id"`
	UserID            uuid.UUID          `json:"partner_id"`
	RoomID            uuid.UUID          `json:"room_id"`
	StartTime         time.Time          `json:"start_time"`
	EndTime           time.Time          `json:"end_time"`
	MaxCapacity       int                `json:"max_capacity"`
	CurrentBookings   int                `json:"current_bookings"`
	Status            AvailabilityStatus `json:"status"`
	ServiceType       string             `json:"service_type,omitempty"`
	PriceCents        *int               `json:"price_cents,omitempty"`
	Notes             string             `json:"notes,omitempty"`
	RecurrencePattern *RecurrencePattern `json:"recurrence_pattern,omitempty"`
	ParentID          *uuid.UUID         `json:"parent_id,omitempty"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
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
	ID          uuid.UUID  `json:"id"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	ServiceType *string    `json:"service_type,omitempty"`
	PriceCents  *int       `json:"price_cents,omitempty"`
	Notes       *string    `json:"notes,omitempty"`
}

func (r *UpdateAvailabilityRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate ID is required
	if r.ID == uuid.Nil {
		errs.Set("id", "availability ID is required")
	}

	// Validate time constraints (now optional due to pointers)
	if r.StartTime != nil && r.EndTime != nil {
		if r.StartTime.Before(time.Now()) {
			errs.Set("start_time", "start time cannot be in the past")
		}

		if r.StartTime.After(*r.EndTime) {
			errs.Set("start_time", "start time must be before end time")
		}

		if r.StartTime.Equal(*r.EndTime) {
			errs.Set("start_time", "start time and end time cannot be the same")
		}

		// Check for reasonable duration (between 15 minutes and 12 hours)
		duration := r.EndTime.Sub(*r.StartTime)
		if duration < 15*time.Minute {
			errs.Set("end_time", "availability duration must be at least 15 minutes")
		}
		if duration > 12*time.Hour {
			errs.Set("end_time", "availability duration cannot exceed 12 hours")
		}
	} else if r.StartTime != nil || r.EndTime != nil {
		// If one time is provided, both must be provided
		errs.Set("time", "both start_time and end_time must be provided together")
	}

	// Validate optional string fields (now pointers)
	if r.ServiceType != nil && len(*r.ServiceType) > 255 {
		errs.Set("service_type", "service type cannot exceed 255 characters")
	}

	if r.Notes != nil && len(*r.Notes) > 1000 {
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

// Gap Detection DTOs

// GetRoomGapsRequest represents a request to find time gaps in a room's schedule
type GetRoomGapsRequest struct {
	RoomID uuid.UUID `json:"room_id" validate:"required"`
	Date   time.Time `json:"date" validate:"required"`
}

// GetRoomGapsResponse represents the time gaps found in a room's schedule
type GetRoomGapsResponse struct {
	RoomID          uuid.UUID              `json:"room_id"`
	Date            time.Time              `json:"date"`
	OperatingHours  OperatingHoursResponse `json:"operating_hours"`
	Gaps            []TimeGapResponse      `json:"gaps"`
	TotalGapMinutes int                    `json:"total_gap_minutes"`
}

// OperatingHoursResponse represents a room's operating hours for a specific date
type OperatingHoursResponse struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// TimeGapResponse represents a time gap in the API response
type TimeGapResponse struct {
	StartTime         time.Time           `json:"start_time"`
	EndTime           time.Time           `json:"end_time"`
	DurationMinutes   int                 `json:"duration_minutes"`
	IsBookable        bool                `json:"is_bookable"`
	SuggestedProducts []ProductSuggestion `json:"suggested_products"`
}

// ProductSuggestion represents a product that fits in a time gap
type ProductSuggestion struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Duration    int       `json:"duration"`
	BufferTime  int       `json:"buffer_time"`
	TotalTime   int       `json:"total_time"` // Duration + BufferTime
}

// Availability Suggestion DTOs

// GetAvailabilitySuggestionsRequest represents a request for availability block suggestions
type GetAvailabilitySuggestionsRequest struct {
	PartnerID uuid.UUID `json:"partner_id" validate:"required"`
	RoomID    uuid.UUID `json:"room_id" validate:"required"`
}

// GetAvailabilitySuggestionsResponse contains recommended availability block durations
type GetAvailabilitySuggestionsResponse struct {
	PartnerID         uuid.UUID                 `json:"partner_id"`
	RoomID            uuid.UUID                 `json:"room_id"`
	AllocationType    string                    `json:"allocation_type"`
	RecommendedBlocks []BlockSuggestionResponse `json:"recommended_blocks"`
}

// BlockSuggestionResponse represents a recommended block in the API response
type BlockSuggestionResponse struct {
	DurationMinutes     int                    `json:"duration_minutes"`
	Rationale           string                 `json:"rationale"`
	ProductCombinations []ProductComboResponse `json:"product_combinations"`
	Priority            int                    `json:"priority"`
}

// ProductComboResponse represents a product combination in the API response
type ProductComboResponse struct {
	Products      []ProductInfoResponse `json:"products"`
	TotalDuration int                   `json:"total_duration"`
	SessionCount  int                   `json:"session_count"`
}

// ProductInfoResponse contains product information in the API response
type ProductInfoResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Duration   int       `json:"duration"`
	BufferTime int       `json:"buffer_time"`
}
