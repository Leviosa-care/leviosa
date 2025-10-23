package userRepository_test

import (
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
)

func NewTestUserEncx() *domain.UserEncx {
	return &domain.UserEncx{
		ID:    uuid.New(),
		State: domain.Unverified,
	}
}
