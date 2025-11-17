package aggregator

import (
	"context"
	"fmt"
	"log"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/google/uuid"
)

// GetAdminAllCategoriesWithImages fetches all categories and their associated images,
// efficiently combining data from separate domain services.
func (s *CategoryImagesService) GetAdminAllCategoriesWithImages(ctx context.Context) ([]*domain.CategoryWithImage, error) {
	// Step 1: Get all categories from the category service.
	categories, err := s.categoryService.GetAllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("aggregator: failed to get all categories: %w", err)
	}

	if len(categories) == 0 {
		return []*domain.CategoryWithImage{}, nil
	}

	allImages, err := s.imageService.GetAllActiveImages(ctx, string(domain.CategoryType))
	if err != nil {
		log.Printf("aggregator: failed to get all active images: %v", err)
		allImages = []*domain.Image{}
	}

	imagesByParentID := make(map[uuid.UUID]*domain.Image)
	for _, image := range allImages {
		imagesByParentID[image.ParentID] = image
	}

	responseSlice := make([]*domain.CategoryWithImage, 0, len(categories))
	for _, category := range categories {
		responseSlice = append(responseSlice, &domain.CategoryWithImage{
			Category: category,
			Image:    imagesByParentID[category.ID], // Can be nil if no image was found
		})
	}

	return responseSlice, nil
}
