package mailService

import (
	"fmt"
	"os"

	"github.com/hengadev/leviosa/internal/domain"
)

type Service struct {
	from     string
	email    string
	password string
}

func New() (*Service, error) {
	email := os.Getenv("GMAIL_EMAIL")
	if email == "" {
		return nil, domain.NewNotFoundErr(fmt.Errorf("environment variable 'GMAIL_EMAIL'"))
	}
	password := os.Getenv("GMAIL_PASSWORD")
	if password == "" {
		return nil, domain.NewNotFoundErr(fmt.Errorf("environment variable 'GMAIL_PASSWORD'"))
	}
	return &Service{
		from:     "gary.testmail.123@gmail.com",
		email:    email,
		password: password,
	}, nil
}
