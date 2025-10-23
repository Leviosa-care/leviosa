package aggregator

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
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

func (s *CategoryImagesService) RemoveCategoryWithImages(ctx context.Context, categoryID string) error {
	if err := s.imageService.DeleteImages(ctx, categoryID, string(domain.CategoryType)); err != nil {
		return err
	}

	if err := s.categoryService.RemoveCategory(ctx, categoryID); err != nil {
		return err
	}
	return nil
}

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

// GetAllPublishedCategoriesWithImages fetches all published categories and their associated images,
// efficiently combining data from separate domain services.
func (s *CategoryImagesService) GetAllPublishedCategoriesWithImages(ctx context.Context) ([]*domain.CategoryWithImage, error) {
	// Step 1: Get all categories from the category service.
	categories, err := s.categoryService.GetAllPublishedCategories(ctx)
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

	imagesByOwnerID := make(map[uuid.UUID]*domain.Image)
	for _, image := range allImages {
		imagesByOwnerID[image.ParentID] = image
	}

	responseSlice := make([]*domain.CategoryWithImage, 0, len(categories))
	for _, category := range categories {
		responseSlice = append(responseSlice, &domain.CategoryWithImage{
			Category: category,
			Image:    imagesByOwnerID[category.ID], // Can be nil if no image was found
		})
	}

	return responseSlice, nil
}
