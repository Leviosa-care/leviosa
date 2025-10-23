package otp

func (s *OTPService) GetOTPDuration() int {
	return s.cache.GetOTPDuration()
}

func (s *OTPService) SetOTPDuration(duration int) {
	s.cache.SetOTPDuration(duration)
}

func (s *OTPService) GetOTPLength() int {
	return s.cache.GetOTPLength()
}

func (s *OTPService) SetOTPLength(length int) {
	s.cache.SetOTPLength(length)
}

func (s *OTPService) GetOTPMaxAttempts() int {
	return s.cache.GetOTPMaxAttempts()
}

func (s *OTPService) SetOTPMaxAttempts(maxAttempts int) {
	s.cache.SetOTPMaxAttempts(maxAttempts)
}
