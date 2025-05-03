package models

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/hengadev/errsx"
)

type User struct {
	ID                  string    `json:"-"`
	Email               string    `json:"-" encx:"encrypt,hash_basic"`
	EmailHash           string    `json:"-"`
	EmailEncrypted      []byte    `json:"-"`
	Password            string    `json:"-" encx:"hash_secure"`
	PasswordHash        string    `json:"-"`
	Picture             string    `json:"-" encx:"encrypt"`
	PictureEncrypted    []byte    `json:"-"`
	CreatedAt           time.Time `json:"-"`
	LoggedInAt          time.Time `json:"-"`
	Role                string    `json:"-"`
	BirthDate           time.Time `json:"-" encx:"encrypt"`
	BirthDateEncrypted  []byte    `json:"-"`
	LastName            string    `json:"-" encx:"encrypt"`
	LastNameEncrypted   []byte    `json:"-"`
	FirstName           string    `json:"-" encx:"encrypt"`
	FirstNameEncrypted  []byte    `json:"-"`
	Gender              string    `json:"-" encx:"encrypt"`
	GenderEncrypted     []byte    `json:"-"`
	Telephone           string    `json:"-" encx:"encrypt,hash_basic"`
	TelephoneHash       string    `json:"-"`
	TelephoneEncrypted  []byte    `json:"-"`
	PostalCode          string    `json:"-" encx:"encrypt"`
	PostalCodeEncrypted []byte    `json:"-"`
	City                string    `json:"-" encx:"encrypt"`
	CityEncrypted       []byte    `json:"-"`
	Address1            string    `json:"-" encx:"encrypt"`
	Address1Encrypted   []byte    `json:"-"`
	Address2            string    `json:"-" encx:"encrypt"`
	Address2Encrypted   []byte    `json:"-"`
	GoogleID            string    `json:"-" encx:"encrypt"`
	GoogleIDEncrypted   []byte    `json:"-"`
	AppleID             string    `json:"-" encx:"encrypt"`
	AppleIDEncrypted    []byte    `json:"-"`
	DEK                 []byte    `json:"-" encx:"encrypt"`
	DEKEncrypted        []byte    `json:"-"`
	KeyVersion          int       `json:"-"`
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
		"EmailEncrypted":      "email_encrypted",
		"EmailHash":           "email_hash",
		"PasswordHash":        "password_hash",
		"PictureEncrypted":    "picture_encrypted",
		"CreatedAtEncrypted":  "created_at_encrypted",
		"LoggedInAtEncrypted": "logged_in_at_encrypted",
		"Role":                "role",
		"BirthDateEncrypted":  "birthdate_encrypted",
		"LastNameEncrypted":   "lastname_encrypted",
		"FirstNameEncrypted":  "firstname_encrypted",
		"GenderEncrypted":     "gender_encrypted",
		"TelephoneEncrypted":  "telephone_encrypted",
		"TelephoneHash":       "telephone_hash",
		"PostalCodeEncrypted": "postal_code_encrypted",
		"CityEncrypted":       "city_encrypted",
		"Address1Encrypted":   "address1_encrypted",
		"Address2Encrypted":   "address2_encrypted",
		"GoogleIDEncrypted":   "google_id_encrypted",
		"AppleIDEncrypted":    "apple_id_encrypted",
	}
}

func (u User) GetProhibitedFields() []string {
	return []string{
		"ID",
		"EmailEncrypted",
		"EmailHash",
		"PasswordHash",
		"GoogleIDEncrypted",
		"AppleIDEncrypted",
	}
}
