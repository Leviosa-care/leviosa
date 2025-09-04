package identity

import (
	"errors"
	"testing"

	"github.com/Leviosa-care/core/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertToRole(t *testing.T) {
	tests := []struct {
		roleStr  string
		expected Role
		isOk     bool
		name     string
	}{
		{roleStr: "visitor", expected: Visitor, name: "Convert to VISITOR", isOk: true},
		{roleStr: "standard", expected: Standard, name: "Convert to STANDARD", isOk: true},
		{roleStr: "premium", expected: Premium, name: "Convert to PREMIUM", isOk: true},
		{roleStr: "guest", expected: Guest, name: "Convert to GUEST", isOk: true},
		{roleStr: "partner", expected: Partner, name: "Convert to PARTNER", isOk: true},
		{roleStr: "administrator", expected: Administrator, name: "Convert to ADMINISTRATOR", isOk: true},
		{roleStr: "random_value", expected: Visitor, name: "Convert to VISITOR", isOk: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ConvertToRole(tt.roleStr)
			assert.Equal(t, tt.expected, got)
			assert.Equal(t, tt.isOk, ok)
		})
	}
}

func TestParseRole(t *testing.T) {
	t.Run("valid roles should parse correctly", func(t *testing.T) {
		tests := []struct {
			input    string
			expected Role
		}{
			{"visitor", Visitor},
			{"standard", Standard},
			{"premium", Premium},
			{"guest", Guest},
			{"partner", Partner},
			{"administrator", Administrator},
			{"VISITOR", Visitor},   // Case insensitive
			{"Standard", Standard}, // Mixed case
			{"  guest  ", Guest},   // With whitespace
		}

		for _, test := range tests {
			t.Run(test.input, func(t *testing.T) {
				role, err := ParseRole(test.input)
				require.NoError(t, err)
				assert.Equal(t, test.expected, role)
			})
		}
	})

	t.Run("invalid roles should return error", func(t *testing.T) {
		invalidRoles := []string{
			"invalid",
			"",
			"admin", // Not "administrator"
			"user",
			"123",
			"null",
		}

		for _, invalid := range invalidRoles {
			t.Run(invalid, func(t *testing.T) {
				role, err := ParseRole(invalid)
				assert.Error(t, err)
				assert.True(t, errors.Is(err, errs.ErrInvalidValue))
				assert.Equal(t, Visitor, role) // Should default to Visitor on error
			})
		}
	})
}

