package aggregator

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// ProductAggregatorService combines data from category and price domains.
type ProductAggregatorService struct {
	productService ports.ProductService
	priceService   ports.PriceService
	imageService   ports.ImageService
}

// NewProductPricesAggregatorService creates a new instance of the aggregator service.
func NewProductPricesAggregatorService(productService ports.ProductService, priceService ports.PriceService, imageService ports.ImageService) ports.ProductAggregatorService {
	return &ProductAggregatorService{
		productService: productService,
		priceService:   priceService,
		imageService:   imageService,
	}
}

func (s *ProductAggregatorService) GetProductByID(ctx context.Context, productIDStr string) (*domain.ProductAggregator, error) {
	image, err := s.imageService.GetActiveImage(ctx, productIDStr, string(domain.ProductType))
	if err != nil {
		if !errors.Is(err, errs.ErrDomainNotFound) {
			return nil, err
		}
		image = nil
	}

	product, err := s.productService.GetProductByID(ctx, productIDStr)
	if err != nil {
		return nil, err
	}
	prices, err := s.priceService.GetPricesByProductID(ctx, productIDStr)
	if err != nil {
		return nil, err
	}

	return &domain.ProductAggregator{
		Product: product,
		Image:   image,
		Prices:  prices,
	}, nil
}

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

func (s *ProductAggregatorService) GetAdminAllProducts(ctx context.Context) ([]*domain.ProductAggregator, error) {
	products, err := s.productService.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}
	prices, err := s.priceService.GetAllPrices(ctx)
	if err != nil {
		return nil, err
	}
	pricesByProductID := make(map[uuid.UUID][]*domain.Price, len(prices))
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

// CreateProductWithPrice creates a new product and an associated price.
// If the price creation fails, it attempts to remove the newly created product.
// The function returns a composite error if the rollback also fails.
func (s *ProductAggregatorService) CreateProductWithPrice(ctx context.Context, request *domain.CreateProductWithPriceRequest) (string, string, error) {
	productID, err := s.productService.CreateProduct(ctx, &request.Product)
	if err != nil {
		// domain error are alredy handled
		return "", "", err
	}
	priceID, err := s.priceService.CreatePrice(ctx, productID, &request.Price)
	if err != nil {
		// rollback if error
		// if err := s.productService.RemoveProduct(ctx, productID); err != nil {
		if rollbackErr := s.productService.RemoveProduct(ctx, productID); !errors.Is(rollbackErr, errs.ErrDomainNotFound) {
			return "", "", fmt.Errorf("failed to create price (%w), AND product rollback also failed: %w", err, rollbackErr)
		}
		return "", "", err
	}
	return productID, priceID, nil
}
