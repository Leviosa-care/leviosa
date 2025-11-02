package oauth

import (
	"fmt"
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// ExchangeCodeForUser exchanges authorization code for user information using HTTP request
func (s *Service) ExchangeCodeForUser(w http.ResponseWriter, r *http.Request) (goth.User, error) {
	// Use gothic's built-in method to complete the authentication
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return goth.User{}, fmt.Errorf("failed to complete user auth: %w", err)
	}

	return user, nil
}
