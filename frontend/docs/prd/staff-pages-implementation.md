# Staff Pages Implementation

**Status:** Draft
**Date:** 2026-06-10

## Problem Statement

The staff-facing portal (`/staff`) has three groups of pages that are either wired to static mock data or missing navigation access, making them unusable in production:

- The main dashboard homepage displays hardcoded KPI values, agenda entries, activity feed, volume chart, and action items — none connected to the real API.
- The `/staff/statistics/finances` page lists transactions as stacked cards rather than a scannable data table, and has no month-picker, making it inferior to the equivalent admin compta page.
- The `/staff/profile` page is fully API-wired but unreachable from the sidebar; there is no navigation link for partners to access it.

Additionally, the analytics, finances, and profile sub-pages use inconsistent design tokens compared to the admin reference pages, reducing visual coherence across the portal.

## Solution

Connect all staff dashboard widgets to real API data, enhance the finances page with a proper data table and month selector, add a profile entry to the sidebar, and align token usage across all three sub-pages with the admin design reference.

## User Stories

1. As a partner, I want the dashboard KPI cards to show my real revenue, booking count, and occupation rate so I can assess my week at a glance.
2. As a partner, I want the "Agenda du jour" widget to show my actual bookings for today so I know what sessions are coming up.
3. As a partner, I want the "Activité récente" widget to reflect real recent events (new bookings, payments received) so I stay informed without navigating away.
4. As a partner, I want the weekly volume chart to reflect my actual booking volume for the past 7 days so I can spot patterns.
5. As a partner, I want the "Actions requises" widget to surface actionable items specific to me (e.g. incomplete Stripe onboarding, pending booking confirmations) so I can act on them immediately.
6. As a partner, I want the finances page to present transactions in a scannable table with date, product, client, payment status, booking status, and amount columns so I can reconcile payments efficiently.
7. As a partner, I want a month picker on the finances page so I can review earnings for any historical month.
8. As a partner, I want a "Mon profil" link in the sidebar so I can reach my profile page without guessing the URL.
9. As a partner, I want the analytics, finances, and profile pages to feel visually consistent with the rest of the staff portal so the experience is coherent.

## Implementation Decisions

### Module 1 — Staff Dashboard Data Layer
Create a `+page.server.ts` for the main staff homepage. It should fetch in parallel:
- Today's bookings from the reservations API (for the agenda widget)
- Weekly KPI summary (revenue, booking count, occupation rate) — derived from the existing partner metrics endpoint or a dedicated summary endpoint
- Recent activity events (last N bookings/payments) for the activity feed
- Weekly booking volume per day (7 values) for the bar chart
- Pending partner actions: incomplete Stripe onboarding flag (from partner profile) and pending booking confirmations count

All fetch failures should be non-fatal: each widget degrades gracefully to an empty state rather than breaking the page.

### Module 2 — Dashboard Widget Components (prop-driven refactor)
Refactor each widget to accept typed props instead of hardcoded literals:

- `_cards.svelte` — accepts `{ revenue, revenueLabel, bookings, occupation, satisfactionScore }` with sensible `null` fallbacks rendered as `—`.
- `_agenda.svelte` — accepts `{ slots: TodaySlot[] }` where each slot has start/end times, product name, client name, and status.
- `_activity.svelte` — accepts `{ events: ActivityEvent[] }` where each event has a title, subtitle, timestamp, and colour key.
- `_volume.svelte` — accepts `{ days: { label: string; pct: number }[] }` (7 entries, normalised to 0–100).

The parent `+page.svelte` spreads data from the server load into each widget component.

### Module 3 — `_actions.svelte` Widget (partner-scoped)
Scope required actions to the authenticated partner:
- Show a Stripe onboarding prompt when `stripeOnboardingComplete === false`.
- Show a count of bookings in `pending` status awaiting partner confirmation with a link to `/staff/agenda/reservations`.
- Hide admin-level actions (practitioner validation, client disputes). This widget is only meaningful for partners; administrators can rely on the admin portal.

### Module 4 — Finances Page Enhancement
Replace the stacked-card transaction list with an `<table>` component consistent with admin compta:
- Columns: Date/Heure · Prestation · Statut paiement · Statut réservation · Montant
- Badge colouring for payment status (`paid` = green, `pending` = yellow, `refunded` = red/orange) matching existing badge conventions.
- Add a month-picker (`<input type="month">`) in the page header that re-navigates with `?month=YYYY-MM` query params, defaulting to the current month. The server load derives the date range from this param.
- The four KPI summary cards (current month revenue, last month revenue, pending amount, next payout) remain unchanged.

### Module 5 — Profile Sidebar Navigation
Add a `UserCircle` (or `User`) icon entry for `/staff/profile` labelled "Mon profil" to both the `desktopNavigation` and `mobileNavigation` arrays in `sidebar.svelte`. Place it between the statistics entries and the settings icon. Role filter: `["administrator", "partner"]`.

### Module 6 — Styling Consistency Pass
Audit and align the three sub-pages against the admin reference:
- Replace ad-hoc `border-border` on card elements with `border-border-card` where the admin uses it.
- Ensure KPI cards use `bg-card` (not `bg-background`) for surfaces that should float above the page background in dark mode, consistent with admin analytics.
- Remove any remaining hardcoded colour classes (`bg-red-50`, `bg-amber-50`) in favour of semantic token equivalents where they exist.
- Keep the staff-specific display typography header pattern (`font-display text-4xl`) — this is intentionally different from admin and should not be homogenised.

## Testing Decisions

No automated test framework is configured for the frontend. Verification strategy:
- `npm run check` (SvelteKit type-checking) must pass after each module.
- Manual smoke-test in the dev environment: navigate to each page, confirm widgets show real data, confirm graceful empty states when the API returns no results.
- The finances month-picker should be tested across a month boundary (switching from a month with data to one without) to confirm the empty state renders correctly.

## Out of Scope

- Adding chart libraries (e.g. Chart.js, Recharts) — the existing CSS bar chart approach is sufficient.
- Pagination on the finances transaction table — load all transactions for the selected month.
- Satisfaction score KPI — no satisfaction/rating API exists; the card should display `—` until the feature is built.
- Admin-role action items in `_actions.svelte` — those belong in the admin portal.
- Any changes to the admin portal pages.

## Further Notes

- The partner metrics endpoint (`/partners/{id}/metrics`) is already used by the analytics sub-page. The dashboard data layer should reuse the same endpoint with a 7-day window rather than introducing a new one.
- The `_actions.svelte` widget currently has a hardcoded badge count of `3` with an amber pulse dot. Once real data is connected, hide the widget entirely when action count is zero to avoid false urgency.
- Profile page already handles the Stripe onboarding prompt inline — the dashboard `_actions.svelte` entry should link to `/staff/profile` rather than duplicating the onboarding flow.
