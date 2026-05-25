package priceHandler

const (
	// Base paths
	ProductsBasePath       = "/products"
	AdminPricesBasePath   = "/admin/prices"
	AdminProductsBasePath = "/admin/products"

	// === Price Resource Paths ===
	IDPath     = "/{id}"
	PricesPath = "/prices"

	// === Public Endpoints ===

	// Get all active prices for a product (public access)
	GetPublicProductPricesEndpoint = ProductsBasePath + IDPath + PricesPath

	// === Admin-Only Endpoints ===
	// (All price management endpoints are admin-only)

	// Get price by ID (admin only)
	GetPriceEndpoint = AdminPricesBasePath + IDPath

	// Get all prices for a product (admin only)
	GetPricesByProductIDEndpoint = AdminProductsBasePath + IDPath + PricesPath

	// Create price for a product (admin only)
	CreatePriceEndpoint = AdminProductsBasePath + IDPath + PricesPath

	// Update price (admin only)
	UpdatePriceEndpoint = AdminPricesBasePath + IDPath
)
