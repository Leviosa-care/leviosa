# 005 — Staff dashboard — actions widget

**Type:** AFK
**Status:** open
**Blocked by:** 003 — Staff dashboard data layer + agenda widget

## What to build

Replace the hardcoded "Actions requises" widget with real partner-scoped action items, and hide it entirely when there is nothing to act on.

The server load (already extended in prior issues) should also return:
- Whether the partner's Stripe onboarding is incomplete (`stripeOnboardingComplete === false`)
- The count of bookings in `pending` status (awaiting partner confirmation)

`_actions.svelte` receives a typed `actions` array. Each action item has a title, description, an age/timestamp string, and a `href` destination. Supported actions:

- **Stripe onboarding incomplete** — title "Configuration Stripe requise", links to `/staff/profile`.
- **Pending bookings** — title "Réservations en attente" with count in description, links to `/staff/agenda/reservations`. Only shown when count > 0.

When the `actions` array is empty the widget renders nothing (no container, no empty state). The amber pulse dot and hardcoded badge count of `3` are removed. The animated dot reappears only when there are real actions.

## Acceptance criteria

- [ ] Widget shows a "Configuration Stripe requise" item when Stripe onboarding is incomplete
- [ ] Widget shows a pending bookings item when the partner has pending bookings, with the correct count
- [ ] Clicking each action item navigates to the correct destination
- [ ] Widget is completely hidden when there are no action items
- [ ] The amber pulse dot only renders when at least one action is present
- [ ] A failed partner profile fetch results in an empty actions list (widget hidden), not an error state
- [ ] `npm run check` passes with no new type errors

## Parent

docs/prd/staff-pages-implementation.md
