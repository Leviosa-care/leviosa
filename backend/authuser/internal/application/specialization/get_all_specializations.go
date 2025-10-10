package specialization

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
)

func (s *SpecializationService) GetAllSpecializations(ctx context.Context) (*domain.GetSpecializationsResponse, error) {
	// Get from database
	specializationsEncx, err := s.repo.GetAllSpecializations(ctx)
	if err != nil {
		return nil, err
	}

	// Decrypt all specializations using the new generated function
	responses := make([]domain.SpecializationResponse, 0, len(specializationsEncx))
	for _, specEncx := range specializationsEncx {
		spec, err := domain.DecryptSpecializationEncx(ctx, s.crypto, specEncx)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *spec.ToResponse())
	}

	return &domain.GetSpecializationsResponse{
		Specializations: responses,
		Total:           len(responses),
	}, nil
}