package availability

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *AvailabilityService) BlockAvailability(ctx context.Context, id uuid.UUID) error {
	// Get existing availability
	availabilityEncx, err := s.availabilityRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.ErrRepositoryNotFound
		}
		return fmt.Errorf("get availability for blocking: %w", err)
	}

	availability, err := domain.DecryptAvailabilityEncx(ctx, s.crypto, availabilityEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("availability", err)
	}

	// Block
	availability.Block()

	updatedAvailabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, s.crypto, availability)
	if err != nil {
		return errs.NewNotEncryptedErr("availability", err)
	}

	// Persist changes
	if err := s.availabilityRepo.Update(ctx, updatedAvailabilityEncx); err != nil {
		return fmt.Errorf("block availability: %w", err)
	}

	return nil
}
