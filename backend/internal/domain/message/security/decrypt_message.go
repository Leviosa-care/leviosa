package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hengadev/leviosa/internal/domain/message/models"

	"github.com/hengadev/errsx"
)

// DecryptMessage decrypts sensitive fields in the provided message model and populates them with their decrypted values.
//
// Parameters:
//   - message: A pointer to the `models.Message` struct containing encrypted fields that need to be decrypted.
//
// Returns:
//   - error: An error from a map containing errors for any decryption failures. The map contains key-value pairs
//     where the key is the name of the field (e.g., "encrypted birthdate") and the value is the corresponding error.
//     If no errors occur, an empty map is returned.
func (s *SecureMessageData) DecryptMessage(message *models.Message) error {
	var errs errsx.Map
	timeFields := []struct {
		name           string
		value          *string
		decryptedValue *time.Time
	}{
		{name: "createdAt", value: &message.EncryptedCreatedAt, decryptedValue: &message.CreatedAt},
	}

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
		{"lastname", &message.Content},
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

	return errs.AsError()
}

// decrypt decrypts sensitive data
func (s *SecureMessageData) decrypt(encryptedData string) (string, error) {
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
