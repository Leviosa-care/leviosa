package aggregator

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/google/uuid"
)

func (s *ProductAggregatorService) GetAllPublishedProducts(ctx context.Context) ([]*domain.ProductAggregator, error) {
	products, err := s.productService.GetAllPublishedProducts(ctx)
	if err != nil {
		return nil, err
	}
	prices, err := s.priceService.GetAllPrices(ctx)
	if err != nil {
		return nil, err
	}
	pricesByProductID := make(map[uuid.UUID][]*domain.Price)
	for _, price := range prices {
		pricesByProductID[price.ProductID] = append(pricesByProductID[price.ProductID], price)
	}

	images, err := s.imageService.GetAllActiveImages(ctx, string(domain.ProductType))
	if err != nil {
		return nil, err
	}
	imagesByParentID := make(map[uuid.UUID]*domain.Image)
	for _, image := range images {
		imagesByParentID[image.ParentID] = image
	}

	responseSlice := make([]*domain.ProductAggregator, 0, len(products))
	for _, product := range products {
		responseSlice = append(responseSlice, &domain.ProductAggregator{
			Product: product,
			Prices:  pricesByProductID[product.ID],
			Image:   imagesByParentID[product.ID],
		})
	}
	return responseSlice, nil
}