package couponPayment

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/coupon"
)

func (s *service) CreateCoupon(ctx context.Context, req *domain.CreateCouponRequest) (*domain.Coupon, error) {
	params := &stripe.CouponParams{
		Name: stripe.String(req.Name),
	}

	// Set discount type (percent_off OR amount_off)
	if req.PercentOff != nil {
		params.PercentOff = stripe.Float64(*req.PercentOff)
	} else if req.AmountOff != nil && req.Currency != nil {
		params.AmountOff = stripe.Int64(int64(*req.AmountOff))
		params.Currency = stripe.String(*req.Currency)
	}

	// Set duration
	params.Duration = stripe.String(req.Duration)
	if req.DurationInMonths != nil {
		params.DurationInMonths = stripe.Int64(int64(*req.DurationInMonths))
	}

	// Set redemption limits
	if req.MaxRedemptions != nil {
		params.MaxRedemptions = stripe.Int64(int64(*req.MaxRedemptions))
	}

	// Set expiry
	if req.RedeemBy != nil {
		params.RedeemBy = stripe.Int64(req.RedeemBy.Unix())
	}

	// Set metadata
	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	stripeCoupon, err := coupon.New(params)
	if err != nil {
		return nil, err
	}

	return mapStripeCouponToDomainCoupon(stripeCoupon), nil
}

func mapStripeCouponToDomainCoupon(stripeCoupon *stripe.Coupon) *domain.Coupon {
	domainCoupon := &domain.Coupon{
		StripeCouponID: stripeCoupon.ID,
		Name:           stripeCoupon.Name,
		Duration:       domain.CouponDuration(stripeCoupon.Duration),
		TimesRedeemed:  int(stripeCoupon.TimesRedeemed),
		IsValid:        stripeCoupon.Valid,
		CreatedAt:      time.Unix(stripeCoupon.Created, 0),
	}

	if stripeCoupon.PercentOff > 0 {
		percentOff := float64(stripeCoupon.PercentOff)
		domainCoupon.PercentOff = &percentOff
	}

	if stripeCoupon.AmountOff > 0 {
		amountOff := int(stripeCoupon.AmountOff)
		domainCoupon.AmountOff = &amountOff
		currency := string(stripeCoupon.Currency)
		domainCoupon.Currency = &currency
	}

	if stripeCoupon.DurationInMonths > 0 {
		durationInMonths := int(stripeCoupon.DurationInMonths)
		domainCoupon.DurationInMonths = &durationInMonths
	}

	if stripeCoupon.MaxRedemptions > 0 {
		maxRedemptions := int(stripeCoupon.MaxRedemptions)
		domainCoupon.MaxRedemptions = &maxRedemptions
	}

	if stripeCoupon.RedeemBy > 0 {
		redeemBy := time.Unix(stripeCoupon.RedeemBy, 0)
		domainCoupon.RedeemBy = &redeemBy
	}

	if len(stripeCoupon.Metadata) > 0 {
		domainCoupon.Metadata = stripeCoupon.Metadata
	}

	return domainCoupon
}
