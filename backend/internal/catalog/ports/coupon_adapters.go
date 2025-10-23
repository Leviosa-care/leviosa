package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type CouponPaymentGateway interface {
	// Coupon management
	CreateCoupon(ctx context.Context, req *domain.CreateCouponRequest) (*domain.Coupon, error)
	GetCoupon(ctx context.Context, stripeCouponID string) (*domain.Coupon, error)
	UpdateCoupon(ctx context.Context, stripeCouponID string, req *domain.UpdateCouponRequest) (*domain.Coupon, error)
	DeleteCoupon(ctx context.Context, stripeCouponID string) error
}
