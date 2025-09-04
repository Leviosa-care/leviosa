package ports

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
)

type UserService interface {
	CreatePendingUser(ctx context.Context, email string) (uuid.UUID, error)
}
