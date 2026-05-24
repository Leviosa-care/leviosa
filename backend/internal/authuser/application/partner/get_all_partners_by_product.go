package partner

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetAllPartnersByProduct retrieves all partners that offer a specific product.
func (s *PartnerService) GetAllPartnersByProduct(ctx context.Context, productID string) ([]*domain.PartnerResponse, error) {
	// Parse and validate product ID
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("product_id %s", err.Error()))
	}

	// Get all partners from repository for the given product
	partnersEncx, err := s.partnerRepo.GetAllPartnersByProduct(ctx, productUUID)
	if err != nil {
		return nil, fmt.Errorf("get partners by product: %w", err)
	}

	// Decrypt partners and build response
	partners := make([]*domain.PartnerResponse, 0, len(partnersEncx))
	for _, partnerEncx := range partnersEncx {
		// Decrypt partner
		partner, err := domain.DecryptPartnerEncx(ctx, s.crypto, partnerEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("partner", err)
		}

		// Build complete partner response
		partners = append(partners, &domain.PartnerResponse{
			ID:                      partner.ID,
			UserID:                  partner.UserID,
			Bio:                     partner.Bio,
			Experience:              partner.Experience,
			// Certifications: partner.Certifications,
			CategoryIDs:             partner.CategoryIDs,
			ProductIDs:              partner.ProductIDs,
			StripeAccountStatus:     partner.StripeAccountStatus,
			StripeOnboardingComplete: partner.StripeOnboardingComplete,
			CreatedAt:               partner.CreatedAt,
			UpdatedAt:               partner.UpdatedAt,
		})
	}

	return partners, nil
}
