package partner

import (
	"context"
	"errors"
	"fmt"

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
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "partner")
		case errors.Is(err, errs.ErrConnectionFailure):
			return nil, fmt.Errorf("get partner by ID - database connection failed: %w", err)
		case errors.Is(err, errs.ErrContext):
			return nil, err
		case errors.Is(err, errs.ErrDatabase):
			return nil, fmt.Errorf("get partner by ID - database error: %w", err)
		default:
			return nil, fmt.Errorf("get partner by ID: %w", err)
		}
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
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "partner")
		case errors.Is(err, errs.ErrRepositoryNotUpdated):
			return nil, fmt.Errorf("update partner in repository failed: %w", err)
		case errors.Is(err, errs.ErrConnectionFailure):
			return nil, fmt.Errorf("update partner - database connection failed: %w", err)
		case errors.Is(err, errs.ErrContext):
			return nil, err
		case errors.Is(err, errs.ErrDatabase):
			return nil, fmt.Errorf("update partner - database error: %w", err)
		default:
			return nil, fmt.Errorf("update partner in repository: %w", err)
		}
	}

	// Return updated partner response
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
