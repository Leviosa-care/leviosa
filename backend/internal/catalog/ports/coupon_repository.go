package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/google/uuid"
)

type CouponRepository interface {
	// reader
	GetCouponByID(ctx context.Context, couponID uuid.UUID) (*domain.Coupon, error)
	GetCouponByStripeID(ctx context.Context, stripeCouponID string) (*domain.Coupon, error)
	GetAllCoupons(ctx context.Context) ([]*domain.Coupon, error)
	GetValidCoupons(ctx context.Context) ([]*domain.Coupon, error)
	CouponExistsByName(ctx context.Context, name string) (bool, error)
	
	// writer
	CreateCoupon(ctx context.Context, coupon *domain.Coupon) (string, error)
	UpdateCoupon(ctx context.Context, couponID uuid.UUID, req *domain.UpdateCouponRequest) error
	IncrementRedemptionCount(ctx context.Context, couponID uuid.UUID) error
	DeactivateCoupon(ctx context.Context, couponID uuid.UUID) error
	DeleteCoupon(ctx context.Context, couponID uuid.UUID) error
}