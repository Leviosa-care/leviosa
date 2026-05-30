package partner

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetPartnerByID retrieves a partner by ID with their associated user information.
func (s *PartnerService) GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.PartnerResponse, error) {
	// Get encrypted partner from repository
	partnerEncx, err := s.partnerRepo.GetPartnerByID(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("get partner by ID from repository: %w", err)
	}

	// Decrypt partner
	partner, err := domain.DecryptPartnerEncx(ctx, s.crypto, partnerEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("partner", err)
	}

	// Build complete response with user info
	return partner.ToResponse(), nil
}
