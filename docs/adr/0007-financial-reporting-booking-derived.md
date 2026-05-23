# ADR 0007 — Financial reporting is booking-derived, not Stripe-driven

**Status:** Accepted  
**Date:** 2026-05-23

## Context

The admin compta page needs to display financial KPIs (gross revenue, refunds, net revenue) and a transaction list. Two sources are available: the internal booking table (which tracks `amount_cents` and `payment_status` per booking) and the Stripe API (which tracks charges, payouts, platform fees, and reconciliation data).

## Decision

Financial reporting is derived exclusively from the booking table in the first pass. A new `GET /admin/bookings/financial-summary` endpoint in the booking service aggregates completed bookings by period.

## Rationale

- No additional Stripe API integration layer needed at this stage.
- Consistent with how partner `Earnings` is already defined (computed from bookings, not Stripe).
- Simpler to test and reason about.

## Future improvement

Stripe-level reconciliation (actual payouts, platform fees, chargebacks, payout schedules) is a known gap. When the platform needs proper accounting (e.g., for VAT, revenue recognition, or partner payout audits), this endpoint should be replaced or supplemented with a Stripe Balance/Payout API integration.
