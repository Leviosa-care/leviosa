package promotionCodePayment

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/promotioncode"
)

func (s *service) CreatePromotionCode(ctx context.Context, req *domain.CreatePromotionCodeRequest) (*domain.PromotionCode, error) {
	params := &stripe.PromotionCodeParams{
		Coupon: stripe.String(req.CouponID), // Stripe coupon ID
		Code:   stripe.String(req.Code),
	}

	// Set redemption limits
	if req.MaxRedemptions != nil {
		params.MaxRedemptions = stripe.Int64(int64(*req.MaxRedemptions))
	}

	// Set expiry
	if req.ExpiresAt != nil {
		params.ExpiresAt = stripe.Int64(req.ExpiresAt.Unix())
	}

	// Set customer eligibility
	params.CustomerEligibility = &stripe.PromotionCodeCustomerEligibilityParams{
		FirstTimeTransaction: stripe.Bool(req.FirstTimeTransaction),
	}

	// Set minimum amount
	if req.MinimumAmount != nil && req.MinimumAmountCurrency != nil {
		params.MinimumAmount = &stripe.PromotionCodeMinimumAmountParams{
			Amount:   stripe.Int64(int64(*req.MinimumAmount)),
			Currency: stripe.String(*req.MinimumAmountCurrency),
		}
	}

	// Set restrictions
	if req.Restrictions != nil && len(req.Restrictions.CurrencyOptions) > 0 {
		params.Restrictions = &stripe.PromotionCodeRestrictionsParams{
			CurrencyOptions: stripe.StringSlice(req.Restrictions.CurrencyOptions),
		}
	}

	// Set metadata
	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	stripePromotionCode, err := promotioncode.New(params)
	if err != nil {
		return nil, err
	}

	return mapStripePromotionCodeToDomainPromotionCode(stripePromotionCode), nil
}

func mapStripePromotionCodeToDomainPromotionCode(stripePC *stripe.PromotionCode) *domain.PromotionCode {
	domainPC := &domain.PromotionCode{
		StripePromotionID: stripePC.ID,
		Code:              stripePC.Code,
		Active:            stripePC.Active,
		TimesRedeemed:     int(stripePC.TimesRedeemed),
		CreatedAt:         time.Unix(stripePC.Created, 0),
	}

	// Handle first time transaction requirement
	if stripePC.CustomerEligibility != nil {
		domainPC.FirstTimeTransaction = stripePC.CustomerEligibility.FirstTimeTransaction
	}

	// Handle max redemptions
	if stripePC.MaxRedemptions > 0 {
		maxRedemptions := int(stripePC.MaxRedemptions)
		domainPC.MaxRedemptions = &maxRedemptions
	}

	// Handle expires at
	if stripePC.ExpiresAt > 0 {
		expiresAt := time.Unix(stripePC.ExpiresAt, 0)
		domainPC.ExpiresAt = &expiresAt
	}

	// Handle minimum amount
	if stripePC.MinimumAmount != nil {
		minimumAmount := int(stripePC.MinimumAmount.Amount)
		domainPC.MinimumAmount = &minimumAmount
		currency := string(stripePC.MinimumAmount.Currency)
		domainPC.MinimumAmountCurrency = &currency
	}

	// Handle restrictions
	if stripePC.Restrictions != nil && len(stripePC.Restrictions.CurrencyOptions) > 0 {
		domainPC.Restrictions = &domain.PromotionCodeRestrictions{
			CurrencyOptions: stripePC.Restrictions.CurrencyOptions,
		}
	}

	// Handle metadata
	if len(stripePC.Metadata) > 0 {
		domainPC.Metadata = stripePC.Metadata
	}

	return domainPC
}
