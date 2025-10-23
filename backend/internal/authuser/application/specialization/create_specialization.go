package specialization

import (
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *SpecializationService) CreateSpecialization(ctx context.Context, request *domain.CreateSpecializationRequest) (*domain.SpecializationResponse, error) {
	// Validate request
	if err := request.Valid(ctx); err != nil {
		return nil, errs.ErrInvalidInput
	}

	// Check if specialization with same name already exists
	_, err := s.repo.GetSpecializationByName(ctx, request.Name)
	if err == nil {
		return nil, errs.ErrUniqueViolation
	}
	if !errors.Is(err, errs.ErrRepositoryNotFound) {
		return nil, err
	}

	// Create specialization entity
	specialization := request.ToSpecialization()

	// Encrypt sensitive fields using the new generated function
	specializationEncx, err := domain.ProcessSpecializationEncx(ctx, s.crypto, specialization)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("specialization during creation", err)
	}

	// Create in database
	if err := s.repo.CreateSpecialization(ctx, specializationEncx); err != nil {
		return nil, err
	}

	// Decrypt for response using the new generated function
	decryptedSpecialization, err := domain.DecryptSpecializationEncx(ctx, s.crypto, specializationEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("specialization for response", err)
	}

	return decryptedSpecialization.ToResponse(), nil
}

