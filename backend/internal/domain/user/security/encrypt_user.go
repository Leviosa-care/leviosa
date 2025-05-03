package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hengadev/leviosa/internal/domain/user/models"

	"github.com/hengadev/errsx"
)

// EncryptUser encrypts sensitive fields in the provided user model and clears the original plaintext values.
//
// Parameters:
//   - user: A pointer to the `models.User` struct containing fields to be encrypted, such as birthdate,
//     last name, first name, email, and password.
//
// Returns:
//   - errsx.Map: A map containing errors for any encryption failures. The map contains key-value pairs
//     where the key is the name of the field (e.g., "encrypt birthdate") and the value is the corresponding error.
//     If no errors occur, an empty map is returned.
func (s *SecureUserData) EncryptUser(user *models.User) error {
	var errs errsx.Map
	timeFields := []struct {
		name           string
		value          *time.Time
		encryptedValue *string
	}{
		{name: "birthdate", value: &user.BirthDate, encryptedValue: &user.EncryptedBirthDate},
		{name: "createdAt", value: &user.CreatedAt, encryptedValue: &user.EncryptedCreatedAt},
		{name: "loggedInAt", value: &user.LoggedInAt, encryptedValue: &user.EncryptedLoggedInAt},
	}

	for _, field := range timeFields {
		if field.value != nil && !field.value.IsZero() {
			dateStr := field.value.Format(time.RFC3339)
			encrypted, err := s.encrypt(dateStr)
			if err != nil {
				errs.Set(field.name, err.Error())
			}
			*field.value = time.Time{}
			*field.encryptedValue = encrypted
		}
	}

	fields := []struct {
		name  string
		value *string
	}{
		{name: "lastname", value: &user.LastName},
		{name: "firstname", value: &user.FirstName},
		{name: "gender", value: &user.Gender},
		{name: "telephone", value: &user.Telephone},
		{name: "postalCode", value: &user.PostalCode},
		{name: "city", value: &user.City},
		{name: "address1", value: &user.Address1},
		{name: "address2", value: &user.Address2},
		{name: "googleID", value: &user.GoogleID},
		{name: "appleID", value: &user.AppleID},
	}

	for _, field := range fields {
		if *field.value != "" {
			encrypted, err := s.encrypt(*field.value)
			if err != nil {
				errs.Set(fmt.Sprintf("encrypt user field %s", field.name), err.Error())
			}
			*field.value = encrypted
		}
	}
	// Handle email specially - we need both a hash for searching and encrypted value
	if user.Email != "" {
		// create hash for searching in database
		user.EmailHash = HashEmail(user.Email)
		// encrypt actual email for storage
		encrypted, err := s.encrypt(user.Email)
		if err != nil {
			errs.Set("encrypt field", err.Error())
		}
		// user.EncryptedEmail = encrypted
		user.Email = encrypted
	}

	// Hash password if present
	if user.Password != "" {
		hash, err := s.hashPassword(user.Password)
		if err != nil {
			errs.Set("hash password", err)
		}
		// TODO: should I encrypt the password ?
		user.PasswordHash = hash
		user.Password = "" // Clear plain text password
	}

	return errs.AsError()
}

func HashEmail(email string) string {
	emailHash := sha256.Sum256([]byte(strings.ToLower(email)))
	return hex.EncodeToString(emailHash[:])
}

// encrypt is a helper function for the EncryptUser function
func (s *SecureUserData) encrypt(data string) (string, error) {
	var errs errsx.Map

	block, err := aes.NewCipher(s.config.EncryptionKey)
	if err != nil {
		errs.Set("aes create cipher", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		errs.Set("cipher create GCM", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		errs.Set("gcm nonce", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), errs.AsError()
}
