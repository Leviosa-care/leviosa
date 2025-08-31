package identity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleConstants(t *testing.T) {
	t.Run("string constants should match role strings", func(t *testing.T) {
		assert.Equal(t, "visitor", VisitorStr)
		assert.Equal(t, "standard", StandardStr)
		assert.Equal(t, "premium", PremiumStr)
		assert.Equal(t, "guest", GuestStr)
		assert.Equal(t, "partner", PartnerStr)
		assert.Equal(t, "administrator", AdministratorStr)
	})

	t.Run("role string should match constants", func(t *testing.T) {
		assert.Equal(t, VisitorStr, Visitor.String())
		assert.Equal(t, StandardStr, Standard.String())
		assert.Equal(t, PremiumStr, Premium.String())
		assert.Equal(t, GuestStr, Guest.String())
		assert.Equal(t, PartnerStr, Partner.String())
		assert.Equal(t, AdministratorStr, Administrator.String())
	})
}