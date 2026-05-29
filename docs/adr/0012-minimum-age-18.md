# ADR-0012: Minimum User Age Set to 18

**Status**: Accepted  
**Date**: 2026-05

---

## Context

The platform requires a minimum age for account creation. The existing code set `MinAgeYears = 13` with a comment citing "GDPR compliance" — 13 is the GDPR floor for data processing consent in most EU member states.

## Decision

The minimum age for account creation is **18 years**. The backend constant `MinAgeYears` is set to `18`. The frontend registration date picker enforces `maxValue = today - 18 years`.

## Rationale

GDPR's age floor governs **data processing consent**, not service eligibility. Leviosa's use case raises two independent concerns:

1. **Contractual capacity**: A booking constitutes a contract between the client and the platform (and implicitly with the partner). In most jurisdictions, minors cannot enter binding contracts independently. Allowing under-18 registrations would require parental consent flows and additional legal safeguards.
2. **Payment**: The booking flow involves Stripe-processed payments. Payment processors and card networks generally require the account holder to be 18+.

13 (the GDPR floor) is the correct minimum for a general-purpose app with no financial transactions. It is not the correct minimum for a paid wellness booking platform.

## Consequences

**Positive:**
- Avoids legal exposure around minor contracts and payment processing
- Simpler onboarding: no parental consent flow needed

**Negative:**
- Excludes 13–17 year olds who might legitimately use wellness services with parental involvement
- If the platform later introduces gift bookings or guardian-managed accounts, this constraint will need revisiting

**Future path**: A guardian-managed account type (adult creates account, books for minor) would restore access for under-18 clients without requiring them to hold their own account.
