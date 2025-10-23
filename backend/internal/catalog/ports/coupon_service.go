package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type CouponService interface {
	// Core CRUD operations
	CreateCoupon(ctx context.Context, req *domain.CreateCouponRequest) (string, error)
	GetCouponByID(ctx context.Context, couponID string) (*domain.CouponResponse, error)
	GetCouponByStripeID(ctx context.Context, stripeCouponID string) (*domain.CouponResponse, error)
	GetAllCoupons(ctx context.Context) ([]*domain.CouponResponse, error)
	GetValidCoupons(ctx context.Context) ([]*domain.CouponResponse, error)
	UpdateCoupon(ctx context.Context, couponID string, req *domain.UpdateCouponRequest) error
	DeactivateCoupon(ctx context.Context, couponID string) error
	DeleteCoupon(ctx context.Context, couponID string) error
	
	// Business operations
	ValidateCoupon(ctx context.Context, stripeCouponID string) (*domain.CouponResponse, error)
	IncrementRedemptionCount(ctx context.Context, couponID string) error
	CheckRedemptionLimit(ctx context.Context, couponID string) (bool, error)
}