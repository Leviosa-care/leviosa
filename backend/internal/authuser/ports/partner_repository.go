package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
)

type PartnerRepository interface {
	GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.PartnerEncx, error)
	GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.PartnerEncx, error)
	GetAllPartners(ctx context.Context) ([]*domain.PartnerEncx, error)
	GetPartnerWithUser(ctx context.Context, partnerID uuid.UUID) (*domain.PartnerEncx, error)
	GetPartnersWithUsers(ctx context.Context) ([]*domain.PartnerEncx, error)
	CreatePartner(ctx context.Context, partner *domain.PartnerEncx) error
	UpdatePartner(ctx context.Context, partner *domain.PartnerEncx) error
	DeletePartner(ctx context.Context, partnerID uuid.UUID) error
	VerifyPartner(ctx context.Context, partnerID uuid.UUID, verifiedByUserID uuid.UUID) error

	// Partner specializations management
	AddPartnerSpecialization(ctx context.Context, partnerID, specializationID uuid.UUID) error
	RemovePartnerSpecialization(ctx context.Context, partnerID, specializationID uuid.UUID) error
	GetPartnerSpecializations(ctx context.Context, partnerID uuid.UUID) ([]*domain.SpecializationEncx, error)
}

