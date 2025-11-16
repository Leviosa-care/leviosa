package aggregator

import (
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
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