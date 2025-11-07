package partner

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetAllPartnersByCategory retrieves all partners that offer services for a specific category.
func (s *PartnerService) GetAllPartnersByCategory(ctx context.Context, categoryID string) ([]*domain.PartnerResponse, error) {
	// Parse and validate category ID
	categoryUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("category_id %s", err.Error()))
	}

	// Get all partners from repository for the given category
	partnersEncx, err := s.partnerRepo.GetAllPartnersByCategory(ctx, categoryUUID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrConnectionFailure):
			return nil, fmt.Errorf("get partners by category - database connection failure: %w", err)
		case errors.Is(err, errs.ErrTooManyConnections):
			return nil, fmt.Errorf("get partners by category - too many database connections: %w", err)
		case errors.Is(err, errs.ErrResourceExhausted):
			return nil, fmt.Errorf("get partners by category - database resources exhausted: %w", err)
		case errors.Is(err, errs.ErrQueryCancelled):
			return nil, fmt.Errorf("get partners by category - query cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure):
			return nil, fmt.Errorf("get partners by category - transaction failure: %w", err)
		case errors.Is(err, errs.ErrDeadlock):
			return nil, fmt.Errorf("get partners by category - database deadlock: %w", err)
		case errors.Is(err, errs.ErrPermissionDenied):
			return nil, fmt.Errorf("get partners by category - permission denied: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, fmt.Errorf("get partners by category - invalid input: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return nil, fmt.Errorf("get partners by category - database error: %w", err)
		case errors.Is(err, errs.ErrContext):
			return nil, fmt.Errorf("get partners by category - context error: %w", err)
		default:
			return nil, fmt.Errorf("get partners by category - unexpected error: %w", err)
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
