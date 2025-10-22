package otp

func (s *OTPService) effectiveOTPLength() int {
	length := defaultOTPLength
	if length <= 0 {
		return 6 // Fallback to hardcoded value (should never happen with defaultOTPLength = 6)
	}
	return length
}
