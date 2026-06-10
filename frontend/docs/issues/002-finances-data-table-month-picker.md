# 002 — Finances page — data table + month picker

**Type:** AFK
**Status:** open
**Blocked by:** None

## What to build

Enhance `/staff/statistics/finances` so transactions are displayed in a scannable data table and the partner can browse any historical month.

**Month picker:** Add a `<input type="month">` control in the page header (matching the pattern used by the admin compta page). Selecting a month re-navigates to the same route with a `?month=YYYY-MM` query parameter. The server load derives its date range from this param, defaulting to the current month when the param is absent.

**Data table:** Replace the current stacked transaction card list with a proper `<table>` with columns: Date/Heure · Prestation · Statut paiement · Statut réservation · Montant. Existing badge colour conventions for payment status (`paid` = green, `pending` = yellow, `refunded` = red/orange) carry over unchanged. The four KPI summary cards at the top of the page remain as-is.

**Empty state:** When the selected month has no transactions, show a centred empty-state (icon + message) consistent with the existing empty states elsewhere on the page.

## Acceptance criteria

- [ ] Month picker renders in the page header with the current month pre-selected
- [ ] Navigating to `?month=2025-11` loads transactions for November 2025
- [ ] Transactions are rendered in a `<table>` with the five specified columns
- [ ] Each row shows the correct payment and booking status badges
- [ ] Selecting a month with no transactions shows an empty state, not a broken table
- [ ] KPI summary cards continue to reflect the selected month's data
- [ ] `npm run check` passes with no new type errors

## Parent

docs/prd/staff-pages-implementation.md
