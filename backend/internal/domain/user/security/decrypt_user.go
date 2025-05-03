package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hengadev/leviosa/internal/domain/user/models"

	"github.com/hengadev/errsx"
)

// DecryptUser decrypts sensitive fields in the provided user model and populates them with their decrypted values.
//
// Parameters:
//   - user: A pointer to the `models.User` struct containing encrypted fields that need to be decrypted.
//
// Returns:
//   - errsx.Map: A map containing errors for any decryption failures. The map contains key-value pairs
//     where the key is the name of the field (e.g., "encrypted birthdate") and the value is the corresponding error.
//     If no errors occur, an empty map is returned.
func (s *SecureUserData) DecryptUser(user *models.User) error {
	var errs errsx.Map
	timeFields := []struct {
		name           string
		value          *string
		decryptedValue *time.Time
	}{
		{name: "birthdate", value: &user.EncryptedBirthDate, decryptedValue: &user.BirthDate},
		{name: "createdAt", value: &user.EncryptedCreatedAt, decryptedValue: &user.CreatedAt},
		{name: "loggedInAt", value: &user.EncryptedLoggedInAt, decryptedValue: &user.LoggedInAt}}

	for _, field := range timeFields {
		if field.value != nil && *field.value != "" {
			decrypted, err := s.decrypt(*field.value)
			if err != nil {
				errs.Set(field.name, err)
			}
			parsedTime, err := time.Parse(time.RFC3339, decrypted)
			if err != nil {
				errs.Set(fmt.Sprintf("parsing decrypted %s", field.name), err)
			}
			*field.decryptedValue = parsedTime
			*field.value = ""
		}
	}

	fields := []struct {
		name  string
		value *string
	}{
		{"lastname", &user.LastName},
		{"firstname", &user.FirstName},
		{"gender", &user.Gender},
		{"telephone", &user.Telephone},
		{"postal code", &user.PostalCode},
		{"city", &user.City},
		{"address 1", &user.Address1},
		{"address 2", &user.Address2},
		{"google ID", &user.GoogleID},
		{"apple ID", &user.AppleID},
	}

	for _, field := range fields {
		if *field.value != "" {
			decrypted, err := s.decrypt(*field.value)
			if err != nil {
				errs.Set(field.name, err)
			}
			*field.value = decrypted
		}
	}

	// Decrypt email if present
	if user.Email != "" {
		decrypted, err := s.decrypt(user.Email)
		if err != nil {
			errs.Set("email", err)
		}
		user.Email = decrypted
	}

	return errs.AsError()
}

// decrypt decrypts sensitive data
func (s *SecureUserData) decrypt(encryptedData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.config.EncryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", err
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
