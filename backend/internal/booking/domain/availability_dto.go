package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type AvailabilityResponse struct {
	ID              uuid.UUID          `json:"id"`
	PartnerID       uuid.UUID          `json:"partner_id"`
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
	PartnerID   uuid.UUID `json:"partner_id" validate:"required"`
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
	return errs.AsError()
}

type CreateRecurringAvailabilityRequest struct {
	PartnerID   uuid.UUID         `json:"partner_id" validate:"required"`
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
