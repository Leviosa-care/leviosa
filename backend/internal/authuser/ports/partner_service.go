package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
)

type PartnerService interface {
	GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.PartnerResponse, error)
	GetAllPartners(ctx context.Context) ([]*domain.PartnerResponse, error)
	GetAllPartnersByCategory(ctx context.Context, categoryID string) ([]*domain.PartnerResponse, error) // fetch the cache to get the cate
	GetAllPartnersByCategories(ctx context.Context, categoryIDs []string) ([]*domain.PartnerResponse, error)
	CreatePartner(ctx context.Context, userID uuid.UUID, bio, experience string, categoryIDs, productIDs []uuid.UUID) (*domain.Partner, error)
	UpdatePartner(ctx context.Context, partnerID uuid.UUID, request *domain.UpdatePartnerRequest) (*domain.PartnerResponse, error)
	DeletePartner(ctx context.Context, partnerID uuid.UUID) error
	VerifyPartner(ctx context.Context, partnerID uuid.UUID, verifiedByUserID uuid.UUID) (*domain.PartnerResponse, error)
}
