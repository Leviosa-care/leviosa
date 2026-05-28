package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID `json:"id"`
	State            UserState `json:"state"`
	Email            string    `json:"email" encx:"encrypt,hash_basic"`
	Password         string    `json:"-" encx:"hash_secure"`
	Picture          string    `json:"picture" encx:"encrypt"`
	CreatedAt        time.Time `json:"created_at" encx:"encrypt"`
	LoggedInAt       time.Time `json:"logged_in_at" encx:"encrypt"`
	Role             string    `json:"-" encx:"encrypt"`
	BirthDate        time.Time `json:"birthdate" encx:"encrypt"`
	LastName         string    `json:"last_name" encx:"encrypt"`
	FirstName        string    `json:"first_name" encx:"encrypt"`
	Gender           string    `json:"gender" encx:"encrypt"`
	Telephone        string    `json:"telephone" encx:"encrypt,hash_basic"`
	PostalCode       string    `json:"postal_code" encx:"encrypt"`
	City             string    `json:"city" encx:"encrypt"`
	Address1         string    `json:"address1" encx:"encrypt"`
	Address2         string    `json:"address2" encx:"encrypt"`
	GoogleID         string    `json:"google_id" encx:"encrypt"`
	AppleID          string    `json:"-" encx:"encrypt"`
	StripeCustomerID string    `json:"-" encx:"encrypt"`
	ProfileIncomplete bool     `json:"profile_incomplete"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:               u.ID,
		State:            u.State,
		Email:            u.Email,
		Picture:          u.Picture,
		CreatedAt:        u.CreatedAt,
		LoggedInAt:       u.LoggedInAt,
		Role:             u.Role,
		BirthDate:        u.BirthDate,
		LastName:         u.LastName,
		FirstName:        u.FirstName,
		Gender:           u.Gender,
		Telephone:        u.Telephone,
		PostalCode:       u.PostalCode,
		City:             u.City,
		Address1:         u.Address1,
		Address2:         u.Address2,
		GoogleID:         u.GoogleID,
		AppleID:          u.AppleID,
		HasPassword:      u.Password != "",
		ProfileIncomplete: u.ProfileIncomplete,
	}
}
