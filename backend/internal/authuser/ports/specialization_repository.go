package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
)

type SpecializationRepository interface {
	GetSpecializationByID(ctx context.Context, specializationID uuid.UUID) (*domain.SpecializationEncx, error)
	GetSpecializationByName(ctx context.Context, name string) (*domain.SpecializationEncx, error)
	GetAllSpecializations(ctx context.Context) ([]*domain.SpecializationEncx, error)
	GetActiveSpecializations(ctx context.Context) ([]*domain.SpecializationEncx, error)
	CreateSpecialization(ctx context.Context, specialization *domain.SpecializationEncx) error
	UpdateSpecialization(ctx context.Context, specialization *domain.SpecializationEncx) error
	DeleteSpecialization(ctx context.Context, specializationID uuid.UUID) error
}