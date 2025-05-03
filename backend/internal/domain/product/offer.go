package productService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/stripe"

	"github.com/hengadev/errsx"
)

// (standard, premium etc...)
type Offer struct {
	ID          string `json:"id"`
	ProductID   string `json:"product_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Picture     string `json:"picture"`
	Duration    int    `json:"duration"`
	Price       int64  `json:"price"`
	PriceID     string `json:"price_id"`
}

func (o Offer) AssertComparable() {}

func (o Offer) Valid(ctx context.Context) error {
	var errs errsx.Map
	if o.Name != "" {
		errs.Set("name", "cannot have an empty name")
	}
	if o.Description != "" {
		errs.Set("name", "cannot have an empty description")
	}
	// TODO: Not sure about that one
	if o.Picture != "" {
		errs.Set("name", "cannot have an empty description")
	}
	if o.Price <= 0 {
		errs.Set("price", "price value should be >= 0")
	}
	return errs.AsError()
}

func (o Offer) GetSQLColumnMapping() map[string]string {
	return map[string]string{
		"ID":          "id",
		"ProductID":   "product_id",
		"Name":        "name",
		"Description": "description",
		"Picture":     "picture",
		"Duration":    "duration",
		"Price":       "price",
		"PriceID":     "encrypted_price_id",
	}
}

func (o Offer) GetProhibitedFields() []string {
	return []string{
		"ID",
		"PriceID",
	}
}

// a helper function that is used to get information when proceeding on payment services
func (p *Offer) GetPaymentInfo(ctx context.Context) *stripeService.PaymentProductInfo {
	return &stripeService.PaymentProductInfo{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}
}

// TODO:
// - create offer (service to create object + store in database)
// - get the *offer
// - then get the Payment info from the *offer
// - use the stripe service to create a new stripe offer with *offer and get the stripe PriceID
// - store the priceID in the database
