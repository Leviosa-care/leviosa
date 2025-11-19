package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type RoomAllocationResponse struct {
	ID             uuid.UUID      `json:"id"`
	RoomID         uuid.UUID      `json:"room_id"`
	UserID         uuid.UUID      `json:"user_id"`
	AllocationType AllocationType `json:"allocation_type"`
	StartDate      *time.Time     `json:"start_date,omitempty"`
	EndDate        *time.Time     `json:"end_date,omitempty"`
	IsActive       bool           `json:"is_active"`
}

// Room Allocation DTOs
type CreateSharedAllocationRequest struct {
	RoomID uuid.UUID `json:"room_id" validate:"required"`
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func (r *CreateSharedAllocationRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate RoomID
	if r.RoomID == uuid.Nil {
		errs.Set("room_id", "room ID is required")
	}

	// Validate UserID
	if r.UserID == uuid.Nil {
		errs.Set("user_id", "user ID is required")
	}

	return errs.AsError()
}

type CreateDedicatedAllocationRequest struct {
	RoomID    uuid.UUID  `json:"room_id" validate:"required"`
	UserID    uuid.UUID  `json:"user_id" validate:"required"`
	StartDate *time.Time `json:"start_date" validate:"required"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

func (r *CreateDedicatedAllocationRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate RoomID
	if r.RoomID == uuid.Nil {
		errs.Set("room_id", "room ID is required")
	}

	// Validate UserID
	if r.UserID == uuid.Nil {
		errs.Set("user_id", "user ID is required")
	}

	// Validate StartDate
	if r.StartDate == nil {
		errs.Set("start_date", "start date is required")
	}

	// Validate EndDate is after StartDate (if both are provided)
	if r.StartDate != nil && r.EndDate != nil {
		if r.EndDate.Before(*r.StartDate) || r.EndDate.Equal(*r.StartDate) {
			errs.Set("end_date", "end date must be after start date")
		}
	}

	return errs.AsError()
}

type UpdateDedicatedAllocationRequest struct {
	StartDate *time.Time `json:"start_date" validate:"required"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

func (r *UpdateDedicatedAllocationRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	return errs.AsError()
}
