package productRepository

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain/product"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (p *repository) AddOffer(ctx context.Context, offer *productService.Offer) error {
	query := `
        INSERT INTO offers (
			id,
			product_id,
            name,
            description,
            encrypted_picture,
            duration,
            price,
            encrypted_price_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
	result, err := p.DB.ExecContext(
		ctx,
		query,
		offer.ID,
		offer.ProductID,
		offer.Name,
		offer.Description,
		offer.Picture,
		offer.Duration,
		offer.Price,
		offer.PriceID,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(err)
		}

	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return rp.NewDatabaseErr(err)
	}
	if rowsAffected == 0 {
		return rp.NewNotCreatedErr(errors.New("no rows affected by insertion statement"), "product type")
	}
	return nil
}
