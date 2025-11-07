# Stripe Payment Options for Partner Platform

## Business Model Overview
- Partners provide services to customers through the platform
- Platform takes a commission percentage from each transaction
- Need to charge customers and pay partners their earnings

## Option 1: Stripe Direct Charges (Simple Implementation)

### Payment Flow
1. **Customer books service** → Customer pays platform directly
2. **Platform processes full payment** via Stripe
3. **Platform takes commission** (your percentage)
4. **Platform pays partner** their remaining share via Stripe Connect

### Implementation Requirements
- One Stripe Customer account per customer
- One Stripe Connected Account per partner
- Manual payout processing from platform to partners
- More control over payment timing and commission calculation

### Advantages
- Simpler to implement initially
- Full control over payment flow
- Can batch process partner payouts
- Easier to handle refunds and disputes

### Disadvantages
- More trust required from partners (platform holds money)
- Complex accounting for platform
- Partners wait for payouts rather than instant payments
- Regulatory burden on platform for money transmission

## Option 2: Stripe Connect + Destination Charges (Recommended)

### Payment Flow
1. **Customer books service** → Single transaction charges customer
2. **Money flows instantly**: Customer → Platform (commission) → Partner (remainder)
3. **Stripe handles automatic** commission deduction and partner payout
4. **Partners receive money directly** in their connected account

### Implementation Requirements
- Stripe Connect platform setup
- Each partner needs Connected Account (Express or Custom)
- Each customer needs Customer account
- Destination charges during checkout

### Advantages
- **Instant partner payments** - partners get paid immediately
- **Better partner trust** - money doesn't flow through platform first
- **Simpler accounting** - Stripe handles commission calculation
- **Reduced regulatory burden** - less money transmission by platform
- **Professional partner experience** - similar to Uber, Airbnb model

### Disadvantages
- More complex Stripe setup
- Less control over payment timing
- Refunds can be more complex
- Need robust partner onboarding for Connected Accounts

## Partner Domain Model for Option 2

### Current Partner Domain Structure

```go
type Partner struct {
	UserID           uuid.UUID   `json:"user_id"`
	Bio              string      `json:"bio" encx:"encrypt"`
	Experience       string      `json:"experience" encx:"encrypt"`
	Certifications   []string    `json:"certifications" encx:"encrypt"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	CategoryIDs      []uuid.UUID `json:"category_ids" encx:"encrypt"`
	ProductIDs       []uuid.UUID `json:"product_ids" encx:"encrypt"`

	// Stripe Connect fields for Option 2
	StripeConnectedAccountID   string `json:"stripe_connected_account_id" encx:"encrypt"`
	StripeAccountStatus        string `json:"stripe_account_status"`
	StripeOnboardingComplete   bool   `json:"stripe_onboarding_complete"`
}
```

### Stripe Account Status Management
- **pending**: Partner started onboarding, not completed
- **active**: Partner can receive payments
- **restricted**: Stripe limited account (needs more info)
- **disabled**: Account cannot receive payments

## Implementation Phases

### Phase 1: Partner Onboarding
- Implement partner registration with basic info
- Add Stripe Connected Account creation during onboarding
- Store `StripeConnectedAccountID` and status tracking
- Handle Stripe verification flow

### Phase 2: Payment Processing
- Implement destination charges during booking
- Handle commission calculation via Stripe
- Add webhook handlers for payment events
- Implement refund handling

### Phase 3: Partner Payouts
- Partner dashboard for earnings viewing
- Automatic payouts via Stripe Connect
- Tax document handling
- Dispute resolution workflow

## Next Steps for Domain Refactoring

1. **✅ Remove old Partner implementation** (`internal/authuser/domain/partner.go`)
2. **✅ Create new Partner struct** based on NewPartner with Stripe Connect fields
3. **Update validation methods** for new fields
4. **Update CompletePartner handler** to use new Partner domain
5. **Add Stripe Connected Account creation** to partner onboarding flow

## Considerations
- **GDPR Compliance**: All Stripe IDs and sensitive data must use `encx:"encrypt"`
- **Partner Onboarding**: Need robust UX for Stripe Connect verification
- **Error Handling**: Stripe API failures need proper handling
- **Testing**: Use Stripe test mode for development and testing
- **Legal**: Consider terms of service and payment processing agreements