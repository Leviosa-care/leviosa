package promotionCodeHandler

const (
	// Base paths
	PromotionCodesBasePath      = "/promotion-codes"
	AdminPromotionCodesBasePath = "/admin/promotion-codes"

	// === Promotion Code Resource Paths ===
	IDPath         = "/{id}"
	CodePath       = "/code"
	CodeValuePath  = "/{code}"
	ValidatePath   = "/validate"
	ActivePath     = "/active"
	DeactivatePath = "/deactivate"

	// === Public Promotion Code Endpoints ===

	// Validate promotion code (public access for checkout)
	ValidatePromotionCodeEndpoint = PromotionCodesBasePath + ValidatePath

	// Get promotion code with associated coupon by code (public access)
	GetPromotionCodeWithCouponEndpoint = PromotionCodesBasePath + CodePath + CodeValuePath

	// === Admin-Only Endpoints ===

	// Get all promotion codes (admin only)
	GetAllPromotionCodesEndpoint = AdminPromotionCodesBasePath

	// Get active promotion codes (admin only)
	GetActivePromotionCodesEndpoint = AdminPromotionCodesBasePath + ActivePath

	// Get promotion code by ID (admin only)
	GetPromotionCodeByIDEndpoint = AdminPromotionCodesBasePath + IDPath

	// Get promotion code by code string (admin only)
	GetPromotionCodeByCodeEndpoint = AdminPromotionCodesBasePath + CodePath + CodeValuePath

	// Create promotion code (admin only)
	CreatePromotionCodeEndpoint = AdminPromotionCodesBasePath

	// Update promotion code (admin only)
	UpdatePromotionCodeEndpoint = AdminPromotionCodesBasePath + IDPath

	// Deactivate promotion code (admin only)
	DeactivatePromotionCodeEndpoint = AdminPromotionCodesBasePath + IDPath + DeactivatePath

	// Delete promotion code (admin only)
	DeletePromotionCodeEndpoint = AdminPromotionCodesBasePath + IDPath
)
