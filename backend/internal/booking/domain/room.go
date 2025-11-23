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
	Name        string `json:"name" encx:"encrypt,hash_basic"`
	Description string `json:"description,omitempty" encx:"encrypt"`
	RoomNumber  string `json:"room_number,omitempty" encx:"encrypt,hash_basic"`

	// Room specifications
	Capacity  int      `json:"capacity"`
	Equipment []string `json:"equipment,omitempty" encx:"encrypt"`

	// Operating hours (time of day in HH:MM format)
	OperatingStartTime time.Time `json:"operating_start_time"`
	OperatingEndTime   time.Time `json:"operating_end_time"`

	// Administrative fields
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *Room) ToResponse() *RoomResponse {
	return &RoomResponse{
		ID:                 r.ID,
		BuildingID:         r.BuildingID,
		Name:               r.Name,
		Description:        r.Description,
		RoomNumber:         r.RoomNumber,
		Capacity:           r.Capacity,
		Equipment:          r.Equipment,
		OperatingStartTime: r.OperatingStartTime,
		OperatingEndTime:   r.OperatingEndTime,
		IsActive:           r.IsActive,
	}
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
