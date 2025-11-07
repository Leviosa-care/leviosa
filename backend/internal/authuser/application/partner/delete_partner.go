package partner

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

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
