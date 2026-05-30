package partnerHandler

const (
	// Base paths
	PartnersBasePath      = "/partners"
	AdminPartnersBasePath = "/admin/partners"

	// === Partner Resource Paths ===

	// Partner category
	CategoryPath = "/categories"

	// Partner product
	ProductPath = "/products"

	// Partner profile
	MePath = "/me"

	// Partner verification
	VerifyPath = "/verify"

	// Catalog validation endpoints
	ValidatePath = "/validate"

	// === Public Partner Endpoints ===

	// Get all public partners (unauthenticated)
	GetPublicPartnersEndpoint = PartnersBasePath

	// Get partner by ID
	GetPartnerByIDEndpoint = PartnersBasePath + "/{id}"

	// Get partner by category
	GetPartnersByCategoryEndpoint = PartnersBasePath + CategoryPath + "/{id}"

	// Get partner by categories
	GetPartnersByCategoriesEndpoint = PartnersBasePath + CategoryPath

	// Get partner by product
	GetPartnersByProductEndpoint = PartnersBasePath + ProductPath + "/{id}"

	// Get partner by products
	GetPartnersByProductsEndpoint = PartnersBasePath + ProductPath

	// === Authenticated Partner Endpoints ===

	// Get authenticated partner's own profile
	GetPartnerMeEndpoint = PartnersBasePath + MePath

	// Delete authenticated partner's own profile
	DeletePartnerMeEndpoint = PartnersBasePath + MePath

	// Update authenticated partner's own profile
	UpdatePartnerMeEndpoint = PartnersBasePath + MePath

	// === Admin-Only Endpoints ===

	// // Create partner
	// CreatePartnerEndpoint = PartnersBasePath

	// Get all partners
	GetAllPartnersEndpoint = AdminPartnersBasePath

	// Delete partner
	DeletePartnerEndpoint = AdminPartnersBasePath + "/{id}"

	// Update partner profile (partner can update their own, admin can update any)
	UpdatePartnerEndpoint = PartnersBasePath + "/{id}"

	// Verify partner credentials
	VerifyPartnerEndpoint = AdminPartnersBasePath + "/{id}" + VerifyPath

	// === Catalog Validation Endpoints ===

	// Validate products exist in catalog
	ValidatePartnerProductsEndpoint = AdminPartnersBasePath + "/products" + ValidatePath

	// === Stripe Endpoints ===

	// Get Stripe onboarding link for the authenticated partner
	GetOnboardingLinkEndpoint = PartnersBasePath + MePath + "/stripe/onboarding-link"
)
