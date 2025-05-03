package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/hengadev/leviosa/internal/domain/message/models"

	"github.com/hengadev/errsx"
)

// EncryptMessage encrypts sensitive fields in the provided message model and clears the original plaintext values.
//
// Parameters:
//   - message: A pointer to the `models.Message` struct containing fields to be encrypted, such as birthdate,
//     content.
//
// Returns:
//   - error: An error from a map containing errors for any encryption failures. The map contains key-value pairs
//     where the key is the name of the field (e.g., "encrypt createdAt") and the value is the corresponding error.
//     If no errors occur, an empty map is returned.
func (s *SecureMessageData) EncryptMessage(message *models.Message) error {
	var errs errsx.Map
	timeFields := []struct {
		name           string
		value          *time.Time
		encryptedValue *string
	}{
		{name: "createdAt", value: &message.CreatedAt, encryptedValue: &message.EncryptedCreatedAt},
	}

	for _, field := range timeFields {
		if field.value != nil && !field.value.IsZero() {
			dateStr := field.value.Format(time.RFC3339)
			encrypted, encryptedErrs := s.encrypt(dateStr)
			if len(encryptedErrs) > 0 {
				errs.Set(field.name, encryptedErrs.Error())
			}
			*field.value = time.Time{}
			*field.encryptedValue = encrypted
		}
	}

	fields := []struct {
		name  string
		value *string
	}{
		{name: "content", value: &message.Content},
	}

	for _, field := range fields {
		if *field.value != "" {
			encrypted, err := s.encrypt(*field.value)
			if err != nil {
				errs.Set(fmt.Sprintf("encrypt message field %s", field.name), err)
			}
			*field.value = encrypted
		}
	}

	return errs.AsError()
}

// encrypt is a helper function for the EncryptMessage function
func (s *SecureMessageData) encrypt(data string) (string, error) {
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
