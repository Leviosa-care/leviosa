package productService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain"
)

func (s *service) UpdateProduct(ctx context.Context, product *Product) error {
	if err := s.repo.ModifyProduct(ctx,
		product,
		map[string]any{"id": product.ID},
	); err != nil {
		switch {
		default:
			return domain.NewUnexpectTypeErr(err)
		}
	}
	return nil
}
