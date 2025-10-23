package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
)

type PartnerService interface {
	GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.CompletePartnerResponse, error)
	GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.CompletePartnerResponse, error)
	GetAllPartners(ctx context.Context) (*domain.GetPartnersResponse, error)
	GetAllPartnersByCategory(ctx context.Context, category string) (*domain.GetPartnersResponse, error) // fetch the cache to get the cate
	GetAllPartnersByCategories(ctx context.Context, category []string) (*domain.GetPartnersResponse, error)
	CreatePartner(ctx context.Context, userID uuid.UUID, bio, experience string, certifications []string, categoryIDs, productIDs []uuid.UUID) error
	UpdatePartner(ctx context.Context, partnerID uuid.UUID, request *domain.UpdatePartnerRequest) (*domain.PartnerResponse, error)
	DeletePartner(ctx context.Context, partnerID uuid.UUID) error
	VerifyPartner(ctx context.Context, partnerID uuid.UUID, verifiedByUserID uuid.UUID) (*domain.PartnerResponse, error)

	UpdateCategories(ctx context.Context, categories []string) error
	UpdateProducts(ctx context.Context, products []string) error
}

// NOTE: the old thing. When I create a partner, the user already exists, I am just adding the informations that are not tracked in user and that makes important the fact to add a new user
// type PartnerService interface {
// 	GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.CompletePartnerResponse, error)
// 	GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.CompletePartnerResponse, error)
// 	GetAllPartners(ctx context.Context) (*domain.GetPartnersResponse, error)
// 	CreatePartner(ctx context.Context, request *domain.CreatePartnerRequest) (*domain.CompletePartnerResponse, error)
// 	UpdatePartner(ctx context.Context, partnerID uuid.UUID, request *domain.UpdatePartnerRequest) (*domain.PartnerResponse, error)
// 	DeletePartner(ctx context.Context, partnerID uuid.UUID) error
// 	VerifyPartner(ctx context.Context, partnerID uuid.UUID, verifiedByUserID uuid.UUID) (*domain.PartnerResponse, error)
//
// 	// Partner specializations management
// 	AddPartnerSpecialization(ctx context.Context, partnerID, specializationID uuid.UUID) error
// 	RemovePartnerSpecialization(ctx context.Context, partnerID, specializationID uuid.UUID) error
// 	GetPartnerSpecializations(ctx context.Context, partnerID uuid.UUID) (*domain.GetPartnerSpecializationsResponse, error)
//
// 	// Catalog validation methods
// 	ValidatePartnerSpecializations(ctx context.Context, specializationIDs []uuid.UUID) error
// 	ValidatePartnerProducts(ctx context.Context, productIDs []uuid.UUID) error
// }
