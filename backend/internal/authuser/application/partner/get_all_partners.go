package partner

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// GetAllPartners retrieves all partners with their associated user information.
func (s *PartnerService) GetAllPartners(ctx context.Context) ([]*domain.PartnerResponse, error) {
	// Get all partners from repository
	partnersEncx, err := s.partnerRepo.GetAllPartners(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrConnectionFailure):
			return nil, fmt.Errorf("get all partners - database connection failure: %w", err)
		case errors.Is(err, errs.ErrTooManyConnections):
			return nil, fmt.Errorf("get all partners - too many database connections: %w", err)
		case errors.Is(err, errs.ErrResourceExhausted):
			return nil, fmt.Errorf("get all partners - database resources exhausted: %w", err)
		case errors.Is(err, errs.ErrQueryCancelled):
			return nil, fmt.Errorf("get all partners - query cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure):
			return nil, fmt.Errorf("get all partners - transaction failure: %w", err)
		case errors.Is(err, errs.ErrDeadlock):
			return nil, fmt.Errorf("get all partners - database deadlock: %w", err)
		case errors.Is(err, errs.ErrPermissionDenied):
			return nil, fmt.Errorf("get all partners - permission denied: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, fmt.Errorf("get all partners - invalid input: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return nil, fmt.Errorf("get all partners - database error: %w", err)
		case errors.Is(err, errs.ErrContext):
			return nil, fmt.Errorf("get all partners - context error: %w", err)
		default:
			return nil, fmt.Errorf("get all partners - unexpected error: %w", err)
		}
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
