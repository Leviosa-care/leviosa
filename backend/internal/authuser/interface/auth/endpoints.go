package aggregatorHandler

const (
	// Base paths
	AuthBasePath      = "/auth"
	AdminAuthBasePath = "/admin/auth"

	// === Authentication Resource Paths ===

	// Email verification and OTP
	EmailPath = "/email"
	OTPPath   = "/otp"

	// User completion
	CompletePath = "/complete"

	// Login/Logout
	LoginPath  = "/login"
	LogoutPath = "/logout"

	// User deletion
	MePath         = "/me"
	AdminUsersPath = "/users"

	// Password reset paths
	PasswordPath        = "/password"
	ResetPath           = "/reset"
	RequestPath         = "/request"
	ValidatePath        = "/validate"
	ConfirmPath         = "/confirm"
	PasswordResetPath   = PasswordPath + ResetPath
	RequestResetPath    = PasswordResetPath + RequestPath
	ValidateResetPath   = PasswordResetPath + ValidatePath
	ConfirmResetPath    = PasswordResetPath + ConfirmPath

	// OAuth paths
	OAuthPath     = "/oauth"
	CallbackPath  = "/callback"
	ProviderParam = "/{provider}"

	// === Public Authentication Endpoints ===

	// Email and OTP verification
	CheckEmailSendOTPEndpoint      = AuthBasePath + EmailPath
	ValidateOTPCreatePendingEndpoint = AuthBasePath + OTPPath

	// User login/logout
	SignInEndpoint  = AuthBasePath + LoginPath
	SignOutEndpoint = AuthBasePath + LogoutPath

	// Password reset flow
	RequestPasswordResetEndpoint  = AuthBasePath + RequestResetPath
	ValidatePasswordResetOTPEndpoint = AuthBasePath + ValidateResetPath
	ConfirmPasswordResetEndpoint  = AuthBasePath + ConfirmResetPath

	// OAuth endpoints
	OAuthStartEndpoint    = AuthBasePath + OAuthPath + ProviderParam
	OAuthCallbackEndpoint = AuthBasePath + OAuthPath + ProviderParam + CallbackPath

	// === Authenticated User Endpoints ===

	// Requires Visitor role or higher
	CompleteUserEndpoint    = AuthBasePath + CompletePath
	CompletePartnerEndpoint = AuthBasePath + CompletePath + "/partner"

	// Requires Standard role or higher
	DeleteOwnAccountEndpoint = AuthBasePath + MePath

	// === Admin-Only Endpoints ===

	// User management (admin only)
	DeleteUserByAdminEndpoint = AdminAuthBasePath + AdminUsersPath + "/{id}"
)
