package userService

import (
	"context"
	// "fmt"

	"github.com/hengadev/leviosa/internal/domain/user/models"
)

// TODO: I am not sure that I no longer need that function

func (s *service) CreateOAuthAccount(
	ctx context.Context,
	userCandidate *models.OAuthUser,
) (*models.User, error) {
	// NOTE: that thing does not work brother
	// user, err := userCandidate.ToUser()
	// if err != nil {
	// 	return nil, fmt.Errorf("convert OAuthUser to User")
	// }
	// // NOTE: We use the same function as the email-password one but we need to make the function better to do account linking.
	// err = s.AddUser(ctx, user)
	// if err != nil {
	// 	return nil, fmt.Errorf("add oauth account %w", err)
	// }
	// user.Role = models.BASIC.String()
	// user.Create()
	// user.Login()
	// fmt.Printf("the user is after the conversion from the OAuth: %#+v\n", user)
	// return user, nil

	return nil, nil
}
