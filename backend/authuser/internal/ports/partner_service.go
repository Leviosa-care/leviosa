package ports

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
)

type PartnerService interface {
	GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.CompletePartnerResponse, error)
	GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.CompletePartnerResponse, error)
	GetAllPartners(ctx context.Context) (*domain.GetPartnersResponse, error)
	CreatePartner(ctx context.Context, request *domain.CreatePartnerRequest) (*domain.CompletePartnerResponse, error)
	UpdatePartner(ctx context.Context, partnerID uuid.UUID, request *domain.UpdatePartnerRequest) (*domain.PartnerResponse, error)
	DeletePartner(ctx context.Context, partnerID uuid.UUID) error
	VerifyPartner(ctx context.Context, partnerID uuid.UUID, verifiedByUserID uuid.UUID) (*domain.PartnerResponse, error)

	// Partner specializations management
	AddPartnerSpecialization(ctx context.Context, partnerID, specializationID uuid.UUID) error
	RemovePartnerSpecialization(ctx context.Context, partnerID, specializationID uuid.UUID) error
	GetPartnerSpecializations(ctx context.Context, partnerID uuid.UUID) (*domain.GetPartnerSpecializationsResponse, error)
}