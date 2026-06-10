# 003 — Staff dashboard data layer + agenda widget

**Type:** AFK
**Status:** open
**Blocked by:** None

## What to build

Create the server data layer for the `/staff` homepage and wire the first widget — the agenda — to real data. This slice establishes the prop-passing pattern that issues 004 and 005 will follow.

**Data layer:** Create `+page.server.ts` for the staff homepage. It must fetch today's bookings from the reservations API (filtered to the authenticated partner, for today's date). All fetches must be non-fatal: a failed API call returns an empty fallback rather than throwing, so no single widget failure can break the whole dashboard.

**Agenda widget refactor:** Convert `_agenda.svelte` from hardcoded literals to a prop-driven component. It receives a typed array of today's time slots (start time, end time, product name, client name, status). The parent `+page.svelte` passes the server data down as props. The visual design and timeline layout remain unchanged. When there are no slots for today, show an empty state message.

The other four widgets (`_cards.svelte`, `_activity.svelte`, `_volume.svelte`, `_actions.svelte`) are intentionally left mocked in this slice — they will be wired in issues 004 and 005.

## Acceptance criteria

- [ ] `+page.server.ts` exists for the staff homepage and fetches today's bookings
- [ ] A fetch failure returns an empty array without breaking the page
- [ ] `_agenda.svelte` accepts a `slots` prop instead of hardcoded entries
- [ ] Real bookings for today are rendered in the agenda widget
- [ ] When the partner has no bookings today, a "Aucune séance aujourd'hui" empty state is shown
- [ ] Remaining widgets still render (even if still mocked)
- [ ] `npm run check` passes with no new type errors

## Parent

docs/prd/staff-pages-implementation.md
