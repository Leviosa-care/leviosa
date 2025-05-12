package mailService

import (
	"context"
	"time"

	// TODO: import the product thing because that is the thing that is the main part
	"github.com/hengadev/leviosa/internal/domain/user/models"
)

func (s *service) WelcomeUser(ctx context.Context, email string, user *models.User, legalAddress, companyInstagram string) error {
	// where to get that thing
	type WelcomeMailService struct {
		Title       string
		Description string
		SVG         string
		Link        string
	}

	data := struct {
		PromoCode int
		Services  []WelcomeMailService
		Year      int
	}{
		PromoCode: 0,
		// Services:,
		Year: time.Now().Year(),
	}

	if err := s.sendMail(
		ctx,
		email,
		"Bienvenue chez Leviosa Care",
		"welcome",
		data,
		nil,
		nil,
	); err != nil {
		return err
	}
	// handle promo code ? use the value and check if the value > 0 using the 'gt' built in function
	return nil
}
