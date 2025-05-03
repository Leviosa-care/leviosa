package productService

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"

	"github.com/google/uuid"
)

func (s *service) CreateOffer(ctx context.Context, offer *Offer) error {
	if errs := offer.Valid(ctx); len(errs) > 0 {
		return domain.NewInvalidValueErr(fmt.Sprintf("offer validation error: %s", errs.Error()))
	}
	offer.ID = uuid.NewString()
	if err := s.repo.AddOffer(ctx, offer); err != nil {
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
