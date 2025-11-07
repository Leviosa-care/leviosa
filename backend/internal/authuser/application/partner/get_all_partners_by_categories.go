package partner

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetAllPartnersByCategories retrieves all partners that offer services in any of the specified categories.
func (s *PartnerService) GetAllPartnersByCategories(ctx context.Context, categoryIDs []string) ([]*domain.PartnerResponse, error) {
	// Validate that at least one category ID is provided
	if len(categoryIDs) == 0 {
		return nil, errs.NewInvalidValueErr("category_ids must contain at least one category")
	}

	// Parse and validate all category IDs
	categoryUUIDs := make([]uuid.UUID, 0, len(categoryIDs))
	for i, categoryID := range categoryIDs {
		categoryUUID, err := uuid.Parse(categoryID)
		if err != nil {
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("category_ids[%d] %s", i, err.Error()))
		}
		categoryUUIDs = append(categoryUUIDs, categoryUUID)
	}

	// Get all partners from repository for the given categories
	partnersEncx, err := s.partnerRepo.GetAllPartnersByCategories(ctx, categoryUUIDs)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrConnectionFailure):
			return nil, fmt.Errorf("get partners by categories - database connection failure: %w", err)
		case errors.Is(err, errs.ErrTooManyConnections):
			return nil, fmt.Errorf("get partners by categories - too many database connections: %w", err)
		case errors.Is(err, errs.ErrResourceExhausted):
			return nil, fmt.Errorf("get partners by categories - database resources exhausted: %w", err)
		case errors.Is(err, errs.ErrQueryCancelled):
			return nil, fmt.Errorf("get partners by categories - query cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure):
			return nil, fmt.Errorf("get partners by categories - transaction failure: %w", err)
		case errors.Is(err, errs.ErrDeadlock):
			return nil, fmt.Errorf("get partners by categories - database deadlock: %w", err)
		case errors.Is(err, errs.ErrPermissionDenied):
			return nil, fmt.Errorf("get partners by categories - permission denied: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, fmt.Errorf("get partners by categories - invalid input: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return nil, fmt.Errorf("get partners by categories - database error: %w", err)
		case errors.Is(err, errs.ErrContext):
			return nil, fmt.Errorf("get partners by categories - context error: %w", err)
		default:
			return nil, fmt.Errorf("get partners by categories - unexpected error: %w", err)
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
