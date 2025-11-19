package partner

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// GetPartnerVerificationStatus checks if a partner is verified.
// A partner is considered verified when:
// - stripe_account_status = 'active'
// - stripe_onboarding_complete = true
//
// This method is part of PublicPartnerService and is safe for inter-service calls
// as it only reads data and doesn't require user session context.
func (s *PartnerService) GetPartnerVerificationStatus(ctx context.Context, partnerID uuid.UUID) (bool, error) {
	// Get partner from repository
	partnerEncx, err := s.partnerRepo.GetPartnerByID(ctx, partnerID)
	if err != nil {
		// If partner not found, return false (not verified) instead of error
		// This allows callers to use this method for existence checks
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("get partner by ID: %w", err)
	}

	// Decrypt partner to check verification status
	partner, err := domain.DecryptPartnerEncx(ctx, s.crypto, partnerEncx)
	if err != nil {
		return false, errs.NewNotDecryptedErr("partner during verification status check", err)
	}

	// Partner is verified when both conditions are met
	isVerified := partner.StripeAccountStatus == domain.StripeAccountStatusActive &&
		partner.StripeOnboardingComplete

	return isVerified, nil
}
