# ADR 0010 — Public booking lookup via signed token

**Status:** Accepted  
**Date:** 2026-05-28

The `/bookings` route is a public, unauthenticated booking status page. Its primary access path is a signed token embedded in the booking confirmation SMS and email (`/bookings?token=xxx`). Visitors who arrive without a token can look up their booking manually by entering a booking reference and the email or phone used when booking. From this page, guests can view their booking details and cancel — subject to the same cancellation window enforced for authenticated clients.

## Why a public lookup page rather than redirecting to `/client/bookings`

SMS is the primary notification channel. A link sent in a reminder SMS must resolve to something useful regardless of whether the recipient is currently logged in. Redirecting unauthenticated visitors to `/auth` and then to `/client/bookings` would strand the significant share of guests who never created an account. The public lookup page lets any guest act on their booking from a single tap.

## Why a signed token rather than a reference + contact form alone

A token embedded in the SMS removes all manual entry on mobile — the most common context for acting on a booking reminder. The reference + contact form exists as a fallback for guests who no longer have the original message. The token is signed (HMAC over booking ID with a server secret) and expires after the booking date plus a grace period, so a leaked token has bounded impact.

## Why cancel is allowed on this page

Requiring guests to create an account to cancel would create a perverse outcome: a guest who wants to free a slot for someone else is blocked unless they register. The cancellation window (defined on the Product) is enforced by the backend regardless of the caller's auth state — the public page simply exposes the same action with token-based authorization instead of session-based.

## Consequences

- The `booking` service must implement token generation (HMAC-signed, booking-scoped) and verification.
- The confirmation SMS and email must include the token URL.
- The manual fallback lookup (reference + email/phone) requires the booking service to accept unauthenticated queries against encrypted guest fields — decryption must happen server-side before comparison.
- Authenticated clients arriving at `/bookings` are redirected to `/client/bookings`.
