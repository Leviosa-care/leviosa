package factories

import (
	"time"

	"github.com/hengadev/leviosa/internal/domain/user/models"
)

var birthdate, _ = time.Parse("2006-01-02", "1998-07-12")

func NewBasicUser(overrides map[string]interface{}) *models.User {
	user := &models.User{
		ID:                 "123e4567-e89b-12d3-a456-426614174000",
		Email:              "john.doe@example.com",
		EmailHash:          "john.doe@example.com",
		Password:           "password",
		PasswordHash:       "hashedpassword",
		Picture:            "picture",
		CreatedAt:          time.Now(),
		LoggedInAt:         time.Now(),
		Role:               "basic",
		BirthDateEncrypted: []byte("1998-07-12"),
		LastName:           "DOE",
		FirstName:          "John",
		Gender:             models.GenderInput{Gender: "male"},
		Telephone:          "0123456789",
		PostalCode:         "75000",
		City:               "Paris",
		Address1:           "01 Avenue Jean DUPONT",
		Address2:           "",
		GoogleID:           "google_id",
		AppleID:            "apple_id",
	}
	// Apply overrides
	for key, value := range overrides {
		switch key {
		case "ID":
			user.ID = value.(string)
		case "Email":
			user.Email = value.(string)
		case "EmailHash":
			user.EmailHash = value.(string)
		case "Password":
			user.Password = value.(string)
		case "PasswordHash":
			user.PasswordHash = value.(string)
		case "CreatedAt":
			user.CreatedAt = value.(time.Time)
		case "LoggedInAt":
			user.LoggedInAt = value.(time.Time)
		case "Role":
			user.Role = value.(string)
		case "BirthDate":
			user.BirthDate = value.(time.Time)
		case "EncryptedBirthDate":
			user.BirthDateEncrypted = value.([]byte)
		case "LastName":
			user.LastName = value.(string)
		case "FirstName":
			user.FirstName = value.(string)
		case "Gender":
			user.Gender = value.(models.GenderInput)
		case "Telephone":
			user.Telephone = value.(string)
		case "PostalCode":
			user.PostalCode = value.(string)
		case "City":
			user.City = value.(string)
		case "Address1":
			user.Address1 = value.(string)
		case "Address2":
			user.Address2 = value.(string)
		case "GoogleID":
			user.GoogleID = value.(string)
		case "AppleID":
			user.AppleID = value.(string)
		}
	}
	return user
}

func NewBasicUserList() []*models.User {
	users := []*models.User{
		NewBasicUser(nil),
		NewBasicUser(map[string]any{
			"ID":                 "123e4567-e89b-12d3-a456-426614174001",
			"Email":              "john.doe@example.com",
			"EmailHash":          "john.doe@example.com",
			"Picture":            "picture1",
			"FirstName":          "Jane",
			"Telephone":          "0123456781",
			"Gender":             models.GenderInput{Gender: models.GenderWoman},
			"EncryptedBirthDate": "1970-01-28",
			"GoogleID":           "google_id1",
			"Apple_ID":           "apple_id1",
		}),
		NewBasicUser(map[string]any{
			"ID":                 "123e4567-e89b-12d3-a456-426614174002",
			"Email":              "jean.doe@example.com",
			"EmailHash":          "jean.doe@example.com",
			"Picture":            "picture2",
			"FirstName":          "Jean",
			"Telephone":          "0123456782",
			"Gender":             models.GenderInput{Gender: models.GenderNonBinary},
			"EncryptedBirthDate": "2000-10-05",
			"GoogleID":           "google_id2",
			"Apple_ID":           "apple_id2",
		}),
	}
	return users
}

func NewBasicUserPendingResponse(overrides map[string]any) *models.UserPendingResponse {
	user := NewBasicUser(nil)
	res := &models.UserPendingResponse{
		Email:    user.Email,
		Role:     user.Role,
		Provider: "email",
	}
	for key, value := range overrides {
		switch key {
		case "Email":
			res.Email = value.(string)
		case "Role":
			res.Role = value.(string)
		case "Provider":
			res.Provider = value.(models.ProviderType)
		}
	}
	return res
}

var (
	Johndoe = &models.User{
		ID:         "1",
		Email:      "john.doe@gmail.com",
		Password:   "$a9rfNhA$N$A78#m",
		CreatedAt:  time.Now().Add(-time.Hour * 4),
		LoggedInAt: time.Now().Add(-time.Hour * 4),
		Role:       models.STANDARD.String(),
		BirthDate:  birthdate,
		LastName:   "DOE",
		FirstName:  "John",
		Gender:     models.GenderInput{Gender: models.GenderMan},
		Telephone:  "0123456789",
	}
	Janedoe = &models.User{
		ID:         "2",
		Email:      "jane.doe@gmail.com",
		Password:   "w4w3f09QF&h)#fwe",
		CreatedAt:  time.Now().Add(-time.Hour * 4),
		LoggedInAt: time.Now().Add(-time.Hour * 4),
		Role:       models.STANDARD.String(),
		BirthDate:  birthdate,
		LastName:   "DOE",
		FirstName:  "Jane",
		Gender:     models.GenderInput{Gender: models.GenderWoman},
		Telephone:  "0123456780",
	}
	Jeandoe = &models.User{
		ID:         "1",
		Email:      "jean.doe@gmail.com",
		Password:   "wf0fT^9f2$$_aewa",
		CreatedAt:  time.Now().Add(-time.Hour * 4),
		LoggedInAt: time.Now().Add(-time.Hour * 4),
		Role:       models.STANDARD.String(),
		BirthDate:  birthdate,
		LastName:   "DOE",
		FirstName:  "Jean",
		Gender:     models.GenderInput{Gender: models.GenderMan},
		Telephone:  "0123456781",
	}
)
