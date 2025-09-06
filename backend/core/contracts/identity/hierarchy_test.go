package identity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleLevel(t *testing.T) {
	tests := []struct {
		role          Role
		expectedLevel int
		name          string
	}{
		{Visitor, 0, "visitor has level 0"},
		{Guest, 1, "guest has level 1"},
		{Standard, 2, "standard has level 2"},
		{Premium, 3, "premium has level 3"},
		{Partner, 4, "partner has level 4"},
		{Administrator, 5, "administrator has level 5"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedLevel, test.role.Level())
		})
	}

	t.Run("invalid role should return -1", func(t *testing.T) {
		invalidRole := Role(99)
		assert.Equal(t, -1, invalidRole.Level())
	})
}

func TestIsAtLeast(t *testing.T) {
	t.Run("role hierarchy comparisons", func(t *testing.T) {
		tests := []struct {
			current  Role
			target   Role
			expected bool
			name     string
		}{
			{Administrator, Visitor, true, "admin >= visitor"},
			{Administrator, Administrator, true, "admin >= admin"},
			{Premium, Standard, true, "premium >= standard"},
			{Standard, Premium, false, "standard < premium"},
			{Guest, Administrator, false, "guest < admin"},
			{Visitor, Guest, false, "visitor < guest"},
			{Partner, Premium, true, "partner >= premium"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				assert.Equal(t, test.expected, test.current.IsAtLeast(test.target))
			})
		}
	})
}

func TestHasPermission(t *testing.T) {
	t.Run("permission checks should work like IsAtLeast", func(t *testing.T) {
		assert.True(t, Administrator.HasPermission(Standard))
		assert.True(t, Premium.HasPermission(Guest))
		assert.False(t, Guest.HasPermission(Premium))
		assert.True(t, Partner.HasPermission(Partner))
	})
}

