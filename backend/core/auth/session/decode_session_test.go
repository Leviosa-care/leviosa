package session_test

import (
	"testing"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/stretchr/testify/assert"
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
			sessionStruct, err := session.DecodeSession(tt.input)

			if tt.expectError {
				assert.Error(t, err, "expected error for input: %s", string(tt.input))
				if tt.expectNil {
					assert.Nil(t, sessionStruct, "session should be nil on error")
				}
			} else {
				assert.NoError(t, err, "unexpected error for input: %s", string(tt.input))
				assert.NotNil(t, sessionStruct, "session should not be nil on success")
			}
		})
	}
}

