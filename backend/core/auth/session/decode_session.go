package session

import (
	"encoding/json"

	"github.com/Leviosa-care/core/errs"
)

// DecodeSession unmarshals JSON bytes back to a session
// This function works with the new encrypted SessionEncx structure
func DecodeSession(data []byte) (*SessionEncx, error) {
	// Decode as SessionEncx (the encrypted version)
	var sessionEncx SessionEncx
	err := json.Unmarshal(data, &sessionEncx)
	if err != nil {
		return nil, errs.NewJSONUnmarshalErr(err)
	}

	// Return the SessionEncx struct
	// The caller will need to use DecryptSessionEncx if they need the decrypted values
	return &sessionEncx, nil
}
