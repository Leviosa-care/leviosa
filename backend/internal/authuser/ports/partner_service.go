package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
)

type PartnerService interface {
	GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.PartnerResponse, error)
	GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.PartnerResponse, error)
	GetAllPartners(ctx context.Context) (*domain.GetPartnersResponse, error)
	GetAllPartnersByCategory(ctx context.Context, category string) (*domain.GetPartnersResponse, error) // fetch the cache to get the cate
	GetAllPartnersByCategories(ctx context.Context, category []string) (*domain.GetPartnersResponse, error)
	CreatePartner(ctx context.Context, userID uuid.UUID, bio, experience string, categoryIDs, productIDs []uuid.UUID) (*domain.Partner, error)
	UpdatePartner(ctx context.Context, partnerID uuid.UUID, request *domain.UpdatePartnerRequest) (*domain.PartnerResponse, error)
	DeletePartner(ctx context.Context, partnerID uuid.UUID) error
	VerifyPartner(ctx context.Context, partnerID uuid.UUID, verifiedByUserID uuid.UUID) (*domain.PartnerResponse, error)

	UpdateCategories(ctx context.Context, categories []string) error
	UpdateProducts(ctx context.Context, products []string) error

	ValidatePartnerProducts(ctx context.Context, productIDs []uuid.UUID) error
}
