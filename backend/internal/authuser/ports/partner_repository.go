package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
)

type PartnerRepository interface {
	GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.PartnerEncx, error)
	GetAllPartners(ctx context.Context) ([]*domain.PartnerEncx, error)
	CreatePartner(ctx context.Context, partner *domain.PartnerEncx) error
	UpdatePartner(ctx context.Context, partner *domain.PartnerEncx) error
	DeletePartner(ctx context.Context, userID uuid.UUID) error
	VerifyPartner(ctx context.Context, userID uuid.UUID, verifiedByUserID uuid.UUID) error
}
