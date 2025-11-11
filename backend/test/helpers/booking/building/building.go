package buildingHelpers

import (
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

func NewTestBuilding(t *testing.T) *domain.Building {
	now := time.Now()
	return &domain.Building{
		ID:          uuid.New(),
		Name:        "Name",
		Address:     "123 Rue de Rivoli",
		City:        "Paris",
		PostalCode:  "75001",
		Country:     "France",
		Description: "A description of the building",
		Phone:       "0612345678",
		Email:       "building@example.fr",
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewTestBuildingWithParams creates a test building with custom parameters
func NewTestBuildingWithParams(t *testing.T, name, city, country string, isActive bool) *domain.Building {
	now := time.Now()
	return &domain.Building{
		ID:          uuid.New(),
		Name:        name,
		Address:     "123 Test Street",
		City:        city,
		PostalCode:  "12345",
		Country:     country,
		Description: "Test building description",
		Phone:       "+1234567890",
		Email:       "test@example.com",
		IsActive:    isActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func NewTestBuildingEncx(t *testing.T) *domain.BuildingEncx {
	now := time.Now()
	return &domain.BuildingEncx{
		ID:                   uuid.New(),
		NameEncrypted:        []byte("encrypted_test_building_name"),
		AddressEncrypted:     []byte("encrypted_123_test_street"),
		CityEncrypted:        []byte("encrypted_test_city"),
		PostalCodeEncrypted:  []byte("encrypted_12345"),
		CountryEncrypted:     []byte("encrypted_test_country"),
		DescriptionEncrypted: []byte("encrypted_test_description"),
		PhoneEncrypted:       []byte("encrypted_+1234567890"),
		EmailEncrypted:       []byte("encrypted_test@example.com"),
		IsActive:             true,
		CreatedAt:            now,
		UpdatedAt:            now,
		DEKEncrypted:         []byte("mock_dek_data"),
		KeyVersion:           1,
		Metadata:             encx.EncryptionMetadata{},
	}
}
