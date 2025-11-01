package http

const (
	// Base paths
	SettingsBasePath        = "/settings"
	AdminSettingsBasePath   = "/admin/settings"
	InternalSettingsBasePath = "/internal/settings"

	// === Company Settings Resource Paths ===

	// Company resource identifiers (relative paths)
	CompanyNamePath     = "/name"
	CompanyEmailPath    = "/email"
	CompanyPhonePath    = "/phone"
	CompanyAddressPath  = "/address"
	CompanyInstagramPath = "/instagram"
	CompanyLogoPath     = "/logo"

	// Public company endpoints (GET only, no authentication)
	GetCompanyNameEndpoint     = SettingsBasePath + CompanyNamePath
	GetCompanyEmailEndpoint    = SettingsBasePath + CompanyEmailPath
	GetCompanyAddressEndpoint  = SettingsBasePath + CompanyAddressPath
	GetCompanyInstagramEndpoint = SettingsBasePath + CompanyInstagramPath
	GetCompanyLogoEndpoint     = SettingsBasePath + CompanyLogoPath

	// Admin company endpoints (Administrator role required)
	AdminGetCompanyPhoneEndpoint = AdminSettingsBasePath + CompanyPhonePath
	SetCompanyNameEndpoint       = AdminSettingsBasePath + CompanyNamePath
	SetCompanyEmailEndpoint      = AdminSettingsBasePath + CompanyEmailPath
	SetCompanyPhoneEndpoint      = AdminSettingsBasePath + CompanyPhonePath
	SetCompanyAddressEndpoint    = AdminSettingsBasePath + CompanyAddressPath
	SetCompanyInstagramEndpoint  = AdminSettingsBasePath + CompanyInstagramPath
	SetCompanyLogoEndpoint       = AdminSettingsBasePath + CompanyLogoPath

	// Internal company endpoints (Service-to-service authentication)
	InternalGetCompanyNameEndpoint     = InternalSettingsBasePath + CompanyNamePath
	InternalGetCompanyEmailEndpoint    = InternalSettingsBasePath + CompanyEmailPath
	InternalGetCompanyPhoneEndpoint    = InternalSettingsBasePath + CompanyPhonePath
	InternalGetCompanyAddressEndpoint  = InternalSettingsBasePath + CompanyAddressPath
	InternalGetCompanyInstagramEndpoint = InternalSettingsBasePath + CompanyInstagramPath
	InternalGetCompanyLogoEndpoint     = InternalSettingsBasePath + CompanyLogoPath

	// === OTP Settings Resource Paths ===

	// OTP resource identifiers
	OTPBasePath        = "/otp"
	OTPDurationPath    = "/duration"
	OTPLengthPath      = "/length"
	OTPMaxAttemptsPath = "/max-attempts"

	// Admin OTP endpoints (Administrator role required)
	AdminGetOTPDurationEndpoint = AdminSettingsBasePath + OTPBasePath + OTPDurationPath
	AdminSetOTPDurationEndpoint = AdminSettingsBasePath + OTPBasePath + OTPDurationPath
	AdminGetOTPLengthEndpoint   = AdminSettingsBasePath + OTPBasePath + OTPLengthPath
	AdminSetOTPLengthEndpoint   = AdminSettingsBasePath + OTPBasePath + OTPLengthPath
	AdminGetOTPMaxAttemptsEndpoint = AdminSettingsBasePath + OTPBasePath + OTPMaxAttemptsPath
	AdminSetOTPMaxAttemptsEndpoint = AdminSettingsBasePath + OTPBasePath + OTPMaxAttemptsPath

	// Internal OTP endpoints (Service-to-service authentication)
	InternalGetOTPDurationEndpoint = InternalSettingsBasePath + OTPBasePath + OTPDurationPath
	InternalGetOTPLengthEndpoint   = InternalSettingsBasePath + OTPBasePath + OTPLengthPath
	InternalGetOTPMaxAttemptsEndpoint = InternalSettingsBasePath + OTPBasePath + OTPMaxAttemptsPath

	// === Token Settings Resource Paths ===

	// Token resource identifiers
	TokensBasePath         = "/tokens"
	AccessTokenPath        = "/access-duration"
	RefreshTokenPath       = "/refresh-duration"

	// Admin token endpoints (Administrator role required)
	AdminGetAccessTokenDurationEndpoint  = AdminSettingsBasePath + TokensBasePath + AccessTokenPath
	AdminSetAccessTokenDurationEndpoint  = AdminSettingsBasePath + TokensBasePath + AccessTokenPath
	AdminGetRefreshTokenDurationEndpoint = AdminSettingsBasePath + TokensBasePath + RefreshTokenPath
	AdminSetRefreshTokenDurationEndpoint = AdminSettingsBasePath + TokensBasePath + RefreshTokenPath

	// Internal token endpoints (Service-to-service authentication)
	InternalGetAccessTokenDurationEndpoint  = InternalSettingsBasePath + TokensBasePath + AccessTokenPath
	InternalGetRefreshTokenDurationEndpoint = InternalSettingsBasePath + TokensBasePath + RefreshTokenPath

	// === Bulk Settings Endpoints ===

	BulkPath = "/bulk"

	// Admin bulk endpoint (Administrator role required)
	AdminBulkEndpoint = AdminSettingsBasePath + BulkPath

	// Internal bulk endpoint (Service-to-service authentication)
	InternalBulkEndpoint = InternalSettingsBasePath + BulkPath
)