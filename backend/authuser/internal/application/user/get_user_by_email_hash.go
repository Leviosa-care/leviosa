package user

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
)

func (s *UserService) GetUserByEmailHash(ctx context.Context, emailHash string) (*domain.UserResponse, error) {
	return nil, nil
}
