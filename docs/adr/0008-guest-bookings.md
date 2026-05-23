# ADR 0008 — Guest bookings without authentication

**Status:** Accepted  
**Date:** 2026-05-23

## Context

Requiring authentication before booking creates unnecessary friction for first-time users who want a single session. A significant share of wellness clients book once, evaluate, and only then create an account. Forcing registration upfront loses these users.

## Decision

Bookings can be made without an authenticated account. `ClientID` becomes nullable on the `Booking` domain model. A booking is valid if it has either a `ClientID` (authenticated user) or guest contact fields:

- `GuestFirstName` + `GuestLastName` (required for guests)
- `GuestEmail` OR `GuestPhone` — at least one required

All guest fields are encrypted at rest via `encx` for GDPR compliance.

After a successful guest booking, the confirmation page shows a prompt to create an account.

## Rationale

- Removes the single biggest friction point in the booking funnel.
- Guest contact info is sufficient for confirmation (email/SMS) and for the partner to identify the client.
- Authenticated users retain advantages: booking history, preferences, future discounts.

## Future work: retroactive booking claim

When a guest later creates an account with the same email or phone number, their past guest bookings should be linked to their new account. This is **not implemented in the first pass** — guest bookings remain standalone until the claim feature is built. Track as a separate issue: "Link guest bookings to account on registration".

## Trade-offs

- Increases complexity of the `Booking` domain (nullable `ClientID`, two valid identity paths).
- `GET /clients/{clientId}/bookings` only returns bookings for authenticated users — guest bookings are not queryable by the guest after the fact until claim is implemented.
- Partner and admin views must handle bookings where client info comes from guest fields rather than the user table.
