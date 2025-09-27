package specialization

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
)

func (s *SpecializationService) GetActiveSpecializations(ctx context.Context) (*domain.GetSpecializationsResponse, error) {
	// Get from database
	specializations, err := s.repo.GetActiveSpecializations(ctx)
	if err != nil {
		return nil, err
	}

	// Decrypt all specializations
	responses := make([]domain.SpecializationResponse, 0, len(specializations))
	for _, spec := range specializations {
		if err := s.crypto.Decrypt(ctx, spec); err != nil {
			return nil, err
		}
		responses = append(responses, *spec.ToResponse())
	}

	return &domain.GetSpecializationsResponse{
		Specializations: responses,
		Total:           len(responses),
	}, nil
}