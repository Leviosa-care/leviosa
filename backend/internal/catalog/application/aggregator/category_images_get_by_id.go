package aggregator

import (
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *CategoryImagesService) GetCategoryByIDWithImage(ctx context.Context, categoryID string) (*domain.CategoryWithImage, error) {
	image, err := s.imageService.GetActiveImage(ctx, categoryID, string(domain.CategoryType))
	if err != nil {
		if !errors.Is(err, errs.ErrDomainNotFound) {
			return nil, err
		}
		image = nil
	}

	category, err := s.categoryService.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	return &domain.CategoryWithImage{
		Category: category,
		Image:    image,
	}, nil
}
