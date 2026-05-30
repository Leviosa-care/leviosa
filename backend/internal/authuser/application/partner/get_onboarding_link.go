package partner

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetOnboardingLink generates a Stripe Account Link URL for the partner to complete onboarding.
// If the partner's StripeConnectedAccountID is empty (e.g. Stripe account creation failed at
// registration), a new Stripe Connect Express account is created first, then persisted back.
func (s *PartnerService) GetOnboardingLink(ctx context.Context, userID uuid.UUID, returnURL, refreshURL string) (string, error) {
	// Get encrypted partner from repository
	partnerEncx, err := s.partnerRepo.GetPartnerByUserID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("get partner by user ID: %w", err)
	}

	// Decrypt partner
	partner, err := domain.DecryptPartnerEncx(ctx, s.crypto, partnerEncx)
	if err != nil {
		return "", errs.NewNotDecryptedErr("partner", err)
	}

	// If the partner has no Stripe account ID, create one first.
	if partner.StripeConnectedAccountID == "" {
		accountID, err := s.stripe.CreateConnectedAccount(ctx, userID)
		if err != nil {
			return "", fmt.Errorf("create connected account: %w", err)
		}
		partner.StripeConnectedAccountID = accountID
		partner.StripeAccountStatus = domain.StripeAccountStatusPending

		// Re-encrypt and persist the updated partner
		updatedEncx, err := domain.ProcessPartnerEncx(ctx, s.crypto, partner)
		if err != nil {
			return "", errs.NewNotEncryptedErr("partner during onboarding link", err)
		}
		if err := s.partnerRepo.UpdatePartner(ctx, updatedEncx); err != nil {
			return "", fmt.Errorf("update partner with stripe account ID: %w", err)
		}
	}

	// Create the Stripe Account Link
	url, err := s.stripe.CreateAccountLink(
		ctx,
		partner.StripeConnectedAccountID,
		"account_onboarding",
		returnURL,
		refreshURL,
	)
	if err != nil {
		return "", fmt.Errorf("create account link: %w", err)
	}

	return url, nil
}
