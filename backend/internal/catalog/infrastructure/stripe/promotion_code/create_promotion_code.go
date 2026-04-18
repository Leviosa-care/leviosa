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
		Coupon: stripe.String(req.CouponID),
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

	// Set metadata
	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	// Build restrictions (consolidates FirstTimeTransaction, MinimumAmount, CurrencyOptions)
	restrictions := &stripe.PromotionCodeRestrictionsParams{
		FirstTimeTransaction: stripe.Bool(req.FirstTimeTransaction),
	}

	if req.MinimumAmount != nil && req.MinimumAmountCurrency != nil {
		restrictions.MinimumAmount = stripe.Int64(int64(*req.MinimumAmount))
		restrictions.MinimumAmountCurrency = stripe.String(*req.MinimumAmountCurrency)
	}

	if req.Restrictions != nil && len(req.Restrictions.CurrencyOptions) > 0 {
		currencyOptions := make(map[string]*stripe.PromotionCodeRestrictionsCurrencyOptionsParams)
		for _, currency := range req.Restrictions.CurrencyOptions {
			currencyOptions[currency] = &stripe.PromotionCodeRestrictionsCurrencyOptionsParams{}
		}
		restrictions.CurrencyOptions = currencyOptions
	}

	params.Restrictions = restrictions

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

	// Handle restrictions
	if stripePC.Restrictions != nil {
		domainPC.FirstTimeTransaction = stripePC.Restrictions.FirstTimeTransaction

		if stripePC.Restrictions.MinimumAmount > 0 {
			minimumAmount := int(stripePC.Restrictions.MinimumAmount)
			domainPC.MinimumAmount = &minimumAmount
			currency := string(stripePC.Restrictions.MinimumAmountCurrency)
			domainPC.MinimumAmountCurrency = &currency
		}

		if len(stripePC.Restrictions.CurrencyOptions) > 0 {
			currencies := make([]string, 0, len(stripePC.Restrictions.CurrencyOptions))
			for currency := range stripePC.Restrictions.CurrencyOptions {
				currencies = append(currencies, currency)
			}
			domainPC.Restrictions = &domain.PromotionCodeRestrictions{
				CurrencyOptions: currencies,
			}
		}
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

	// Handle metadata
	if len(stripePC.Metadata) > 0 {
		domainPC.Metadata = stripePC.Metadata
	}

	return domainPC
}
