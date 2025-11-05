package partner

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	catalogDomain "github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
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
//   - certifications: List of certifications (optional) - REMOVED
//   - categoryIDs: List of catalog category UUIDs the partner offers services for
//   - productIDs: List of catalog product UUIDs the partner offers services for
//
// Behavior:
//   - Invalid category/product IDs are silently filtered out
//   - Only valid IDs that exist in published catalog are stored
//   - Empty arrays are allowed (partner has no products/categories initially)
//
// Returns error if:
//   - Catalog service call fails
//   - Database operation fails
//   - Encryption fails
func (s *PartnerService) CreatePartner(ctx context.Context, userID uuid.UUID, bio, experience string, categoryIDs, productIDs []uuid.UUID) (*domain.Partner, error) {
	// Filter catalog IDs to only include valid published items
	validCategoryIDs, validProductIDs, err := s.verifyCatalogIDs(ctx, categoryIDs, productIDs)
	if err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	now := time.Now()

	// Create partner entity with filtered valid IDs
	partner := &domain.Partner{
		ID:          uuid.New(),
		UserID:      userID,
		Bio:         bio,
		Experience:  experience,
		CategoryIDs: validCategoryIDs,
		ProductIDs:  validProductIDs,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// TODO: do the stripe related account creation operations

	// Encrypt partner data
	partnerEncx, err := domain.ProcessPartnerEncx(ctx, s.crypto, partner)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("partner during creation", err)
	}

	// Create partner in database
	if err := s.partnerRepo.CreatePartner(ctx, partnerEncx); err != nil {
		return nil, err
	}

	return partner, nil
}

// verifyCatalogIDs filters provided category and product IDs to only include those that exist in the published catalog.
//
// This ensures partners can only be associated with published catalog items that are currently available.
// Invalid IDs are silently filtered out rather than causing errors.
//
// Returns:
//   - validCategoryIDs: Subset of input categoryIDs that exist in published catalog
//   - validProductIDs: Subset of input productIDs that exist in published catalog
//   - error: Only if catalog service calls fail
func (s *PartnerService) verifyCatalogIDs(ctx context.Context, categoryIDs, productIDs []uuid.UUID) ([]uuid.UUID, []uuid.UUID, error) {
	var (
		retrievedCategories []*catalogDomain.Category
		retrievedProducts   []*catalogDomain.ProductRes
	)

	g, ctx := errgroup.WithContext(ctx)

	// Fetch categories concurrently
	g.Go(func() error {
		cats, err := s.categoryService.GetAllPublishedCategories(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch published categories: %w", err)
		}
		retrievedCategories = cats
		return nil
	})

	// Fetch products concurrently
	g.Go(func() error {
		prods, err := s.productService.GetAllPublishedProducts(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch published products: %w", err)
		}
		retrievedProducts = prods
		return nil
	})

	// Wait for both goroutines to finish
	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	// Build category set and filter valid IDs
	var validCategoryIDs []uuid.UUID
	if len(categoryIDs) > 0 && len(retrievedCategories) > 0 {
		categorySet := make(map[uuid.UUID]struct{}, len(retrievedCategories))
		for _, c := range retrievedCategories {
			categorySet[c.ID] = struct{}{}
		}
		validCategoryIDs = make([]uuid.UUID, 0, len(categoryIDs))
		for _, id := range categoryIDs {
			if _, exists := categorySet[id]; exists {
				validCategoryIDs = append(validCategoryIDs, id)
			}
		}
	}

	// Build product set and filter valid IDs
	var validProductIDs []uuid.UUID
	if len(productIDs) > 0 && len(retrievedProducts) > 0 {
		productSet := make(map[uuid.UUID]struct{}, len(retrievedProducts))
		for _, p := range retrievedProducts {
			productSet[p.ID] = struct{}{}
		}
		validProductIDs = make([]uuid.UUID, 0, len(productIDs))
		for _, id := range productIDs {
			if _, exists := productSet[id]; exists {
				validProductIDs = append(validProductIDs, id)
			}
		}
	}

	return validCategoryIDs, validProductIDs, nil
}
