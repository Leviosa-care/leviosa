package specialization

import (
	"context"

	"github.com/google/uuid"
)

func (s *SpecializationService) DeleteSpecialization(ctx context.Context, specializationID uuid.UUID) error {
	// Delete from database
	return s.repo.DeleteSpecialization(ctx, specializationID)
}

