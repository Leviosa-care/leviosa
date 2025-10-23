package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	// reader
	GetAllCategories(ctx context.Context) ([]*domain.Category, error)
	GetAllPublishedCategories(ctx context.Context) ([]*domain.Category, error)
	CategoryExistsByName(ctx context.Context, name string) (bool, error)
	CountProductsInCategory(ctx context.Context, categoryID uuid.UUID) (int, error)
	// writer
	AddCategory(ctx context.Context, c *domain.Category) (string, error)
	UpdateCategory(ctx context.Context, category *domain.UpdateCategoryRequest) error
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
}
