---
target: frontend/src/routes/(public)
total_score: 24
p0_count: 2
p1_count: 2
timestamp: 2026-07-20T20-48-38Z
slug: frontend-src-routes-public
---
# Critique — Leviosa Public Storefront

**Method: dual-agent (A: a18b28aa87f88f49a · B: a0815e5e1c679a791)**

Target: `frontend/src/routes/(public)/` (homepage + about, services, team, book, bookings, legal). Booking-flow mechanics excluded per scope (separate track). Browser inspection skipped — no dev server running.

## Design Health Score

| # | Heuristic | Score | Key Issue |
|---|-----------|-------|-----------|
| 1 | Visibility of System Status | 3/4 | Mostly SSR'd; no real gaps found |
| 2 | Match Between System and Real World | 1/4 | Mission/values copy reads like B2B startup-consulting, not a wellness practice |
| 3 | User Control and Freedom | 3/4 | Drawer, cancel flow, back links all solid |
| 4 | Consistency and Standards | 3/4 | Tokens consistent; "Nos valeurs" duplicated near-verbatim across two pages |
| 5 | Error Prevention | 2/4 | Cancel gated; booking-lookup contact fields not required |
| 6 | Recognition Rather Than Recall | 2/4 | Manual booking lookup demands a memorized/pasted raw UUID |
| 7 | Flexibility and Efficiency | 3/4 | Booking entry points everywhere |
| 8 | Aesthetic and Minimalist Design | 3/4 | Clean by minimalist standards — but minimalism isn't the brief anymore |
| 9 | Error Recovery | 3/4 | Inline errors with icon + message throughout |
| 10 | Help and Documentation | 1/4 | FAQ solid; legal/privacy and legal/terms are empty files |
| **Total** | | **24/40** | **Acceptable band — two weakest scores are both trust-critical** |

## Anti-Patterns Verdict

LLM: eyebrow-above-every-section default, hero-metric stat-row implemented 3x, identical 3-up icon-card grids repeated 4+ times, one decorative dot-grid texture copy-pasted across pages, sanctioned glass treatment leaked into a second decorative use. No cream/beige backgrounds, no gradient text. Deeper problem: copy for `_why.svelte`, `about` mission/values, `_services.svelte` "Consultation Stratégique" reads like B2B consultancy, not this business.

Deterministic scan: detect.mjs exit code 2, 5 advisory findings, no false positives — 4x off-ramp text-[10px], 1x undocumented hex color. Static-evidence pass (grep) found the real scope of token drift: hardcoded status colors across 7+ files (no status-color tokens exist at all in app.css), duplicated dot-grid inline style across 4 files, eyebrow pattern quantified at 11 occurrences across 8 files.

## Overall Impression

Competently built, not visually broken — but carrying two different businesses' worth of copy, and the brand-energy shift PRODUCT.md just committed to hasn't landed anywhere in the shipped code yet. Two P0s (mismatched copy, empty legal pages) sit exactly at the trust moments PRODUCT.md's belief ladder calls most important, and are currently the weakest-executed parts of the site.

## What's Working

1. `reveal.ts` correctly short-circuits under prefers-reduced-motion (not just softer durations), applied uniformly.
2. Cancel-booking flow gates a destructive action behind explicit confirm + inline error + loading spinner — right friction at a stakes-bearing moment.
3. Dual-CTA card pattern ("Réserver" + "Voir le profil/détails") faithfully implements PRODUCT.md's secondary-CTA strategy.

## Priority Issues

**[P0] Copy is written for the wrong business**
Why it matters: `_why.svelte`, `about`'s mission/values, `_services.svelte`'s "Consultation Stratégique" use B2B startup-consulting language. Violates PRODUCT.md's anti-reference directly, breaks belief ladder at step one.
Fix: rewrite around wellness-appropriate value props; replace generic strategy card with an actual bookable session type.
Suggested command: /impeccable clarify

**[P0] Legal pages are empty**
Why it matters: legal/privacy and legal/terms are 2-line stubs, linked from footer, reachable during payment flow. PRODUCT.md states trust is a floor not a stretch goal.
Fix: populate with real content or a clearly-labeled placeholder with a contact channel.
Suggested command: /impeccable harden

**[P1] Brand personality hasn't moved off the stated floor**
Why it matters: DESIGN.md calls the current system "the floor to build up from, not the ceiling." Eyebrow default (11 occurrences/8 files), stat-row implemented 3x, identical 3-up grids 4+ times, hardcoded status colors (no tokens exist) — default scaffolding, not a considered brand push.
Fix: extend Cal Sans display type beyond _how.svelte, retire uppercase-eyebrow default, restrict glass to hero only, formalize real status-color tokens.
Suggested command: /impeccable bolder

**[P1] Homepage front-loads nine sections before any commitment point**
Why it matters: Hero → How → Services → Stats → Why → Team → Testimonials → FAQ → Ready stacks three redundant "trust us" sections before a decided visitor can act — cuts against "get out of the way of the flow" positioning.
Fix: consolidate repeated trust sections, or lean on persistent nav CTA as explicit fast path.
Suggested command: /impeccable distill

**[P2] Trust-critical content buried or vague exactly where it should be prominent**
Why it matters: cancellation-fee FAQ answer non-committal; verification answer at FAQ position 4 of 6; only "Vérifié" badge on the site is on a fake placeholder person in dead code (_hero.card.svelte, imported nowhere); real team cards carry no verification signal.
Fix: surface concrete verification statement near hero/_how.svelte; add real "Vérifié" indicator to actual team cards; replace vague cancellation line with a specific rule or policy link.
Suggested command: /impeccable clarify

## Persona Red Flags

Jordan (first-timer): hero headline doesn't say expert-in-what until scrolling; hits B2B copy mid-scroll and wonders if on wrong site; "Mes réservations" nav item has no visual distinction for a non-customer.

Alex (project-specific — first-time client deciding whether to trust the site enough to pay): clicks privacy policy pre-payment, lands on blank page, likely abandons; no concrete cancellation policy visible before booking; no verified-practitioner signal on real team roster; asked to retain/retype a raw UUID as "reference number."

Casey (distracted mobile): hero autoplays background video with no poster/fallback for slow connections; mobile drawer nav genuinely well-built (44px targets, Escape support); services/team category sidebar becomes horizontally-scrolling pill row on mobile with no visible scroll affordance.

## Minor Observations

- `_hero.card.svelte` is dead code — imported nowhere, yet holds the site's only "Vérifié" badge.
- `team/[id]/+page.svelte:79` uses border-l-2 on a blockquote — token-based/neutral, soft case of the banned accent-stripe shape.
- FAQ accordion, category filters, and bookings contact-method toggle all lack proper ARIA state.
- Footer omits "Mes réservations," present in header nav.
- Footer has 3 instances of off-ramp text-[10px] type (detector-caught).
- Dot-grid decorative background duplicated as inline style across 4 separate page files.

## Questions to Consider

1. If `_why.svelte` and "Nos valeurs" were deleted tonight, would a wellness client lose anything they needed?
2. Where, concretely, in the shipped code is the stated brand-energy effort visible?
3. If verification is the core trust lever, why doesn't a single real practitioner card carry the "Vérifié" badge?
4. Is the nine-section homepage serving "Doctolib-simple," or a different, unstated goal now in tension with it?
