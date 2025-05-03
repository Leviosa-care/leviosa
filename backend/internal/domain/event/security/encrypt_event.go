package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/hengadev/leviosa/internal/domain/event/models"

	"github.com/hengadev/errsx"
)

// EncryptEvent encrypts sensitive fields in the provided user model and clears the original plaintext values.
//
// Parameters:
//   - event: A pointer to the `models.Event` struct containing fields to be encrypted, such as title,
//     begin at, description.
//
// Returns:
//   - errsx.Map: A map containing errors for any encryption failures. The map contains key-value pairs
//     where the key is the name of the field (e.g., "encrypt event field beginAt") and the value is the corresponding error.
//     If no errors occur, an empty map is returned.
func (s *SecureEventData) EncryptEvent(event *models.Event) error {
	var errs errsx.Map
	timeFields := []struct {
		name           string
		value          *time.Time
		encryptedValue *string
	}{
		{name: "beginAt", value: &event.BeginAt, encryptedValue: &event.EncryptedBeginAt},
		{name: "endAt", value: &event.EndAt, encryptedValue: &event.EncryptedEndAt},
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
		{name: "title", value: &event.Title},
		{name: "description", value: &event.Description},
		{name: "city", value: &event.City},
		{name: "postalCode", value: &event.PostalCode},
		{name: "address1", value: &event.Address1},
		{name: "address2", value: &event.Address2},
		{name: "priceID", value: &event.PriceID},
	}

	for _, field := range fields {
		if *field.value != "" {
			encrypted, err := s.encrypt(*field.value)
			if err != nil {
				errs.Set(fmt.Sprintf("encrypt event field %s", field.name), err.Error())
			}
			*field.value = encrypted
		}
	}

	// TODO: not so sure about that one
	listfields := []struct {
		name  string
		value *[]string
	}{
		{name: "productID", value: &event.Products},
		{name: "offerID", value: &event.Offers},
	}

	for _, field := range listfields {
		for _, pid := range *field.value {
			if pid != "" {
				encrypted, err := s.encrypt(pid)
				if err != nil {
					errs.Set("encrypt product ID", err.Error())
				}
				pid = encrypted
			}
		}
	}
	return errs.AsError()
}

// encrypt is a helper function for the EncryptEvent function
func (s *SecureEventData) encrypt(data string) (string, error) {
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
