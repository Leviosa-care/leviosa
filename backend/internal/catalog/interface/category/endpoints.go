package categoryHandler

const (
	// Base paths
	CategoriesBasePath      = "/categories"
	AdminCategoriesBasePath = "/admin/categories"

	// === Category Resource Paths ===
	IDPath = "/{id}"

	// === Public Category Endpoints ===

	// Get all published categories (public access)
	GetAllPublishedCategoriesEndpoint = CategoriesBasePath

	// Get category by ID (public access)
	GetCategoryByIDEndpoint = CategoriesBasePath + IDPath

	// === Admin-Only Endpoints ===

	// Get all categories including drafts (admin only)
	GetAdminAllCategoriesEndpoint = AdminCategoriesBasePath

	// Create category (admin only)
	CreateCategoryEndpoint = AdminCategoriesBasePath

	// Modify category (admin only)
	ModifyCategoryEndpoint = AdminCategoriesBasePath + IDPath

	// Remove category (admin only)
	RemoveCategoryEndpoint = AdminCategoriesBasePath + IDPath
)
