package roomHelpers

import (
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

func NewTestRoom(t *testing.T) *domain.Room {
	now := time.Now()
	return &domain.Room{
		ID:              uuid.New(),
		BuildingID:      uuid.New(),
		Name:            "Consultation Room A",
		Description:     "Spacious consultation room with examination equipment",
		RoomNumber:      "101",
		Capacity:        1,
		Equipment:       []string{"examination table", "chair", "sink", "cabinet"},
		HourlyRateCents: th.IntPtr(5000), // 50.00 EUR in cents
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// NewTestRoomWithParams creates a test room with custom parameters
func NewTestRoomWithParams(t *testing.T, buildingID uuid.UUID, name, roomNumber string, capacity int, hourlyRate int, isActive bool) *domain.Room {
	now := time.Now()
	return &domain.Room{
		ID:              uuid.New(),
		BuildingID:      buildingID,
		Name:            name,
		Description:     "Test room description",
		RoomNumber:      roomNumber,
		Capacity:        capacity,
		Equipment:       []string{"basic equipment"},
		HourlyRateCents: th.IntPtr(hourlyRate),
		IsActive:        isActive,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// NewTestRoomWithBuilding creates a test room associated with a specific building
func NewTestRoomWithBuilding(t *testing.T, buildingID uuid.UUID) *domain.Room {
	now := time.Now()
	return &domain.Room{
		ID:              uuid.New(),
		BuildingID:      buildingID,
		Name:            "Treatment Room B",
		Description:     "Standard treatment room",
		RoomNumber:      "201",
		Capacity:        2,
		Equipment:       []string{"treatment table", "lights", "storage"},
		HourlyRateCents: th.IntPtr(7500), // 75.00 EUR in cents
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func NewTestRoomEncx(t *testing.T) *domain.RoomEncx {
	now := time.Now()
	return &domain.RoomEncx{
		ID:                   uuid.New(),
		BuildingID:           uuid.New(),
		NameEncrypted:        []byte("encrypted_test_room_name"),
		NameHash:             "hashed_room_name",
		DescriptionEncrypted: []byte("encrypted_test_description"),
		RoomNumberEncrypted:  []byte("encrypted_101"),
		RoomNumberHash:       "hashed_101",
		Capacity:             1,
		EquipmentEncrypted:   []byte(`["equipment1", "equipment2"]`),
		HourlyRateCents:      &[]int{5000}[0],
		IsActive:             true,
		CreatedAt:            now,
		UpdatedAt:            now,
		DEKEncrypted:         []byte("mock_dek_data"),
		KeyVersion:           1,
		Metadata:             encx.EncryptionMetadata{},
	}
}

// NewTestRoomEncxWithBuilding creates an encrypted test room with specific building ID
func NewTestRoomEncxWithBuilding(t *testing.T, buildingID uuid.UUID) *domain.RoomEncx {
	now := time.Now()
	rate := 7500
	return &domain.RoomEncx{
		ID:                   uuid.New(),
		BuildingID:           buildingID,
		NameEncrypted:        []byte("encrypted_treatment_room"),
		NameHash:             "hashed_treatment_room",
		DescriptionEncrypted: []byte("encrypted_treatment_desc"),
		RoomNumberEncrypted:  []byte("encrypted_201"),
		RoomNumberHash:       "hashed_201",
		Capacity:             2,
		EquipmentEncrypted:   []byte(`["treatment table", "lights"]`),
		HourlyRateCents:      &rate,
		IsActive:             true,
		CreatedAt:            now,
		UpdatedAt:            now,
		DEKEncrypted:         []byte("mock_dek_data"),
		KeyVersion:           1,
		Metadata:             encx.EncryptionMetadata{},
	}
}

// NewInactiveTestRoomEncx creates an inactive encrypted test room
func NewInactiveTestRoomEncx(t *testing.T, buildingID uuid.UUID) *domain.RoomEncx {
	now := time.Now()
	return &domain.RoomEncx{
		ID:                   uuid.New(),
		BuildingID:           buildingID,
		NameEncrypted:        []byte("encrypted_inactive_room"),
		NameHash:             "hashed_inactive_room",
		DescriptionEncrypted: []byte("encrypted_inactive_desc"),
		RoomNumberEncrypted:  []byte("encrypted_999"),
		RoomNumberHash:       "hashed_999",
		Capacity:             1,
		EquipmentEncrypted:   []byte(`[]`),
		HourlyRateCents:      nil, // No rate set
		IsActive:             false,
		CreatedAt:            now,
		UpdatedAt:            now,
		DEKEncrypted:         []byte("mock_dek_data"),
		KeyVersion:           1,
		Metadata:             encx.EncryptionMetadata{},
	}
}
