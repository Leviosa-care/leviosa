package partner

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetAllPartnersByProduct retrieves all partners that offer a specific product.
func (s *PartnerService) GetAllPartnersByProduct(ctx context.Context, productID string) ([]*domain.PartnerResponse, error) {
	// Parse and validate product ID
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("product_id %s", err.Error()))
	}

	// Get all partners from repository for the given product
	partnersEncx, err := s.partnerRepo.GetAllPartnersByProduct(ctx, productUUID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrConnectionFailure):
			return nil, fmt.Errorf("get partners by product - database connection failure: %w", err)
		case errors.Is(err, errs.ErrTooManyConnections):
			return nil, fmt.Errorf("get partners by product - too many database connections: %w", err)
		case errors.Is(err, errs.ErrResourceExhausted):
			return nil, fmt.Errorf("get partners by product - database resources exhausted: %w", err)
		case errors.Is(err, errs.ErrQueryCancelled):
			return nil, fmt.Errorf("get partners by product - query cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure):
			return nil, fmt.Errorf("get partners by product - transaction failure: %w", err)
		case errors.Is(err, errs.ErrDeadlock):
			return nil, fmt.Errorf("get partners by product - database deadlock: %w", err)
		case errors.Is(err, errs.ErrPermissionDenied):
			return nil, fmt.Errorf("get partners by product - permission denied: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, fmt.Errorf("get partners by product - invalid input: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return nil, fmt.Errorf("get partners by product - database error: %w", err)
		case errors.Is(err, errs.ErrContext):
			return nil, fmt.Errorf("get partners by product - context error: %w", err)
		default:
			return nil, fmt.Errorf("get partners by product - unexpected error: %w", err)
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
