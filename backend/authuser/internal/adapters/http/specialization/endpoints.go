package specializationHandler

const (
	// Base paths
	SpecializationsBasePath      = "/specializations"
	AdminSpecializationsBasePath = "/admin/specializations"

	// === Specialization Resource Paths ===
	// (All endpoints use base paths with optional ID parameter)

	// === Public Specialization Endpoints ===

	// Get all active specializations (any authenticated user can view for selection)
	GetAllSpecializationsEndpoint = SpecializationsBasePath

	// === Admin-Only Endpoints ===

	// Create new specialization
	CreateSpecializationEndpoint = AdminSpecializationsBasePath

	// Get specialization by ID
	GetSpecializationByIDEndpoint = AdminSpecializationsBasePath + "/{id}"

	// Update specialization
	UpdateSpecializationEndpoint = AdminSpecializationsBasePath + "/{id}"

	// Delete specialization
	DeleteSpecializationEndpoint = AdminSpecializationsBasePath + "/{id}"
)
