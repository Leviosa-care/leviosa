package ports

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
)

type SpecializationService interface {
	GetSpecializationByID(ctx context.Context, specializationID uuid.UUID) (*domain.SpecializationResponse, error)
	GetAllSpecializations(ctx context.Context) (*domain.GetSpecializationsResponse, error)
	GetActiveSpecializations(ctx context.Context) (*domain.GetSpecializationsResponse, error)
	CreateSpecialization(ctx context.Context, request *domain.CreateSpecializationRequest) (*domain.SpecializationResponse, error)
	UpdateSpecialization(ctx context.Context, specializationID uuid.UUID, request *domain.UpdateSpecializationRequest) (*domain.SpecializationResponse, error)
	DeleteSpecialization(ctx context.Context, specializationID uuid.UUID) error
}