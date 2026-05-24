package authuseradapter

import (
	"context"
	"strings"

	authuserPorts "github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	"github.com/google/uuid"
)

type nameFetcher struct {
	svc authuserPorts.UserService
}

// New returns a UserNameFetcher backed by the authuser service.
func New(svc authuserPorts.UserService) *nameFetcher {
	return &nameFetcher{svc: svc}
}

// FetchName resolves a user ID to "FirstName LastName", falling back to "Utilisateur".
func (f *nameFetcher) FetchName(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := f.svc.GetUserByID(ctx, userID)
	if err != nil {
		return "Utilisateur", nil
	}
	name := strings.TrimSpace(user.FirstName + " " + user.LastName)
	if name == "" {
		return "Utilisateur", nil
	}
	return name, nil
}
