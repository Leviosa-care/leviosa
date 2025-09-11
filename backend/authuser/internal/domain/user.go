package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                        uuid.UUID `json:"id"`
	State                     UserState `json:"state"`
	Email                     string    `json:"email" encx:"encrypt,hash_basic"`
	EmailHash                 string    `json:"-"`
	EmailEncrypted            []byte    `json:"-"`
	Password                  string    `json:"-" encx:"hash_secure"`
	PasswordHash              string    `json:"-"`
	Picture                   string    `json:"picture" encx:"encrypt"`
	PictureEncrypted          []byte    `json:"-"`
	CreatedAt                 time.Time `json:"created_at" encx:"encrypt"`
	CreatedAtEncrypted        []byte    `json:"-"`
	LoggedInAt                time.Time `json:"logged_in_at" encx:"encrypt"`
	LoggedInAtEncrypted       []byte    `json:"-"`
	Role                      string    `json:"-" encx:"encrypt"`
	RoleEncrypted             []byte    `json:"-"`
	BirthDate                 time.Time `json:"birthdate" encx:"encrypt"`
	BirthDateEncrypted        []byte    `json:"-"`
	LastName                  string    `json:"last_name" encx:"encrypt"`
	LastNameEncrypted         []byte    `json:"-"`
	FirstName                 string    `json:"first_name" encx:"encrypt"`
	FirstNameEncrypted        []byte    `json:"-"`
	Gender                    string    `json:"gender" encx:"encrypt"`
	GenderEncrypted           []byte    `json:"-"`
	Telephone                 string    `json:"telephone" encx:"encrypt,hash_basic"`
	TelephoneHash             string    `json:"-"`
	TelephoneEncrypted        []byte    `json:"-"`
	PostalCode                string    `json:"postal_code" encx:"encrypt"`
	PostalCodeEncrypted       []byte    `json:"-"`
	City                      string    `json:"city" encx:"encrypt"`
	CityEncrypted             []byte    `json:"-"`
	Address1                  string    `json:"address1" encx:"encrypt"`
	Address1Encrypted         []byte    `json:"-"`
	Address2                  string    `json:"address2" encx:"encrypt"`
	Address2Encrypted         []byte    `json:"-"`
	GoogleID                  string    `json:"google_id" encx:"encrypt"`
	GoogleIDEncrypted         []byte    `json:"-"`
	AppleID                   string    `json:"-" encx:"encrypt"`
	AppleIDEncrypted          []byte    `json:"-"`
	StripeCustomerID          string    `json:"-" encx:"encrypt"`
	StripeCustomerIDEncrypted []byte    `json:"-"`
	DEK                       []byte    `json:"-" encx:"encrypt"`
	DEKEncrypted              []byte    `json:"-"`
	KeyVersion                int       `json:"-"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		State:      u.State,
		Email:      u.Email,
		Picture:    u.Picture,
		CreatedAt:  u.CreatedAt,
		LoggedInAt: u.LoggedInAt,
		Role:       u.Role,
		BirthDate:  u.BirthDate,
		LastName:   u.LastName,
		FirstName:  u.FirstName,
		Gender:     u.Gender,
		Telephone:  u.Telephone,
		PostalCode: u.PostalCode,
		City:       u.City,
		Address1:   u.Address1,
		Address2:   u.Address2,
		GoogleID:   u.GoogleID,
		AppleID:    u.AppleID,
	}
}
