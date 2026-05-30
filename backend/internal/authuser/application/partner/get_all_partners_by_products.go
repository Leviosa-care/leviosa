package partner

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetAllPartnersByProducts retrieves all partners that offer any of the specified products.
func (s *PartnerService) GetAllPartnersByProducts(ctx context.Context, productIDs []string) ([]*domain.PartnerResponse, error) {
	// Validate that at least one product ID is provided
	if len(productIDs) == 0 {
		return nil, errs.NewInvalidValueErr("product_ids must contain at least one product")
	}

	// Parse and validate all product IDs
	productUUIDs := make([]uuid.UUID, 0, len(productIDs))
	for i, productID := range productIDs {
		productUUID, err := uuid.Parse(productID)
		if err != nil {
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("product_ids[%d] %s", i, err.Error()))
		}
		productUUIDs = append(productUUIDs, productUUID)
	}

	// Get all partners from repository for the given products
	partnersEncx, err := s.partnerRepo.GetAllPartnersByProducts(ctx, productUUIDs)
	if err != nil {
		return nil, fmt.Errorf("get partners by products: %w", err)
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
		partners = append(partners, partner.ToResponse())
	}

	return partners, nil
}
