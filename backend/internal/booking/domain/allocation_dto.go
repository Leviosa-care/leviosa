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
