package productRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/product"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

func (p *repository) GetOffer(ctx context.Context, offerID string) (*productService.Offer, error) {
	var offer productService.Offer
	query := fmt.Sprintf(`
        SELECT
			product_id,
            name,
            description,
            encrypted_picture,
            duration,
            price,
            encrypted_price_id
        FROM offers
        WHERE id = $1;`, pg.QualifiedTable(p.schema, "offers"))
	err := p.DB.QueryRowContext(ctx, query, offerID).Scan(
		&offer.ProductID,
		&offer.Name,
		&offer.Description,
		&offer.Picture,
		&offer.Duration,
		&offer.Price,
		&offer.PriceID,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, rp.NewNotFoundErr(err, "offer")
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return nil, rp.NewContextErr(err)
		default:
			return nil, rp.NewDatabaseErr(err)
		}
	}
	offer.ID = offerID
	return &offer, nil
}
