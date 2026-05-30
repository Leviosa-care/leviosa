package partner

import (
	"context"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// UpdatePartner updates an existing partner's profile fields.
// Only updates fields that are provided (non-nil) in the request.
func (s *PartnerService) UpdatePartner(ctx context.Context, partnerID uuid.UUID, request *domain.UpdatePartnerRequest) (*domain.PartnerResponse, error) {
	// Validate request
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get existing partner
	partnerEncx, err := s.partnerRepo.GetPartnerByID(ctx, partnerID)
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
	if request.Occupation != nil {
		partner.Occupation = strings.TrimSpace(*request.Occupation)
	}
	if request.Quote != nil {
		partner.Quote = strings.TrimSpace(*request.Quote)
	}
	if request.Tags != nil {
		tags := make([]string, 0, len(*request.Tags))
		for _, tag := range *request.Tags {
			tags = append(tags, strings.TrimSpace(tag))
		}
		partner.Tags = tags
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
		return nil, fmt.Errorf("update partner: %w", err)
	}

	// Return updated partner response
	return partner.ToResponse(), nil
}
