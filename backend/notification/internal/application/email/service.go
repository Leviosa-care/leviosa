package email

import (
	"github.com/Leviosa-care/notification/internal/domain"
	"github.com/Leviosa-care/notification/internal/ports"
)

type EmailService struct {
	emailClient ports.EmailService
	cache       *domain.CompanyCache
}

func New(emailClient ports.EmailService, cache *domain.CompanyCache) *EmailService {
	return &EmailService{
		emailClient: emailClient,
		cache:       cache,
	}
}
