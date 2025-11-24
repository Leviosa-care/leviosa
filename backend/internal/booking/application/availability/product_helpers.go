package availability

import (
	"context"
	"fmt"

	catalogDomain "github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

// getPartnerProducts fetches all published products for a partner
// This is used for availability duration validation
func (s *AvailabilityService) getPartnerProducts(ctx context.Context) ([]*catalogDomain.ProductRes, error) {
	// Note: Currently fetches all published products
	// TODO: Filter by partner when partner-product association is implemented
	products, err := s.productService.GetAllPublishedProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch published products: %w", err)
	}

	return products, nil
}
