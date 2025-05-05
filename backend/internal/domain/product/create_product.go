package productService

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"

	"github.com/google/uuid"
)

func (s *service) CreateProduct(ctx context.Context, product *Product) error {
	if err := product.Valid(ctx); err != nil {
		return domain.NewInvalidValueErr(fmt.Sprintf("product validation error: %s", err.Error()))
	}
	product.ID = uuid.NewString()
	if err := s.repo.AddProduct(ctx, product); err != nil {
		switch {
		case errors.Is(err, rp.ErrNotCreated):
			return domain.NewNotCreatedErr(err)
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrContext):
			return err
		}
	}
	return nil
}
