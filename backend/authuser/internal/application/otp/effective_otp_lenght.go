package otp

const DefaultOTPLength = 6

func (s *OTPService) effectiveOTPLength() int {
	length := s.GetOTPLength()
	if length <= 0 {
		return DefaultOTPLength
	}
	return length
}
