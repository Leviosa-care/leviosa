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
	GetAllPartnersByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.PartnerEncx, error)
	GetAllPartnersByCategories(ctx context.Context, categoryIDs []uuid.UUID) ([]*domain.PartnerEncx, error)
	GetAllPartnersByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.PartnerEncx, error)
	CreatePartner(ctx context.Context, partner *domain.PartnerEncx) error
	UpdatePartner(ctx context.Context, partner *domain.PartnerEncx) error
	DeletePartner(ctx context.Context, userID uuid.UUID) error
	VerifyPartner(ctx context.Context, userID uuid.UUID, verifiedByUserID uuid.UUID) error
}
