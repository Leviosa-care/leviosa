package partner

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

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
		ID:         partner.ID,
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
