package webhookHandler

import (
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/stretchr/testify/assert"
)

func TestMapStripeAccountStatus(t *testing.T) {
	tests := []struct {
		name           string
		chargesEnabled bool
		payoutsEnabled bool
		expected       domain.StripeAccountStatus
	}{
		{
			name:           "active when both charges and payouts enabled",
			chargesEnabled: true,
			payoutsEnabled: true,
			expected:       domain.StripeAccountStatusActive,
		},
		{
			name:           "restricted when charges enabled but payouts disabled",
			chargesEnabled: true,
			payoutsEnabled: false,
			expected:       domain.StripeAccountStatusRestricted,
		},
		{
			name:           "disabled when both charges and payouts disabled",
			chargesEnabled: false,
			payoutsEnabled: false,
			expected:       domain.StripeAccountStatusDisabled,
		},
		{
			name:           "disabled when charges disabled but payouts enabled (unusual)",
			chargesEnabled: false,
			payoutsEnabled: true,
			expected:       domain.StripeAccountStatusDisabled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapStripeAccountStatus(tt.chargesEnabled, tt.payoutsEnabled)
			assert.Equal(t, tt.expected, result)
		})
	}
}
