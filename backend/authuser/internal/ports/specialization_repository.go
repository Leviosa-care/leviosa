package ports

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
)

type SpecializationRepository interface {
	GetSpecializationByID(ctx context.Context, specializationID uuid.UUID) (*domain.Specialization, error)
	GetSpecializationByName(ctx context.Context, name string) (*domain.Specialization, error)
	GetAllSpecializations(ctx context.Context) ([]*domain.Specialization, error)
	GetActiveSpecializations(ctx context.Context) ([]*domain.Specialization, error)
	CreateSpecialization(ctx context.Context, specialization *domain.Specialization) error
	UpdateSpecialization(ctx context.Context, specialization *domain.Specialization) error
	DeleteSpecialization(ctx context.Context, specializationID uuid.UUID) error
}