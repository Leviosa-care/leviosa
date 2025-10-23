package partner

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

// CreatePartner creates a new partner profile for an existing user during registration.
//
// This method is called by the aggregator service during partner registration flow.
// It assumes the user already exists and has been validated by the caller.
//
// Parameters:
//   - userID: The UUID of the existing user
//   - bio: Partner bio/description (optional)
//   - experience: Partner experience description (optional)
//   - certifications: List of certifications (optional)
//   - categoryIDs: List of catalog category UUIDs the partner offers services for
//   - productIDs: List of catalog product UUIDs the partner offers services for
//
// Returns error if:
//   - Catalog validation fails (invalid category/product IDs)
//   - Database operation fails
//   - Encryption fails
func (s *PartnerService) CreatePartner(ctx context.Context, userID uuid.UUID, bio, experience string, certifications []string, categoryIDs, productIDs []uuid.UUID) error {
	// Validate catalog IDs against cache
	if err := s.verifyCatalogIDs(categoryIDs, productIDs); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	// Create partner entity
	partner := &domain.Partner{
		ID:             uuid.New(),
		UserID:         userID,
		Bio:            bio,
		Experience:     experience,
		Certifications: certifications,
		CategoryIDs:    categoryIDs,
		ProductIDs:     productIDs,
		IsVerified:     false, // Partners start unverified, admin must verify
	}

	// Encrypt partner data
	partnerEncx, err := domain.ProcessPartnerEncx(ctx, s.crypto, partner)
	if err != nil {
		return errs.NewNotEncryptedErr("partner during creation", err)
	}

	// Create partner in database
	if err := s.partnerRepo.CreatePartner(ctx, partnerEncx); err != nil {
		return err
	}

	return nil
}

// verifyCatalogIDs validates that all provided category and product IDs exist in the catalog cache.
//
// This ensures partners can only be associated with published catalog items that are currently available.
// The catalog cache is kept up-to-date via RabbitMQ events, so this validation reflects real-time catalog state.
func (s *PartnerService) verifyCatalogIDs(categoryIDs, productIDs []uuid.UUID) error {
	var errs errsx.Map

	// Validate all category IDs exist in catalog cache
	for _, categoryID := range categoryIDs {
		if categoryID == uuid.Nil {
			errs.Set("category_ids", "category ID cannot be nil")
			continue
		}

		if !s.catalogCache.IsValidCategory(categoryID) {
			errs.Set("category_ids", "category ID "+categoryID.String()+" does not exist or is not published")
		}
	}

	// Validate all product IDs exist in catalog cache
	for _, productID := range productIDs {
		if productID == uuid.Nil {
			errs.Set("product_ids", "product ID cannot be nil")
			continue
		}

		if !s.catalogCache.IsValidProduct(productID) {
			errs.Set("product_ids", "product ID "+productID.String()+" does not exist or is not published")
		}
	}

	return errs.AsError()
}
