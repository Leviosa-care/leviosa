package domain

import (
	"time"

	"github.com/google/uuid"
)

// Room represents an individual treatment space within a building
type Room struct {
	ID         uuid.UUID `json:"id"`
	BuildingID uuid.UUID `json:"building_id"`

	// Room identification (encrypted)
	Name                   string `json:"name" encx:"encrypt"`
	NameEncrypted          []byte `json:"-"`
	Description            string `json:"description,omitempty" encx:"encrypt"`
	DescriptionEncrypted   []byte `json:"-"`
	RoomNumber             string `json:"room_number,omitempty" encx:"encrypt"`
	RoomNumberEncrypted    []byte `json:"-"`

	// Room specifications
	Capacity              int      `json:"capacity"`
	Equipment             []string `json:"equipment,omitempty" encx:"encrypt"`
	EquipmentEncrypted    []byte   `json:"-"`

	// Pricing
	HourlyRateCents *int `json:"hourly_rate_cents,omitempty"`

	// Administrative fields
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewRoom creates a new Room with validated data
func NewRoom(buildingID uuid.UUID, name string, capacity int) (*Room, error) {
	if buildingID == uuid.Nil {
		return nil, ErrInvalidBuildingID
	}
	if name == "" {
		return nil, ErrInvalidRoomName
	}
	if capacity <= 0 {
		return nil, ErrInvalidRoomCapacity
	}

	return &Room{
		ID:         uuid.New(),
		BuildingID: buildingID,
		Name:       name,
		Capacity:   capacity,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

// UpdateDetails updates the room's basic information
func (r *Room) UpdateDetails(name, description, roomNumber string, capacity int) error {
	if name == "" {
		return ErrInvalidRoomName
	}
	if capacity <= 0 {
		return ErrInvalidRoomCapacity
	}

	r.Name = name
	r.Description = description
	r.RoomNumber = roomNumber
	r.Capacity = capacity
	r.UpdatedAt = time.Now()
	return nil
}

// SetEquipment updates the room's equipment list
func (r *Room) SetEquipment(equipment []string) {
	r.Equipment = make([]string, len(equipment))
	copy(r.Equipment, equipment)
	r.UpdatedAt = time.Now()
}

// SetHourlyRate sets the base hourly rate for the room
func (r *Room) SetHourlyRate(rateCents int) error {
	if rateCents < 0 {
		return ErrInvalidRoomRate
	}
	r.HourlyRateCents = &rateCents
	r.UpdatedAt = time.Now()
	return nil
}

// ClearHourlyRate removes the base hourly rate
func (r *Room) ClearHourlyRate() {
	r.HourlyRateCents = nil
	r.UpdatedAt = time.Now()
}

// Deactivate marks the room as inactive
func (r *Room) Deactivate() {
	r.IsActive = false
	r.UpdatedAt = time.Now()
}

// Activate marks the room as active
func (r *Room) Activate() {
	r.IsActive = true
	r.UpdatedAt = time.Now()
}