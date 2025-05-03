package productService

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) GetProduct(ctx context.Context, productID string) (*Product, error) {
	// check in here if the produdctID is ""
	product, err := s.repo.GetProduct(ctx, productID)
	if err != nil {
		// TODO: add other errors
		switch {
		case errors.Is(err, rp.ErrNotFound):
			return nil, domain.NewNotFoundErr(err)
		default:
			return nil, domain.NewUnexpectTypeErr(err)
		}
	}
	return product, nil
}
