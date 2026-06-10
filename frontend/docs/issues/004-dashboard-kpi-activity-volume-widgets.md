# 004 — Staff dashboard — KPI cards, activity feed, volume chart

**Type:** AFK
**Status:** open
**Blocked by:** 003 — Staff dashboard data layer + agenda widget

## What to build

Extend the server load created in issue 003 and wire three remaining dashboard widgets to real data.

**Server load extensions:**
- Fetch 7-day partner metrics from the existing `/partners/{id}/metrics` endpoint (same one used by the analytics sub-page, with a 7-day window) to derive revenue, booking count, and occupation rate for the KPI cards and the volume chart.
- Derive a recent activity feed from the last N bookings/payments already available in the metrics or reservations response — no new endpoint required. Each event has a title (e.g. "Nouvelle réservation"), a subtitle (client name or amount), a relative timestamp, and a semantic colour key.

**Widget refactors:**
- `_cards.svelte` — accepts `{ revenueCents, revenueGrowthPct, bookingsCount, occupationPct }` props. Fields that cannot be derived from available API data (satisfaction score) render `—`. Trend labels render `—` when the previous period value is zero.
- `_activity.svelte` — accepts `{ events: ActivityEvent[] }` props. Shows an empty state when the array is empty.
- `_volume.svelte` — accepts `{ days: { label: string; pct: number }[] }` props (7 entries, normalised 0–100 relative to the week's maximum). Labels are single-letter day abbreviations. When no data is available, all bars render at 0%.

## Acceptance criteria

- [ ] KPI cards display real revenue and booking count for the past 7 days
- [ ] KPI cards display `—` for satisfaction score (no data source yet)
- [ ] Activity feed shows real recent booking/payment events
- [ ] Activity feed shows an empty state when there are no recent events
- [ ] Volume chart bars reflect actual daily booking counts for the past 7 days
- [ ] A metrics API failure leaves the affected widget in a graceful empty/zero state
- [ ] `npm run check` passes with no new type errors

## Parent

docs/prd/staff-pages-implementation.md
