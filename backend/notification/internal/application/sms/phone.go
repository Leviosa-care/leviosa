package sms

func (s *SMSService) SetPhone(phone string) {
	s.cache.SetCompanyPhone(phone)
}

func (s *SMSService) GetPhone() string {
	return s.cache.GetCompanyPhone()
}
