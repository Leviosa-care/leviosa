package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type RoomResponse struct {
	ID          uuid.UUID `json:"id"`
	BuildingID  uuid.UUID `json:"building_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Capacity    int       `json:"capacity"`
	Equipment   []string  `json:"equipment,omitempty"`
	PriceCents  *int      `json:"price_cents,omitempty"`
	Currency    string    `json:"currency,omitempty"`
	IsActive    bool      `json:"is_active"`
}

// Room DTOs
type CreateRoomRequest struct {
	BuildingID  uuid.UUID `json:"building_id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	Description string    `json:"description,omitempty" validate:"max=1000"`
	RoomNumber  string    `json:"room_number,omitempty" encx:"encrypt"`
	Capacity    int       `json:"capacity" validate:"required,min=1,max=50"`
	Equipment   []string  `json:"equipment,omitempty"`
	PriceCents  *int      `json:"price_cents,omitempty" validate:"omitempty,min=0"`
	Currency    string    `json:"currency,omitempty" validate:"omitempty,len=3"`
	IsActive    bool      `json:"is_active"`
}

func (r *CreateRoomRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	return errs.AsError()
}

type UpdateRoomRequest struct {
	Name        string   `json:"name" validate:"required,min=1,max=255"`
	Description string   `json:"description,omitempty" validate:"max=1000"`
	Capacity    int      `json:"capacity" validate:"required,min=1,max=50"`
	Equipment   []string `json:"equipment,omitempty"`
	PriceCents  *int     `json:"price_cents,omitempty" validate:"omitempty,min=0"`
	Currency    string   `json:"currency,omitempty" validate:"omitempty,len=3"`
}

func (r *UpdateRoomRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	return errs.AsError()
}
