package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type CategoryImagesService interface {
	GetAdminAllCategoriesWithImages(ctx context.Context) ([]*domain.CategoryWithImage, error)
	GetAllPublishedCategoriesWithImages(ctx context.Context) ([]*domain.CategoryWithImage, error)
	GetCategoryByIDWithImage(ctx context.Context, categoryID string) (*domain.CategoryWithImage, error)
	RemoveCategoryWithImages(ctx context.Context, categoryID string) error
}
