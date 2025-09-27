package ports

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
)

type PartnerRepository interface {
	GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.Partner, error)
	GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.Partner, error)
	GetAllPartners(ctx context.Context) ([]*domain.Partner, error)
	GetPartnerWithUser(ctx context.Context, partnerID uuid.UUID) (*domain.Partner, error)
	GetPartnersWithUsers(ctx context.Context) ([]*domain.Partner, error)
	CreatePartner(ctx context.Context, partner *domain.Partner) error
	UpdatePartner(ctx context.Context, partner *domain.Partner) error
	DeletePartner(ctx context.Context, partnerID uuid.UUID) error
	VerifyPartner(ctx context.Context, partnerID uuid.UUID, verifiedByUserID uuid.UUID) error

	// Partner specializations management
	AddPartnerSpecialization(ctx context.Context, partnerID, specializationID uuid.UUID) error
	RemovePartnerSpecialization(ctx context.Context, partnerID, specializationID uuid.UUID) error
	GetPartnerSpecializations(ctx context.Context, partnerID uuid.UUID) ([]*domain.Specialization, error)
}