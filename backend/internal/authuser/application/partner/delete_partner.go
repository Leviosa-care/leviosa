package partner

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// DeletePartner deletes a partner by ID.
// This is an admin-only operation that removes the partner profile but does NOT delete the user account.
func (s *PartnerService) DeletePartner(ctx context.Context, partnerID uuid.UUID) error {
	// Verify partner exists and get their user_id
	partnerEncx, err := s.partnerRepo.GetPartnerByID(ctx, partnerID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "partner")
		case errors.Is(err, errs.ErrConnectionFailure):
			return fmt.Errorf("get partner by ID - database connection failed: %w", err)
		case errors.Is(err, errs.ErrContext):
			return err
		case errors.Is(err, errs.ErrDatabase):
			return fmt.Errorf("get partner by ID - database error: %w", err)
		default:
			return fmt.Errorf("get partner by ID: %w", err)
		}
	}

	// Delete partner using user_id (repository expects user_id, not partner_id)
	if err := s.partnerRepo.DeletePartner(ctx, partnerEncx.UserID); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			// Partner was found but deletion failed - shouldn't happen but handle gracefully
			return errs.NewNotFoundErr(err, "partner")
		case errors.Is(err, errs.ErrRepositoryNotDeleted):
			return fmt.Errorf("delete partner from repository failed: %w", err)
		case errors.Is(err, errs.ErrConnectionFailure):
			return fmt.Errorf("delete partner - database connection failed: %w", err)
		case errors.Is(err, errs.ErrContext):
			return err
		case errors.Is(err, errs.ErrDatabase):
			return fmt.Errorf("delete partner - database error: %w", err)
		default:
			return fmt.Errorf("delete partner from repository: %w", err)
		}
	}

	return nil
}
