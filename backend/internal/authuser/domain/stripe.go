package domain

// StripeAccountStatus represents the possible statuses of a Stripe Connected Account
type StripeAccountStatus string

const (
	// StripeAccountStatusPending indicates the account has been created but not yet completed onboarding
	StripeAccountStatusPending StripeAccountStatus = "pending"

	// StripeAccountStatusActive indicates the account is fully onboarded and can receive payments
	StripeAccountStatusActive StripeAccountStatus = "active"

	// StripeAccountStatusRestricted indicates the account has restrictions (needs more information)
	StripeAccountStatusRestricted StripeAccountStatus = "restricted"

	// StripeAccountStatusDisabled indicates the account has been disabled by Stripe
	StripeAccountStatusDisabled StripeAccountStatus = "disabled"
)

// IsValid returns true if the StripeAccountStatus is a valid status
func (s StripeAccountStatus) IsValid() bool {
	switch s {
	case StripeAccountStatusPending, StripeAccountStatusActive, StripeAccountStatusRestricted, StripeAccountStatusDisabled:
		return true
	default:
		return false
	}
}

// String returns the string representation of the StripeAccountStatus
func (s StripeAccountStatus) String() string {
	return string(s)
}

// CanReceivePayments returns true if the account status allows receiving payments
func (s StripeAccountStatus) CanReceivePayments() bool {
	return s == StripeAccountStatusActive
}

// NeedsOnboarding returns true if the account needs to complete or update onboarding
func (s StripeAccountStatus) NeedsOnboarding() bool {
	return s == StripeAccountStatusPending || s == StripeAccountStatusRestricted
}

// IsProblematic returns true if the account has issues that need attention
func (s StripeAccountStatus) IsProblematic() bool {
	return s == StripeAccountStatusRestricted || s == StripeAccountStatusDisabled
}