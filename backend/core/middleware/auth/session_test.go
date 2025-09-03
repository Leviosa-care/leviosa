package auth

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeSession(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectError bool
		expectNil   bool
	}{
		{
			name: "valid session JSON",
			input: []byte(`{
				"user_id_encrypted": "dXNlcl9pZA==",
				"role_encrypted": "cm9sZQ==",
				"state_encrypted": "c3RhdGU=",
				"created_at_encrypted": "Y3JlYXRlZF9hdA==",
				"expires_at_encrypted": "ZXhwaXJlc19hdA==",
				"access_token_hash": "test_access_token_hash",
				"dek_encrypted": "ZGVr",
				"key_version": 1
			}`),
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "empty JSON",
			input:       []byte(`{}`),
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "invalid JSON",
			input:       []byte(`{invalid json`),
			expectError: true,
			expectNil:   true,
		},
		{
			name:        "null input",
			input:       []byte(`null`),
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "empty input",
			input:       []byte(``),
			expectError: true,
			expectNil:   true,
		},
		{
			name:        "malformed JSON - missing quotes",
			input:       []byte(`{access_token_hash: test}`),
			expectError: true,
			expectNil:   true,
		},
		{
			name:        "malformed JSON - trailing comma",
			input:       []byte(`{"access_token_hash": "test",}`),
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := DecodeSession(tt.input)

			if tt.expectError {
				assert.Error(t, err, "expected error for input: %s", string(tt.input))
				if tt.expectNil {
					assert.Nil(t, session, "session should be nil on error")
				}
			} else {
				assert.NoError(t, err, "unexpected error for input: %s", string(tt.input))
				assert.NotNil(t, session, "session should not be nil on success")
			}
		})
	}
}

func TestDecodeSession_ValidSession(t *testing.T) {
	// Test with a complete, valid session
	sessionData := map[string]interface{}{
		"user_id_encrypted":     []byte("encrypted_user_id"),
		"role_encrypted":        []byte("encrypted_role"),
		"state_encrypted":       []byte("encrypted_state"),
		"created_at_encrypted":  []byte("encrypted_created_at"),
		"expires_at_encrypted":  []byte("encrypted_expires_at"),
		"access_access_token_hash":     "test_access_access_token_hash_123",
		"refresh_access_token_hash":    "test_refresh_access_token_hash_123",
		"dek_encrypted":         []byte("encrypted_dek"),
		"key_version":           42,
	}

	jsonData, err := json.Marshal(sessionData)
	require.NoError(t, err)

	session, err := DecodeSession(jsonData)
	require.NoError(t, err)
	require.NotNil(t, session)

	// Verify fields are correctly unmarshaled
	assert.Equal(t, "test_access_access_token_hash_123", session.AccessTokenHash)
	assert.Equal(t, "test_refresh_access_token_hash_123", session.RefreshTokenHash)
	assert.Equal(t, 42, session.KeyVersion)
	assert.Equal(t, []byte("encrypted_user_id"), session.UserIDEncrypted)
	assert.Equal(t, []byte("encrypted_role"), session.RoleEncrypted)
	assert.Equal(t, []byte("encrypted_state"), session.StateEncrypted)
	assert.Equal(t, []byte("encrypted_created_at"), session.CreatedAtEncrypted)
	assert.Equal(t, []byte("encrypted_expires_at"), session.ExpiresAtEncrypted)
	assert.Equal(t, []byte("encrypted_dek"), session.DEKEncrypted)
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
				"access_token_hash": "test",
				"key_version": "not_a_number"
			}`,
			expectErr: true,
		},
		{
			name: "key_version as float",
			jsonInput: `{
				"access_token_hash": "test",
				"key_version": 1.5
			}`,
			expectErr: true, // JSON unmarshaling does NOT convert float to int
		},
		{
			name: "access_token_hash as number",
			jsonInput: `{
				"access_token_hash": 123,
				"key_version": 1
			}`,
			expectErr: true,
		},
		{
			name: "encrypted fields as strings",
			jsonInput: `{
				"user_id_encrypted": "should_be_bytes",
				"access_token_hash": "test",
				"key_version": 1
			}`,
			expectErr: true,
		},
		{
			name: "encrypted fields as base64 strings",
			jsonInput: `{
				"user_id_encrypted": "dGVzdA==",
				"access_token_hash": "test",
				"key_version": 1
			}`,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeSession([]byte(tt.jsonInput))

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
	assert.Equal(t, SessionState("pending"), SessionPending)
	assert.Equal(t, SessionState("active"), SessionActive)

	// Test that they're different
	assert.NotEqual(t, SessionPending, SessionActive)

	// Test string conversion
	assert.Equal(t, "pending", string(SessionPending))
	assert.Equal(t, "active", string(SessionActive))
}

func TestSession_FieldTags(t *testing.T) {
	// Test that Session struct has correct JSON tags using reflection
	// This ensures encrypted fields are properly marked for JSON serialization

	session := Session{
		ID:                 uuid.New(),
		UserIDEncrypted:    []byte("encrypted_user"),
		RoleEncrypted:      []byte("encrypted_role"),
		StateEncrypted:     []byte("encrypted_state"),
		CreatedAtEncrypted: []byte("encrypted_created"),
		ExpiresAtEncrypted: []byte("encrypted_expires"),
		AccessTokenHash:          "access_access_token_hash_value",
		RefreshTokenHash:         "refresh_access_token_hash_value",
		DEKEncrypted:       []byte("encrypted_dek"),
		KeyVersion:         1,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(session)
	require.NoError(t, err)

	// Unmarshal back to verify field mapping
	var unmarshaled map[string]interface{}
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Check that encrypted fields are included in JSON
	assert.Contains(t, unmarshaled, "user_id_encrypted")
	assert.Contains(t, unmarshaled, "role_encrypted")
	assert.Contains(t, unmarshaled, "state_encrypted")
	assert.Contains(t, unmarshaled, "created_at_encrypted")
	assert.Contains(t, unmarshaled, "expires_at_encrypted")
	assert.Contains(t, unmarshaled, "access_token_hash")
	assert.Contains(t, unmarshaled, "refresh_token_hash")
	assert.Contains(t, unmarshaled, "dek_encrypted")
	assert.Contains(t, unmarshaled, "key_version")

	// Check that plaintext fields are excluded (json:"-" tag)
	assert.NotContains(t, unmarshaled, "id")
	assert.NotContains(t, unmarshaled, "user_id")
	assert.NotContains(t, unmarshaled, "role")
	assert.NotContains(t, unmarshaled, "state")
	assert.NotContains(t, unmarshaled, "created_at")
	assert.NotContains(t, unmarshaled, "expires_at")
	assert.NotContains(t, unmarshaled, "token")
	assert.NotContains(t, unmarshaled, "dek")
}

func TestSession_Constants(t *testing.T) {
	// Test that session constants are reasonable
	assert.Equal(t, 24*time.Hour, SessionDuration)

	// Verify duration is positive
	assert.Positive(t, SessionDuration)
}

func TestDecodeSession_RoundTrip(t *testing.T) {
	// Test that we can marshal a session and then decode it back
	original := Session{
		ID:                 uuid.New(),
		UserIDEncrypted:    []byte("encrypted_user_id_data"),
		RoleEncrypted:      []byte("encrypted_role_data"),
		StateEncrypted:     []byte("encrypted_state_data"),
		CreatedAtEncrypted: []byte("encrypted_created_at_data"),
		ExpiresAtEncrypted: []byte("encrypted_expires_at_data"),
		AccessTokenHash:          "original_access_access_token_hash",
		RefreshTokenHash:         "original_refresh_access_token_hash",
		DEKEncrypted:       []byte("encrypted_dek_data"),
		KeyVersion:         123,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	require.NoError(t, err)

	// Decode back
	decoded, err := DecodeSession(jsonData)
	require.NoError(t, err)
	require.NotNil(t, decoded)

	// Verify encrypted fields are preserved
	assert.Equal(t, original.UserIDEncrypted, decoded.UserIDEncrypted)
	assert.Equal(t, original.RoleEncrypted, decoded.RoleEncrypted)
	assert.Equal(t, original.StateEncrypted, decoded.StateEncrypted)
	assert.Equal(t, original.CreatedAtEncrypted, decoded.CreatedAtEncrypted)
	assert.Equal(t, original.ExpiresAtEncrypted, decoded.ExpiresAtEncrypted)
	assert.Equal(t, original.AccessTokenHash, decoded.AccessTokenHash)
	assert.Equal(t, original.DEKEncrypted, decoded.DEKEncrypted)
	assert.Equal(t, original.KeyVersion, decoded.KeyVersion)

	// Verify plaintext fields are zero values (not serialized)
	assert.Equal(t, uuid.Nil, decoded.ID)
	assert.Equal(t, uuid.Nil, decoded.UserID)
	assert.Equal(t, identity.Role(0), decoded.Role) // Zero value is 0, not 1
	assert.Equal(t, SessionState(""), decoded.State)
	assert.True(t, decoded.CreatedAt.IsZero())
	assert.True(t, decoded.ExpiresAt.IsZero())
	assert.Empty(t, decoded.AccessToken)
	assert.Empty(t, decoded.RefreshToken)
	assert.Nil(t, decoded.DEK)
}
