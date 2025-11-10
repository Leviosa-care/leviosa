package domain

import (
	"time"

	"github.com/google/uuid"
)

// Building represents a physical location containing treatment rooms
type Building struct {
	ID uuid.UUID `json:"id"`

	// Name and address (encrypted for GDPR compliance)
	Name       string `json:"name" encx:"encrypt"`
	Address    string `json:"address" encx:"encrypt"`
	City       string `json:"city" encx:"encrypt"`
	PostalCode string `json:"postal_code" encx:"encrypt"`
	Country    string `json:"country" encx:"encrypt"`

	// Business information
	Description string `json:"description,omitempty" encx:"encrypt"`
	Phone       string `json:"phone,omitempty" encx:"encrypt"`
	Email       string `json:"email,omitempty" encx:"encrypt"`

	// Administrative fields
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *Building) ToResponse() *BuildingResponse {
	return &BuildingResponse{
		ID:          b.ID,
		Name:        b.Name,
		Address:     b.Address,
		City:        b.City,
		PostalCode:  b.PostalCode,
		Country:     b.Country,
		Description: b.Description,
		Phone:       b.Phone,
		Email:       b.Email,
	}

}

// UpdateDetails updates the building's details
func (b *Building) UpdateDetails(name, address, city, postalCode, country string) error {
	if name == "" {
		return ErrInvalidBuildingName
	}
	if address == "" {
		return ErrInvalidBuildingAddress
	}
	if city == "" {
		return ErrInvalidBuildingCity
	}
	if country == "" {
		return ErrInvalidBuildingCountry
	}

	b.Name = name
	b.Address = address
	b.City = city
	b.PostalCode = postalCode
	b.Country = country
	b.UpdatedAt = time.Now()
	return nil
}

// SetContactInfo sets optional contact information
func (b *Building) SetContactInfo(description, phone, email string) {
	b.Description = description
	b.Phone = phone
	b.Email = email
	b.UpdatedAt = time.Now()
}

// Deactivate marks the building as inactive
func (b *Building) Deactivate() {
	b.IsActive = false
	b.UpdatedAt = time.Now()
}

// Activate marks the building as active
func (b *Building) Activate() {
	b.IsActive = true
	b.UpdatedAt = time.Now()
}
