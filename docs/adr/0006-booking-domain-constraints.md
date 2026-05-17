# ADR-0006: Booking Domain Constraints and Scheduling Rules

**Status**: Accepted  
**Date**: 2025-07

---

## Context

The scheduling domain has several non-obvious constraints that must be consistently enforced across APIs, background jobs, and future tooling. Encoding them only in HTTP handlers or database triggers creates drift. They belong in the domain layer.

## Decision

The following constraints are enforced in `internal/booking/domain/` and `internal/catalog/domain/`:

**Product duration**: Must be in `[20, 240]` minutes and a multiple of 10. Validated in the `Product` domain constructor; rejected before reaching persistence.

**10-minute alignment**: All availability slots and booking start times must align to 10-minute boundaries (e.g., 09:00, 09:10, 09:20). This simplifies gap detection and utilization math.

**Buffer time**: After a session ends, the room is unavailable for the product's buffer duration. The booking system accounts for this when checking slot availability.

**Cancellation window**: Bookings cannot be cancelled within N hours of the session start time, where N is the product's `cancellation_window`. Enforced as a domain method on `Booking`.

**Booking status FSM**: 
```
confirmed → completed
confirmed → cancelled
confirmed → no_show
```
Transitions are methods on the `Booking` struct. Invalid transitions (e.g., `completed → cancelled`) return a domain error without hitting the database.

**Room allocation types**:
- `dedicated`: a room is assigned to one partner for a fixed time period; only that partner can book it during that period
- `shared`: a room is available to any allocated partner on an ongoing basis; no fixed end date

**Gap detection**: The application layer (`booking/application/`) scans room schedules to identify unused time between bookings. Gaps shorter than a configurable threshold are surfaced as upsell opportunities. This is a read-only analytical operation, not a booking mutation.

## Consequences

**Positive:**
- Business rules live in one place and are testable without infrastructure
- Domain errors are explicit and typed, not inferred from database constraint violations
- Future APIs (mobile, webhook) automatically inherit the same rules

**Negative:**
- Duration and alignment constraints must be re-documented when onboarding partners who expect arbitrary durations
- The FSM must be extended explicitly for any new booking status
