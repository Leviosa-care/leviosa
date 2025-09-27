package specialization

import (
	"context"
	"errors"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
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

	// Encrypt sensitive fields
	if err := s.crypto.Encrypt(ctx, specialization); err != nil {
		return nil, errs.ErrInvalidValue
	}

	// Create in database
	if err := s.repo.CreateSpecialization(ctx, specialization); err != nil {
		return nil, err
	}

	// Decrypt for response
	if err := s.crypto.Decrypt(ctx, specialization); err != nil {
		return nil, errs.ErrInvalidValue
	}

	return specialization.ToResponse(), nil
}