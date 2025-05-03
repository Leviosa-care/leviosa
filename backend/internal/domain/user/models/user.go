package models

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/hengadev/errsx"
)

type User struct {
	ID                  string    `json:"-"`
	Email               string    `json:"-"`
	EmailHash           string    `json:"-"`
	Password            string    `json:"-"`
	PasswordHash        string    `json:"-"`
	Picture             string    `json:"-"`
	CreatedAt           time.Time `json:"-"`
	EncryptedCreatedAt  string    `json:"-"`
	LoggedInAt          time.Time `json:"-"`
	EncryptedLoggedInAt string    `json:"-"`
	Role                string    `json:"-"`
	BirthDate           time.Time `json:"-"`
	EncryptedBirthDate  string    `json:"-"`
	LastName            string    `json:"-"`
	FirstName           string    `json:"-"`
	Gender              string    `json:"-"`
	Telephone           string    `json:"-"`
	PostalCode          string    `json:"-"`
	City                string    `json:"-"`
	Address1            string    `json:"-"`
	Address2            string    `json:"-"`
	GoogleID            string    `json:"-"`
	AppleID             string    `json:"-"`
}

func (a *User) Create() {
	a.CreatedAt = time.Now().UTC()
}

func (a *User) Login() {
	a.LoggedInAt = time.Now().UTC()
}

// I need to use the hash things in that function
func NewUser(
	user UserSignUp,
	role Role,
) *User {
	return &User{
		ID:         uuid.NewString(),
		Email:      user.Email,
		Password:   user.Password,
		Role:       BASIC.String(),
		BirthDate:  user.BirthDate,
		LastName:   user.LastName,
		FirstName:  user.FirstName,
		Gender:     user.Gender,
		Telephone:  user.Telephone,
		PostalCode: user.PostalCode,
		City:       user.City,
		Address1:   user.Address1,
		Address2:   user.Address2,
	}
}

func (u User) Valid(ctx context.Context) error {
	var errs errsx.Map
	return errs.AsError()
}

// Interface implementation

func (u User) AssertComparable() {}

func (u User) GetSQLColumnMapping() map[string]string {
	return map[string]string{
		"ID":                  "id",
		"Email":               "encrypted_email",
		"EmailHash":           "email_hash",
		"PasswordHash":        "password_hash",
		"Picture":             "encrypted_picture",
		"EncryptedCreatedAt":  "encrypted_created_at",
		"EncryptedLoggedInAt": "encrypted_logged_in_at",
		"Role":                "role",
		"EncryptedBirthDate":  "encrypted_birthdate",
		"LastName":            "encrypted_lastname",
		"FirstName":           "encrypted_firstname",
		"Gender":              "encrypted_gender",
		"Telephone":           "encrypted_telephone",
		"PostalCode":          "encrypted_postal_code",
		"City":                "encrypted_city",
		"Address1":            "encrypted_address1",
		"Address2":            "encrypted_address2",
		"GoogleID":            "encrypted_google_id",
		"AppleID":             "encrypted_apple_id",
	}
}

func (u User) GetProhibitedFields() []string {
	return []string{
		"ID",
		"Email",
		"Password",
		"GoogleID",
		"AppleID",
	}
}
