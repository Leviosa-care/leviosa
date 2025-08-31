package identity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		role     Role
		expected string
		name     string
	}{
		{role: Visitor, expected: "visitor", name: "Get string visitor"},
		{role: Standard, expected: "standard", name: "Get string standard"},
		{role: Premium, expected: "premium", name: "Get string premium"},
		{role: Guest, expected: "guest", name: "Get string guest"},
		{role: Partner, expected: "partner", name: "Get string partner"},
		{role: Administrator, expected: "administrator", name: "Get string administrator"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.role.String()
			assert.Equal(t, tt.expected, got)
		})
	}

	t.Run("invalid role should return unknown", func(t *testing.T) {
		invalidRole := Role(99)
		assert.Equal(t, "unknown", invalidRole.String())
	})
}

func TestIsValid(t *testing.T) {
	t.Run("valid roles", func(t *testing.T) {
		validRoles := []Role{Visitor, Standard, Premium, Guest, Partner, Administrator}
		for _, role := range validRoles {
			t.Run(role.String(), func(t *testing.T) {
				assert.True(t, role.IsValid())
			})
		}
	})

	t.Run("invalid roles", func(t *testing.T) {
		invalidRoles := []Role{Role(-1), Role(99), Role(100)}
		for _, role := range invalidRoles {
			t.Run("invalid_role", func(t *testing.T) {
				assert.False(t, role.IsValid())
			})
		}
	})
}

func TestIsAdmin(t *testing.T) {
	t.Run("Administrator should be admin", func(t *testing.T) {
		assert.True(t, Administrator.IsAdmin())
	})

	t.Run("Non-administrator roles should not be admin", func(t *testing.T) {
		nonAdminRoles := []Role{Visitor, Standard, Premium, Guest, Partner}
		for _, role := range nonAdminRoles {
			t.Run(role.String(), func(t *testing.T) {
				assert.False(t, role.IsAdmin())
			})
		}
	})
}