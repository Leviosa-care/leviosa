package aggregator

import (
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
)

// CategoryImagesService combines data from category and price domains.
type CategoryImagesService struct {
	categoryService ports.CategoryService
	imageService    ports.ImageParentService
}

// NewCategoryAggregatorService creates a new instance of the aggregator service.
func NewCategoryAggregatorService(categoryService ports.CategoryService, imageService ports.ImageParentService) ports.CategoryImagesService {
	return &CategoryImagesService{
		categoryService: categoryService,
		imageService:    imageService,
	}
}