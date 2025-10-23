package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/google/uuid"
)

type PromotionCodeRepository interface {
	// reader
	GetPromotionCodeByID(ctx context.Context, promotionCodeID uuid.UUID) (*domain.PromotionCode, error)
	GetPromotionCodeByCode(ctx context.Context, code string) (*domain.PromotionCode, error)
	GetPromotionCodeByStripeID(ctx context.Context, stripePromotionID string) (*domain.PromotionCode, error)
	GetPromotionCodesByCouponID(ctx context.Context, couponID uuid.UUID) ([]*domain.PromotionCode, error)
	GetAllPromotionCodes(ctx context.Context) ([]*domain.PromotionCode, error)
	GetActivePromotionCodes(ctx context.Context) ([]*domain.PromotionCode, error)
	PromotionCodeExistsByCode(ctx context.Context, code string) (bool, error)
	
	// writer
	CreatePromotionCode(ctx context.Context, promotionCode *domain.PromotionCode) (string, error)
	UpdatePromotionCode(ctx context.Context, promotionCodeID uuid.UUID, req *domain.UpdatePromotionCodeRequest) error
	IncrementRedemptionCount(ctx context.Context, promotionCodeID uuid.UUID) error
	DeactivatePromotionCode(ctx context.Context, promotionCodeID uuid.UUID) error
	DeletePromotionCode(ctx context.Context, promotionCodeID uuid.UUID) error
	DeletePromotionCodesByCouponID(ctx context.Context, couponID uuid.UUID) error
}