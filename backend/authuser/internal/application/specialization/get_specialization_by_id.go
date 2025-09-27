package specialization

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
)

func (s *SpecializationService) GetSpecializationByID(ctx context.Context, specializationID uuid.UUID) (*domain.SpecializationResponse, error) {
	// Get from database
	specialization, err := s.repo.GetSpecializationByID(ctx, specializationID)
	if err != nil {
		return nil, err
	}

	// Decrypt sensitive fields
	if err := s.crypto.Decrypt(ctx, specialization); err != nil {
		return nil, err
	}

	return specialization.ToResponse(), nil
}