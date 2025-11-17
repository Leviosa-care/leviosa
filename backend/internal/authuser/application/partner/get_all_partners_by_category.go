package partner

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetAllPartnersByCategory retrieves all partners that offer services for a specific category.
func (s *PartnerService) GetAllPartnersByCategory(ctx context.Context, categoryID string) ([]*domain.PartnerResponse, error) {
	// Parse and validate category ID
	categoryUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("category_id %s", err.Error()))
	}

	// Get all partners from repository for the given category
	partnersEncx, err := s.partnerRepo.GetAllPartnersByCategory(ctx, categoryUUID)
	if err != nil {
		return nil, fmt.Errorf("get partners by category: %w", err)
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
			ID:         partner.ID,
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

	return partners, nil
}
