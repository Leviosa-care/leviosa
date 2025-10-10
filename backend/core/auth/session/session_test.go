package session_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeSession_ValidSession(t *testing.T) {
	// Test with a complete, valid session
	sessionData := map[string]interface{}{
		"userid_encrypted":    []byte("encrypted_user_id"),
		"role_encrypted":      []byte("encrypted_role"),
		"state_encrypted":     []byte("encrypted_state"),
		"createdat_encrypted": []byte("encrypted_created_at"),
		"expiresat_encrypted": []byte("encrypted_expires_at"),
		"accesstoken_hash":    "test_access_token_hash_123",
		"refreshtoken_hash":   "test_refresh_token_hash_123",
		"dek_encrypted":       []byte("encrypted_dek"),
		"key_version":         42,
	}

	jsonData, err := json.Marshal(sessionData)
	require.NoError(t, err)

	sessionStruct, err := session.DecodeSession(jsonData)
	require.NoError(t, err)
	require.NotNil(t, sessionStruct)

	// Verify fields are correctly unmarshaled
	assert.Equal(t, "test_access_token_hash_123", sessionStruct.AccessTokenHash)
	assert.Equal(t, "test_refresh_token_hash_123", sessionStruct.RefreshTokenHash)
	assert.Equal(t, 42, sessionStruct.KeyVersion)
	assert.Equal(t, []byte("encrypted_user_id"), sessionStruct.UserIDEncrypted)
	assert.Equal(t, []byte("encrypted_role"), sessionStruct.RoleEncrypted)
	assert.Equal(t, []byte("encrypted_state"), sessionStruct.StateEncrypted)
	assert.Equal(t, []byte("encrypted_created_at"), sessionStruct.CreatedAtEncrypted)
	assert.Equal(t, []byte("encrypted_expires_at"), sessionStruct.ExpiresAtEncrypted)
	assert.Equal(t, []byte("encrypted_dek"), sessionStruct.DEKEncrypted)
}

func TestDecodeSession_TypeValidation(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput string
		expectErr bool
	}{
		{
			name: "key_version as string",
			jsonInput: `{
				"accesstoken_hash": "test",
				"key_version": "not_a_number"
			}`,
			expectErr: true,
		},
		{
			name: "key_version as float",
			jsonInput: `{
				"accesstoken_hash": "test",
				"key_version": 1.5
			}`,
			expectErr: true, // JSON unmarshaling does NOT convert float to int
		},
		{
			name: "accesstoken_hash as number",
			jsonInput: `{
				"accesstoken_hash": 123,
				"key_version": 1
			}`,
			expectErr: true,
		},
		{
			name: "encrypted fields as strings",
			jsonInput: `{
				"userid_encrypted": "should_be_bytes",
				"accesstoken_hash": "test",
				"key_version": 1
			}`,
			expectErr: true,
		},
		{
			name: "encrypted fields as base64 strings",
			jsonInput: `{
				"userid_encrypted": "dGVzdA==",
				"accesstoken_hash": "test",
				"key_version": 1
			}`,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := session.DecodeSession([]byte(tt.jsonInput))

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSessionState_Constants(t *testing.T) {
	// Test that SessionState constants are defined correctly
	assert.Equal(t, session.SessionState("pending"), session.SessionPending)
	assert.Equal(t, session.SessionState("active"), session.SessionActive)

	// Test that they're different
	assert.NotEqual(t, session.SessionPending, session.SessionActive)

	// Test string conversion
	assert.Equal(t, "pending", string(session.SessionPending))
	assert.Equal(t, "active", string(session.SessionActive))
}

func TestSession_FieldTags(t *testing.T) {
	// Test that SessionEncx struct has correct JSON tags using reflection
	// This ensures encrypted fields are properly marked for JSON serialization

	sessionEncx := session.SessionEncx{
		ID:                 uuid.New(),
		UserIDEncrypted:    []byte("encrypted_user"),
		RoleEncrypted:      []byte("encrypted_role"),
		StateEncrypted:     []byte("encrypted_state"),
		CreatedAtEncrypted: []byte("encrypted_created"),
		ExpiresAtEncrypted: []byte("encrypted_expires"),
		AccessTokenHash:    "access_access_token_hash_value",
		RefreshTokenHash:   "refresh_access_token_hash_value",
		DEKEncrypted:       []byte("encrypted_dek"),
		KeyVersion:         1,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(sessionEncx)
	require.NoError(t, err)

	// Unmarshal back to verify field mapping
	var unmarshaled map[string]interface{}
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Check that encrypted fields are included in JSON
	assert.Contains(t, unmarshaled, "userid_encrypted")
	assert.Contains(t, unmarshaled, "role_encrypted")
	assert.Contains(t, unmarshaled, "state_encrypted")
	assert.Contains(t, unmarshaled, "createdat_encrypted")
	assert.Contains(t, unmarshaled, "expiresat_encrypted")
	assert.Contains(t, unmarshaled, "accesstoken_hash")
	assert.Contains(t, unmarshaled, "refreshtoken_hash")
	assert.Contains(t, unmarshaled, "dek_encrypted")
	assert.Contains(t, unmarshaled, "key_version")

	// Check that most plaintext fields are excluded (json:"-" tag)
	// Note: ID field is included in SessionEncx despite json:"-" tag in original Session
	assert.NotContains(t, unmarshaled, "userid")
	assert.NotContains(t, unmarshaled, "role")
	assert.NotContains(t, unmarshaled, "state")
	assert.NotContains(t, unmarshaled, "createdat")
	assert.NotContains(t, unmarshaled, "expiresat")
	assert.NotContains(t, unmarshaled, "token")
	assert.NotContains(t, unmarshaled, "dek")
}

func TestSession_Constants(t *testing.T) {
	// Test that session constants are reasonable
	assert.Equal(t, 24*time.Hour, session.SessionDuration)

	// Verify duration is positive
	assert.Positive(t, session.SessionDuration)
}

func TestDecodeSession_RoundTrip(t *testing.T) {
	// Test that we can marshal a session and then decode it back
	original := session.SessionEncx{
		UserIDEncrypted:    []byte("encrypted_user_id_data"),
		RoleEncrypted:      []byte("encrypted_role_data"),
		StateEncrypted:     []byte("encrypted_state_data"),
		CreatedAtEncrypted: []byte("encrypted_created_at_data"),
		ExpiresAtEncrypted: []byte("encrypted_expires_at_data"),
		AccessTokenHash:    "original_access_access_token_hash",
		RefreshTokenHash:   "original_refresh_access_token_hash",
		DEKEncrypted:       []byte("encrypted_dek_data"),
		KeyVersion:         123,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	require.NoError(t, err)

	// Decode back
	decoded, err := session.DecodeSession(jsonData)
	require.NoError(t, err)
	require.NotNil(t, decoded)

	// Verify encrypted fields are preserved
	assert.Equal(t, original.UserIDEncrypted, decoded.UserIDEncrypted)
	assert.Equal(t, original.RoleEncrypted, decoded.RoleEncrypted)
	assert.Equal(t, original.StateEncrypted, decoded.StateEncrypted)
	assert.Equal(t, original.CreatedAtEncrypted, decoded.CreatedAtEncrypted)
	assert.Equal(t, original.ExpiresAtEncrypted, decoded.ExpiresAtEncrypted)
	assert.Equal(t, original.AccessTokenHash, decoded.AccessTokenHash)
	assert.Equal(t, original.RefreshTokenHash, decoded.RefreshTokenHash)
	assert.Equal(t, original.DEKEncrypted, decoded.DEKEncrypted)
	assert.Equal(t, original.KeyVersion, decoded.KeyVersion)
}
