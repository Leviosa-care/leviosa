package domain

import (
	"context"
	"fmt"
	"strings"

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
	IsActive    bool      `json:"is_active"`
}

func (r *CreateRoomRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// BuildingID validation
	if r.BuildingID == uuid.Nil {
		errs.Set("building_id", fmt.Errorf("building ID is required"))
	}

	// Name validation
	if strings.TrimSpace(r.Name) == "" {
		errs.Set("name", fmt.Errorf("name is required"))
	} else if len([]rune(strings.TrimSpace(r.Name))) > 255 {
		errs.Set("name", fmt.Errorf("name cannot exceed 255 characters"))
	}

	// Description validation
	if r.Description != "" {
		if len([]rune(strings.TrimSpace(r.Description))) > 1000 {
			errs.Set("description", fmt.Errorf("description cannot exceed 1000 characters"))
		}
	}

	// RoomNumber validation (optional)
	if r.RoomNumber != "" {
		if len([]rune(strings.TrimSpace(r.RoomNumber))) > 50 {
			errs.Set("room_number", fmt.Errorf("room number cannot exceed 50 characters"))
		}
	}

	// Capacity validation
	if r.Capacity <= 0 {
		errs.Set("capacity", fmt.Errorf("capacity must be at least 1"))
	} else if r.Capacity > 50 {
		errs.Set("capacity", fmt.Errorf("capacity cannot exceed 50"))
	}

	// Equipment validation (optional)
	if r.Equipment != nil {
		if len(r.Equipment) > 20 {
			errs.Set("equipment", fmt.Errorf("cannot have more than 20 equipment items"))
		}
		for i, item := range r.Equipment {
			if strings.TrimSpace(item) == "" {
				errs.Set("equipment", fmt.Errorf("equipment item %d cannot be empty", i+1))
				break
			}
			if len([]rune(strings.TrimSpace(item))) > 100 {
				errs.Set("equipment", fmt.Errorf("equipment item %d cannot exceed 100 characters", i+1))
				break
			}
		}
	}

	return errs.AsError()
}

type UpdateRoomRequest struct {
	ID          uuid.UUID `json:"id"`
	BuildingID  uuid.UUID `json:"building_id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description,omitempty"`
	RoomNumber  *string   `json:"room_number,omitempty"`
	Capacity    *int      `json:"capacity"`
	Equipment   *[]string `json:"equipment,omitempty"`
	IsActive    *bool     `json:"is_active,omitempty"`
}

func (r *UpdateRoomRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// ID validation (required)
	if r.ID == uuid.Nil {
		errs.Set("id", fmt.Errorf("room ID is required"))
	}

	// BuildingID validation (if not nil UUID, validate it)
	if r.BuildingID != uuid.Nil {
		// BuildingID is being updated, no additional validation needed for UUID
	}

	// Name validation (only if provided)
	if r.Name != nil {
		if strings.TrimSpace(*r.Name) == "" {
			errs.Set("name", fmt.Errorf("name cannot be empty"))
		} else if len([]rune(strings.TrimSpace(*r.Name))) > 255 {
			errs.Set("name", fmt.Errorf("name cannot exceed 255 characters"))
		}
	}

	// Description validation (only if provided)
	if r.Description != nil {
		if len([]rune(strings.TrimSpace(*r.Description))) > 1000 {
			errs.Set("description", fmt.Errorf("description cannot exceed 1000 characters"))
		}
	}

	// RoomNumber validation (only if provided)
	if r.RoomNumber != nil {
		if len([]rune(strings.TrimSpace(*r.RoomNumber))) > 50 {
			errs.Set("room_number", fmt.Errorf("room number cannot exceed 50 characters"))
		}
	}

	// Capacity validation (only if provided)
	if r.Capacity != nil {
		if *r.Capacity <= 0 {
			errs.Set("capacity", fmt.Errorf("capacity must be at least 1"))
		} else if *r.Capacity > 50 {
			errs.Set("capacity", fmt.Errorf("capacity cannot exceed 50"))
		}
	}

	// Equipment validation (only if provided)
	if r.Equipment != nil {
		equipmentList := *r.Equipment
		if len(equipmentList) > 20 {
			errs.Set("equipment", fmt.Errorf("cannot have more than 20 equipment items"))
		}
		for i, item := range equipmentList {
			if strings.TrimSpace(item) == "" {
				errs.Set("equipment", fmt.Errorf("equipment item %d cannot be empty", i+1))
				break
			}
			if len([]rune(strings.TrimSpace(item))) > 100 {
				errs.Set("equipment", fmt.Errorf("equipment item %d cannot exceed 100 characters", i+1))
				break
			}
		}
	}

	return errs.AsError()
}
