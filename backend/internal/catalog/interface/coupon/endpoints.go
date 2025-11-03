package couponHandler

const (
	// Base paths
	CouponsBasePath      = "/coupons"
	AdminCouponsBasePath = "/admin/coupons"

	// === Coupon Resource Paths ===
	IDPath           = "/{id}"
	ValidatePath     = "/validate"
	ValidPath        = "/valid"
	StripePath       = "/stripe"
	StripeIDPath     = "/{stripeId}"
	DeactivatePath   = "/deactivate"

	// === Public Coupon Endpoints ===

	// Validate coupon code (public access for checkout)
	ValidateCouponEndpoint = CouponsBasePath + ValidatePath

	// Get all valid/active coupons (public access)
	GetValidCouponsEndpoint = CouponsBasePath + ValidPath

	// === Admin-Only Endpoints ===

	// Get all coupons (admin only)
	GetAllCouponsEndpoint = AdminCouponsBasePath

	// Get coupon by ID (admin only)
	GetCouponByIDEndpoint = AdminCouponsBasePath + IDPath

	// Get coupon by Stripe ID (admin only)
	GetCouponByStripeIDEndpoint = AdminCouponsBasePath + StripePath + StripeIDPath

	// Create coupon (admin only)
	CreateCouponEndpoint = AdminCouponsBasePath

	// Update coupon (admin only)
	UpdateCouponEndpoint = AdminCouponsBasePath + IDPath

	// Deactivate coupon (admin only)
	DeactivateCouponEndpoint = AdminCouponsBasePath + IDPath + DeactivatePath

	// Delete coupon (admin only)
	DeleteCouponEndpoint = AdminCouponsBasePath + IDPath
)
