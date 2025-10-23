package domain

// CreateProductWithPriceRequest combines the product and price creation inputs.
type CreateProductWithPriceRequest struct {
	Product CreateProductRequest `json:"product"`
	Price   CreatePriceRequest   `json:"price"`
}

type ProductWithPrices struct {
	Product *ProductRes `json:"product"`
	Prices  []*Price    `json:"prices,omitempty"`
}

type ProductAggregator struct {
	Product *ProductRes `json:"product"`
	Prices  []*Price    `json:"prices,omitempty"`
	Image   *Image      `json:"images,omitempty"`
}
