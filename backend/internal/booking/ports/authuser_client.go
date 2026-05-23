package ports

import (
	"context"

	"github.com/google/uuid"
)

// AuthUserClient defines the interface for communicating with the authuser service.
// This interface provides access to partner information needed by the booking module.
//
// Implementation Strategy:
// - Modular Monolith: In-process implementation delegates to authuser.PartnerService
// - Microservices: HTTP-based implementation uses ServiceClient for authenticated calls
type AuthUserClient interface {
	// GetPartnerVerificationStatus checks if a partner is verified by partner ID.
	// Returns false if partner doesn't exist (no error for not found).
	GetPartnerVerificationStatus(ctx context.Context, partnerID uuid.UUID) (bool, error)

	// GetPartnerByUserID retrieves a partner by their user ID.
	// This is needed because allocations are created using user IDs, not partner IDs.
	// Returns error if partner not found or if there's a system error.
	GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*PartnerInfo, error)

	// GetUserName returns the display name ("FirstName LastName") for a user.
	// Returns an empty string if the user is not found.
	GetUserName(ctx context.Context, userID uuid.UUID) (string, error)
}

// PartnerInfo represents basic partner information needed by booking module
type PartnerInfo struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	IsVerified bool      `json:"is_verified"`
}
