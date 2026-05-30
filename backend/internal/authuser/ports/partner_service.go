package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
)

// PublicPartnerService defines read-only operations safe for inter-service communication.
// This interface is exposed to other services (e.g., booking, notification) for partner
// information retrieval without requiring user authentication context.
type PublicPartnerService interface {
	// GetPartnerVerificationStatus checks if a partner is verified.
	// A partner is considered verified when:
	// - stripe_account_status = 'active'
	// - stripe_onboarding_complete = true
	// - user.role = 'partner'
	GetPartnerVerificationStatus(ctx context.Context, partnerID uuid.UUID) (bool, error)

	// GetPartnerByUserID retrieves basic partner information by user ID.
	// Returns error if partner not found or system error occurs.
	GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.PartnerResponse, error)
}

// PartnerService defines the full partner service interface including both
// public operations (read-only, safe for inter-service calls) and private
// operations (mutations, require user session context).
type PartnerService interface {
	PublicPartnerService
	GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.PartnerResponse, error)
	GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.PartnerResponse, error)
	GetAllPartners(ctx context.Context) ([]*domain.PartnerResponse, error)
	GetPublicPartners(ctx context.Context) ([]*domain.PublicPartnerResponse, error)
	GetAllPartnersByCategory(ctx context.Context, categoryID string) ([]*domain.PartnerResponse, error) // fetch the cache to get the cate
	GetAllPartnersByCategories(ctx context.Context, categoryIDs []string) ([]*domain.PartnerResponse, error)
	GetAllPartnersByProduct(ctx context.Context, productID string) ([]*domain.PartnerResponse, error) // fetch the cache to get the cate
	GetAllPartnersByProducts(ctx context.Context, productIDs []string) ([]*domain.PartnerResponse, error)
	CreatePartner(ctx context.Context, userID uuid.UUID, bio, experience string, categoryIDs, productIDs []uuid.UUID) (*domain.Partner, error)
	UpdatePartner(ctx context.Context, partnerID uuid.UUID, request *domain.UpdatePartnerRequest) (*domain.PartnerResponse, error)
	DeletePartner(ctx context.Context, partnerID uuid.UUID) error
	VerifyPartner(ctx context.Context, partnerID uuid.UUID, verifiedByUserID uuid.UUID) (*domain.PartnerResponse, error)

	// GetOnboardingLink generates a Stripe Account Link URL for the partner to complete onboarding.
	// If the partner has no StripeConnectedAccountID, a new Stripe account is created first.
	// Returns the URL to redirect the partner to.
	GetOnboardingLink(ctx context.Context, userID uuid.UUID, returnURL, refreshURL string) (string, error)

	// UpdateStripeAccountStatus updates a partner's stripe_account_status by looking up
	// the partner whose StripeConnectedAccountID matches the given accountID.
	// Returns the partner ID of the updated partner, or an error if not found.
	UpdateStripeAccountStatus(ctx context.Context, stripeAccountID string, status domain.StripeAccountStatus) (uuid.UUID, error)
}
