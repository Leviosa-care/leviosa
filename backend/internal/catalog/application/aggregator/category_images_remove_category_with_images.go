package aggregator

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

func (s *CategoryImagesService) RemoveCategoryWithImages(ctx context.Context, categoryID string) error {
	if err := s.imageService.DeleteImages(ctx, categoryID, string(domain.CategoryType)); err != nil {
		return err
	}

	if err := s.categoryService.RemoveCategory(ctx, categoryID); err != nil {
		return err
	}
	return nil
}