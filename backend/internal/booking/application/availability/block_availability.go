package availability

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *AvailabilityService) BlockAvailability(ctx context.Context, id uuid.UUID) error {
	// Get existing availability
	availability, err := s.availabilityRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.ErrRepositoryNotFound
		}
		return fmt.Errorf("get availability for blocking: %w", err)
	}

	// Block
	availability.Block()

	// Persist changes
	if err := s.availabilityRepo.Update(ctx, availability); err != nil {
		return fmt.Errorf("block availability: %w", err)
	}

	return nil
}