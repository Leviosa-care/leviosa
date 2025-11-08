package partner

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetAllPartnersByProducts retrieves all partners that offer any of the specified products.
func (s *PartnerService) GetAllPartnersByProducts(ctx context.Context, productIDs []string) ([]*domain.PartnerResponse, error) {
	// Validate that at least one product ID is provided
	if len(productIDs) == 0 {
		return nil, errs.NewInvalidValueErr("product_ids must contain at least one product")
	}

	// Parse and validate all product IDs
	productUUIDs := make([]uuid.UUID, 0, len(productIDs))
	for i, productID := range productIDs {
		productUUID, err := uuid.Parse(productID)
		if err != nil {
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("product_ids[%d] %s", i, err.Error()))
		}
		productUUIDs = append(productUUIDs, productUUID)
	}

	// Get all partners from repository for the given products
	partnersEncx, err := s.partnerRepo.GetAllPartnersByProducts(ctx, productUUIDs)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrConnectionFailure):
			return nil, fmt.Errorf("get partners by products - database connection failure: %w", err)
		case errors.Is(err, errs.ErrTooManyConnections):
			return nil, fmt.Errorf("get partners by products - too many database connections: %w", err)
		case errors.Is(err, errs.ErrResourceExhausted):
			return nil, fmt.Errorf("get partners by products - database resources exhausted: %w", err)
		case errors.Is(err, errs.ErrQueryCancelled):
			return nil, fmt.Errorf("get partners by products - query cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure):
			return nil, fmt.Errorf("get partners by products - transaction failure: %w", err)
		case errors.Is(err, errs.ErrDeadlock):
			return nil, fmt.Errorf("get partners by products - database deadlock: %w", err)
		case errors.Is(err, errs.ErrPermissionDenied):
			return nil, fmt.Errorf("get partners by products - permission denied: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, fmt.Errorf("get partners by products - invalid input: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return nil, fmt.Errorf("get partners by products - database error: %w", err)
		case errors.Is(err, errs.ErrContext):
			return nil, fmt.Errorf("get partners by products - context error: %w", err)
		default:
			return nil, fmt.Errorf("get partners by products - unexpected error: %w", err)
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
