package productService

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) RemoveProduct(ctx context.Context, productID string) error {
	if err := s.repo.RemoveProduct(ctx, productID); err != nil {
		// TODO: do all the error handling for that thing
		switch {
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrContext):
			return err
		default:
			return domain.NewUnexpectTypeErr(err)
		}
	}
	return nil
}
