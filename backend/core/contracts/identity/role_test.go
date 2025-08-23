package identity_test

import (
	"testing"

	"github.com/Leviosa-care/core/contracts/identity"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		role     identity.Role
		expected string
		name     string
	}{
		{role: identity.VISITOR, expected: "visitor", name: "Get string unknown"},
		{role: identity.STANDARD, expected: "standard", name: "Get string standard"},
		{role: identity.PREMIUM, expected: "premium", name: "Get string premium"},
		{role: identity.GUEST, expected: "guest", name: "Get string guest"},
		{role: identity.PARTNER, expected: "partner", name: "Get string partner"},
		{role: identity.ADMINISTRATOR, expected: "administrator", name: "Get string administrator"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.role.String()
			assert.Equal(t, got, tt.expected)
		})
	}
}

func TestConvertToRole(t *testing.T) {
	tests := []struct {
		roleStr  string
		expected identity.Role
		name     string
	}{
		{roleStr: "visitor", expected: identity.VISITOR, name: "Convert to VISITOR"},
		{roleStr: "standard", expected: identity.STANDARD, name: "Convert to STANDARD"},
		{roleStr: "premium", expected: identity.PREMIUM, name: "Convert to PREMIUM"},
		{roleStr: "guest", expected: identity.GUEST, name: "Convert to GUEST"},
		{roleStr: "partner", expected: identity.PARTNER, name: "Convert to PARTNER"},
		{roleStr: "administrator", expected: identity.ADMINISTRATOR, name: "Convert to ADMINISTRATOR"},
		// {roleStr: test.GenerateRandomString(5), expected: identity.VISITOR, name: "Convert to VISITOR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := identity.ConvertToRole(tt.roleStr)
			assert.Equal(t, got, tt.expected)
		})
	}
}
