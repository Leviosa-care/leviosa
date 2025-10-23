package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type PromotionCodeService interface {
	// Core CRUD operations
	CreatePromotionCode(ctx context.Context, req *domain.CreatePromotionCodeRequest) (string, error)
	GetPromotionCodeByID(ctx context.Context, promotionCodeID string) (*domain.PromotionCode, error)
	GetPromotionCodeByCode(ctx context.Context, code string) (*domain.PromotionCode, error)
	GetPromotionCodeByStripeID(ctx context.Context, stripePromotionID string) (*domain.PromotionCode, error)
	GetPromotionCodesByCouponID(ctx context.Context, couponID string) ([]*domain.PromotionCode, error)
	GetAllPromotionCodes(ctx context.Context) ([]*domain.PromotionCode, error)
	GetActivePromotionCodes(ctx context.Context) ([]*domain.PromotionCode, error)
	UpdatePromotionCode(ctx context.Context, promotionCodeID string, req *domain.UpdatePromotionCodeRequest) error
	DeactivatePromotionCode(ctx context.Context, promotionCodeID string) error
	DeletePromotionCode(ctx context.Context, promotionCodeID string) error
	DeletePromotionCodesByCouponID(ctx context.Context, couponID string) error
	
	// Business operations
	ValidatePromotionCode(ctx context.Context, req *domain.ValidatePromotionCodeRequest) (*domain.ValidatePromotionCodeResponse, error)
	IncrementRedemptionCount(ctx context.Context, promotionCodeID string) error
	CheckRedemptionLimit(ctx context.Context, promotionCodeID string) (bool, error)
	GetPromotionCodeWithCoupon(ctx context.Context, code string) (*domain.PromotionCodeWithCouponResponse, error)
}
