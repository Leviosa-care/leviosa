package priceHandler

const (
	// Base paths
	AdminPricesBasePath   = "/admin/prices"
	AdminProductsBasePath = "/admin/products"

	// === Price Resource Paths ===
	IDPath     = "/{id}"
	PricesPath = "/prices"

	// === Admin-Only Endpoints ===
	// (All price endpoints are admin-only)

	// Get price by ID (admin only)
	GetPriceEndpoint = AdminPricesBasePath + IDPath

	// Get all prices for a product (admin only)
	GetPricesByProductIDEndpoint = AdminProductsBasePath + IDPath + PricesPath

	// Create price for a product (admin only)
	CreatePriceEndpoint = AdminProductsBasePath + IDPath + PricesPath

	// Update price (admin only)
	UpdatePriceEndpoint = AdminPricesBasePath + IDPath
)
