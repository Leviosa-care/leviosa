package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, c *domain.CreateCategoryRequest) (string, error)
	GetCategoryByID(ctx context.Context, ID string) (*domain.Category, error)
	GetAllPublishedCategories(ctx context.Context) ([]*domain.Category, error)
	GetAllCategories(ctx context.Context) ([]*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.UpdateCategoryRequest) error
	RemoveCategory(ctx context.Context, ID string) error
}
