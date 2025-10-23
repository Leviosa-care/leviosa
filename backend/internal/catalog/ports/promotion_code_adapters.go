package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type PromotionCodePaymentGateway interface {
	// Promotion Code management
	CreatePromotionCode(ctx context.Context, req *domain.CreatePromotionCodeRequest) (*domain.PromotionCode, error)
	GetPromotionCode(ctx context.Context, stripePromotionID string) (*domain.PromotionCode, error)
	UpdatePromotionCode(ctx context.Context, stripePromotionID string, req *domain.UpdatePromotionCodeRequest) (*domain.PromotionCode, error)
	DeletePromotionCode(ctx context.Context, stripePromotionID string) error
	ListPromotionCodes(ctx context.Context, stripeCouponID string) ([]*domain.PromotionCode, error)
}