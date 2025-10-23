package partnerHandler

const (
	// Base paths
	PartnersBasePath      = "/partners"
	AdminPartnersBasePath = "/admin/partners"

	// === Partner Resource Paths ===

	// Partner identification
	UserPath = "/user"

	// Partner verification
	VerifyPath = "/verify"

	// Partner specializations
	SpecializationsPath = "/specializations"

	// Catalog validation endpoints
	ValidatePath = "/validate"

	// === Public Partner Endpoints ===
	// (None - all partner management requires authentication)

	// === Authenticated Partner Endpoints ===

	// Get partner specializations (any authenticated user can view)
	GetPartnerSpecializationsEndpoint = PartnersBasePath + "/{id}" + SpecializationsPath

	// Update partner profile (partner can update their own, admin can update any)
	UpdatePartnerEndpoint = PartnersBasePath + "/{id}"

	// === Admin-Only Endpoints ===

	// Create partner
	CreatePartnerEndpoint = PartnersBasePath

	// Get partner by ID
	GetPartnerByIDEndpoint = AdminPartnersBasePath + "/{id}"

	// Get partner by user ID
	GetPartnerByUserIDEndpoint = AdminPartnersBasePath + UserPath + "/{userID}"

	// Get all partners
	GetAllPartnersEndpoint = AdminPartnersBasePath

	// Delete partner
	DeletePartnerEndpoint = AdminPartnersBasePath + "/{id}"

	// Verify partner credentials
	VerifyPartnerEndpoint = AdminPartnersBasePath + "/{id}" + VerifyPath

	// Add specialization to partner
	AddPartnerSpecializationEndpoint = AdminPartnersBasePath + "/{id}" + SpecializationsPath + "/{specializationID}"

	// Remove specialization from partner
	RemovePartnerSpecializationEndpoint = AdminPartnersBasePath + "/{id}" + SpecializationsPath + "/{specializationID}"

	// === Catalog Validation Endpoints ===

	// Validate specializations exist in catalog
	ValidatePartnerSpecializationsEndpoint = AdminPartnersBasePath + SpecializationsPath + ValidatePath

	// Validate products exist in catalog
	ValidatePartnerProductsEndpoint = AdminPartnersBasePath + "/products" + ValidatePath
)
