package partner

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// NOTE: The methods below are stubs for future implementation.
// Partner registration is now handled via /auth/complete/partner endpoint.

// GetPartnerByID retrieves a partner by ID with their associated user information.
func (s *PartnerService) GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.PartnerResponse, error) {
	// Get encrypted partner from repository
	partnerEncx, err := s.partnerRepo.GetPartnerByUserID(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("get partner by ID from repository: %w", err)
	}

	// Decrypt partner
	partner, err := domain.DecryptPartnerEncx(ctx, s.crypto, partnerEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("partner", err)
	}

	// Build complete response with user info
	// return partner, nil
	return &domain.PartnerResponse{
		UserID:     partner.UserID,
		Bio:        partner.Bio,
		Experience: partner.Experience,
		// Certifications: partner.Certifications,
		CategoryIDs: partner.CategoryIDs,
		ProductIDs:  partner.ProductIDs,
		CreatedAt:   partner.CreatedAt,
		UpdatedAt:   partner.UpdatedAt,
	}, nil
}

// GetPartnerByUserID retrieves a partner by user ID with their associated user information.
func (s *PartnerService) GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.PartnerResponse, error) {
	// Get encrypted partner from repository
	partnerEncx, err := s.partnerRepo.GetPartnerByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get partner by user ID from repository: %w", err)
	}

	// Decrypt partner
	partner, err := domain.DecryptPartnerEncx(ctx, s.crypto, partnerEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("partner", err)
	}

	// Build complete response with user info
	return &domain.PartnerResponse{
		UserID:     partner.UserID,
		Bio:        partner.Bio,
		Experience: partner.Experience,
		// Certifications: partner.Certifications,
		CategoryIDs: partner.CategoryIDs,
		ProductIDs:  partner.ProductIDs,
		CreatedAt:   partner.CreatedAt,
		UpdatedAt:   partner.UpdatedAt,
	}, nil
}

// GetAllPartners retrieves all partners with their associated user information.
func (s *PartnerService) GetAllPartners(ctx context.Context) (*domain.GetPartnersResponse, error) {
	// Get all partners from repository
	partnersEncx, err := s.partnerRepo.GetAllPartners(ctx)
	if err != nil {
		return nil, fmt.Errorf("get partners with users from repository: %w", err)
	}

	// Decrypt partners and build response
	partners := make([]domain.PartnerResponse, 0, len(partnersEncx))
	for _, partnerEncx := range partnersEncx {
		// Decrypt partner
		partner, err := domain.DecryptPartnerEncx(ctx, s.crypto, partnerEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("partner", err)
		}

		// Build complete partner response
		partners = append(partners, domain.PartnerResponse{
			UserID:     partner.UserID,
			Bio:        partner.Bio,
			Experience: partner.Experience,
			// Certifications: partner.Certifications,
			CategoryIDs: partner.CategoryIDs,
			ProductIDs:  partner.ProductIDs,
			CreatedAt:   partner.CreatedAt,
			UpdatedAt:   partner.UpdatedAt,
		})
	}

	return &domain.GetPartnersResponse{
		Partners: partners,
		Total:    len(partners),
	}, nil
}

// UpdatePartner updates an existing partner's profile fields.
// Only updates fields that are provided (non-nil) in the request.
func (s *PartnerService) UpdatePartner(ctx context.Context, partnerID uuid.UUID, request *domain.UpdatePartnerRequest) (*domain.PartnerResponse, error) {
	// Validate request
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get existing partner
	partnerEncx, err := s.partnerRepo.GetPartnerByUserID(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("get partner by ID: %w", err)
	}

	// Decrypt partner
	partner, err := domain.DecryptPartnerEncx(ctx, s.crypto, partnerEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("partner", err)
	}

	// Update only provided fields
	if request.Bio != nil {
		partner.Bio = *request.Bio
	}
	if request.Experience != nil {
		partner.Experience = *request.Experience
	}
	// if request.Certifications != nil {
	//	partner.Certifications = *request.Certifications
	// }

	// Re-encrypt partner with updated fields
	updatedPartnerEncx, err := domain.ProcessPartnerEncx(ctx, s.crypto, partner)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("partner during update", err)
	}

	// Save updated partner
	if err := s.partnerRepo.UpdatePartner(ctx, updatedPartnerEncx); err != nil {
		return nil, fmt.Errorf("update partner in repository: %w", err)
	}

	// Return updated partner response
	return &domain.PartnerResponse{
		UserID:     partner.UserID,
		Bio:        partner.Bio,
		Experience: partner.Experience,
		// Certifications: partner.Certifications,
		CategoryIDs: partner.CategoryIDs,
		ProductIDs:  partner.ProductIDs,
		CreatedAt:   partner.CreatedAt,
		UpdatedAt:   partner.UpdatedAt,
	}, nil
}

// DeletePartner deletes a partner by ID.
// This is an admin-only operation that removes the partner profile but does NOT delete the user account.
func (s *PartnerService) DeletePartner(ctx context.Context, partnerID uuid.UUID) error {
	// Verify partner exists
	_, err := s.partnerRepo.GetPartnerByUserID(ctx, partnerID)
	if err != nil {
		return fmt.Errorf("get partner by ID: %w", err)
	}

	// Delete partner
	if err := s.partnerRepo.DeletePartner(ctx, partnerID); err != nil {
		return fmt.Errorf("delete partner from repository: %w", err)
	}

	return nil
}

// VerifyPartner verifies a partner and updates their user role to "partner".
// This is an admin-only operation that:
// - Sets partner.IsVerified = true
// - Sets partner.VerifiedAt = time.Now()
// - Sets partner.VerifiedByUserID = verifiedByUserID
// - Updates user.Role = "partner"
// - Updates user.State = "active"
func (s *PartnerService) VerifyPartner(ctx context.Context, partnerID uuid.UUID, verifiedByUserID uuid.UUID) (*domain.PartnerResponse, error) {
	// Get partner to verify it exists
	partnerEncx, err := s.partnerRepo.GetPartnerByUserID(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("get partner by ID: %w", err)
	}

	// Decrypt partner to check if already verified
	partner, err := domain.DecryptPartnerEncx(ctx, s.crypto, partnerEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("partner during verification check", err)
	}

	// Update partner verification status in repository
	if err := s.partnerRepo.VerifyPartner(ctx, partnerID, verifiedByUserID); err != nil {
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
	updatedPartnerEncx, err := s.partnerRepo.GetPartnerByUserID(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("get updated partner: %w", err)
	}

	// Decrypt updated partner for response
	updatedPartner, err := domain.DecryptPartnerEncx(ctx, s.crypto, updatedPartnerEncx)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("partner after verification", err)
	}

	// 7. Return partner response
	return &domain.PartnerResponse{
		UserID:      updatedPartner.UserID,
		Bio:         updatedPartner.Bio,
		Experience:  updatedPartner.Experience,
		CategoryIDs: updatedPartner.CategoryIDs,
		ProductIDs:  updatedPartner.ProductIDs,
		CreatedAt:   updatedPartner.CreatedAt,
		UpdatedAt:   updatedPartner.UpdatedAt,
	}, nil
}

func (s *PartnerService) GetAllPartnersByCategory(ctx context.Context, category string) (*domain.GetPartnersResponse, error) {
	return nil, nil
}
func (s *PartnerService) GetAllPartnersByCategories(ctx context.Context, category []string) (*domain.GetPartnersResponse, error) {
	return nil, nil
}

func (s *PartnerService) UpdateCategories(ctx context.Context, categories []string) error {
	return nil
}
func (s *PartnerService) UpdateProducts(ctx context.Context, products []string) error {
	return nil
}

func (s *PartnerService) ValidatePartnerProducts(ctx context.Context, productIDs []uuid.UUID) error {
	return nil
}
