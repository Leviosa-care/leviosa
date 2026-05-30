package partner

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// VerifyPartner verifies a partner and updates their user role to "partner".
// This is an admin-only operation that:
// - Sets partner.IsVerified = true
// - Sets partner.VerifiedAt = time.Now()
// - Sets partner.VerifiedByUserID = verifiedByUserID
// - Updates user.Role = "partner"
// - Updates user.State = "active"
func (s *PartnerService) VerifyPartner(ctx context.Context, partnerID uuid.UUID, verifiedByUserID uuid.UUID) (*domain.PartnerResponse, error) {
	// Get partner to verify it exists
	partnerEncx, err := s.partnerRepo.GetPartnerByID(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("get partner by ID: %w", err)
	}

	// Decrypt partner to check if already verified
	partner, err := domain.DecryptPartnerEncx(ctx, s.crypto, partnerEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("partner during verification check", err)
	}

	// Check if partner is already verified
	if partner.StripeAccountStatus == domain.StripeAccountStatusActive && partner.StripeOnboardingComplete {
		return nil, errs.NewConflictErr(fmt.Errorf("partner is already verified"))
	}

	// Update partner verification status in repository (uses user_id, not partner_id)
	if err := s.partnerRepo.VerifyPartner(ctx, partner.UserID, verifiedByUserID); err != nil {
		return nil, fmt.Errorf("verify partner in repository: %w", err)
	}

	// Get user associated with the partner
	userEncx, err := s.userRepo.GetUserByID(ctx, partner.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	// Decrypt user for modification
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("user during partner verification", err)
	}

	// Update user role and state
	user.Role = identity.PartnerStr
	user.State = domain.Active

	// Re-encrypt user with modifications
	updatedUserEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("user during role update", err)
	}

	// Save updated user
	if err := s.userRepo.UpdateUser(ctx, updatedUserEncx); err != nil {
		return nil, fmt.Errorf("update user role and state: %w", err)
	}

	// Get updated partner with all fields
	updatedPartnerEncx, err := s.partnerRepo.GetPartnerByID(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("get updated partner: %w", err)
	}

	// Decrypt updated partner for response
	updatedPartner, err := domain.DecryptPartnerEncx(ctx, s.crypto, updatedPartnerEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("partner after verification", err)
	}

	// Return partner response
	return updatedPartner.ToResponse(), nil
}
