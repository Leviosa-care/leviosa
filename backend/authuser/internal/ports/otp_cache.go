package ports

// OTPCache defines the interface for OTP configuration caching
type OTPCache interface {
	// GetOTPDuration retrieves the cached OTP duration setting in minutes
	GetOTPDuration() int

	// SetOTPDuration updates the cached OTP duration setting
	SetOTPDuration(duration int)

	// GetOTPLength retrieves the cached OTP length setting
	GetOTPLength() int

	// SetOTPLength updates the cached OTP length setting
	SetOTPLength(length int)

	// GetOTPMaxAttempts retrieves the cached OTP max attempts setting
	GetOTPMaxAttempts() int

	// SetOTPMaxAttempts updates the cached OTP max attempts setting
	SetOTPMaxAttempts(maxAttempts int)
}
