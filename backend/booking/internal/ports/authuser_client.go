package ports

import (
	"context"

	"github.com/google/uuid"
)

// AuthUserClient defines the interface for communicating with the authuser service
type AuthUserClient interface {
	// GetPartnerByID retrieves a partner by their ID from the authuser service
	GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*PartnerInfo, error)

	// GetPartnerByUserID retrieves a partner by their user ID from the authuser service
	GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*PartnerInfo, error)

	// ValidatePartnerExists checks if a partner exists and is verified
	ValidatePartnerExists(ctx context.Context, partnerID uuid.UUID) (bool, error)

	// GetPartnerSpecializations retrieves a partner's specializations
	GetPartnerSpecializations(ctx context.Context, partnerID uuid.UUID) ([]SpecializationInfo, error)
}

// PartnerInfo represents partner information from the authuser service
type PartnerInfo struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	Bio              string    `json:"bio"`
	Experience       string    `json:"experience"`
	Certifications   []string  `json:"certifications"`
	IsVerified       bool      `json:"is_verified"`
	Email            string    `json:"email"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Telephone        string    `json:"telephone"`
}

// SpecializationInfo represents specialization information from the authuser service
type SpecializationInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
}