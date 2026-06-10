# 006 — Styling consistency pass

**Type:** AFK
**Status:** open
**Blocked by:** 002 — Finances data table + month picker, 004 — Staff dashboard KPI cards / activity / volume widgets, 005 — Staff dashboard actions widget

## What to build

Audit and align design token usage across the staff analytics, finances, and profile sub-pages against the admin reference pages, so the portal feels visually coherent without changing the staff-specific typography style.

Changes to make:

- **Card borders:** Replace `border-border` with `border-border-card` on card elements in the analytics, finances, and profile pages where the admin pages use `border-border-card` for the same element role.
- **Card backgrounds:** Ensure KPI cards use `bg-card` (not `bg-background`) for surfaces that should float above the page background in dark mode, matching admin analytics.
- **Hardcoded colour classes in alert banners:** Replace one-off `bg-red-50 border-red-200` and `bg-amber-50 border-amber-200` patterns in the profile page with semantic token equivalents where they exist in the design system.
- **Do not touch** the staff display typography header pattern (`font-display text-4xl tracking-tight`) — this is intentionally distinct from the admin header style and should remain.
- **Do not touch** the admin pages.

## Acceptance criteria

- [ ] Analytics, finances, and profile pages use `border-border-card` on card components consistently with admin reference
- [ ] KPI cards on analytics and finances pages use `bg-card` for card background
- [ ] Alert banners on the profile page use semantic tokens rather than hardcoded Tailwind colour classes where a suitable token exists
- [ ] Staff display typography header style is unchanged
- [ ] No visual regressions on the dashboard homepage widgets
- [ ] `npm run check` passes with no new type errors

## Parent

docs/prd/staff-pages-implementation.md
