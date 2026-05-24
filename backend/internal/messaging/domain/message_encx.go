package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/hengadev/errsx"
)

// MessageEncx is the at-rest encrypted form stored in the database.
type MessageEncx struct {
	ID            uuid.UUID
	ThreadID      uuid.UUID
	SenderID      uuid.UUID
	BodyEncrypted []byte
	CreatedAt     time.Time
	ReadAt        *time.Time

	DEKEncrypted []byte
	KeyVersion   int
}

// ProcessMessageEncx encrypts the message body using a freshly generated DEK.
func ProcessMessageEncx(ctx context.Context, crypto encx.CryptoService, source *Message) (*MessageEncx, error) {
	var errs errsx.Map

	result := &MessageEncx{
		ID:        source.ID,
		ThreadID:  source.ThreadID,
		SenderID:  source.SenderID,
		CreatedAt: source.CreatedAt,
		ReadAt:    source.ReadAt,
	}

	dek, err := crypto.GenerateDEK()
	if err != nil {
		errs.Set("DEK generation", err)
		return result, errs.AsError()
	}

	bodyBytes, err := encx.SerializeValue(source.Body)
	if err != nil {
		errs.Set("Body serialization", err)
	} else {
		result.BodyEncrypted, err = crypto.EncryptData(ctx, bodyBytes, dek)
		if err != nil {
			errs.Set("Body encryption", err)
		}
	}

	result.DEKEncrypted, err = crypto.EncryptDEK(ctx, dek)
	if err != nil {
		errs.Set("DEK encryption", err)
	}

	result.KeyVersion, err = crypto.GetCurrentKEKVersion(ctx, crypto.GetAlias())
	if err != nil {
		errs.Set("KEK version retrieval", err)
	}

	return result, errs.AsError()
}

// DecryptMessageEncx decrypts a MessageEncx back into a plain Message.
func DecryptMessageEncx(ctx context.Context, crypto encx.CryptoService, source *MessageEncx) (*Message, error) {
	var errs errsx.Map

	result := &Message{
		ID:        source.ID,
		ThreadID:  source.ThreadID,
		SenderID:  source.SenderID,
		CreatedAt: source.CreatedAt,
		ReadAt:    source.ReadAt,
	}

	dek, err := crypto.DecryptDEKWithVersion(ctx, source.DEKEncrypted, source.KeyVersion)
	if err != nil {
		errs.Set("DEK decryption", err)
		return result, errs.AsError()
	}

	if len(source.BodyEncrypted) > 0 {
		bodyBytes, err := crypto.DecryptData(ctx, source.BodyEncrypted, dek)
		if err != nil {
			errs.Set("Body decryption", err)
		} else if err := encx.DeserializeValue(bodyBytes, &result.Body); err != nil {
			errs.Set("Body deserialization", err)
		}
	}

	return result, errs.AsError()
}
