---
name: Leviosa
description: Appointment booking storefront and staff dashboard for a service-based wellness business
colors:
  ink: "#18181b"
  charcoal-text: "#171717"
  charcoal-text-soft: "#525252"
  paper: "#ffffff"
  mist: "#f4f4f5"
  cloud: "#ededed"
  cloud-hover: "#e0e0e0"
  cloud-active: "#d4d4d4"
  hairline-white: "#fafafa"
  sky-pale: "#e1f3fe"
  sky-deep: "#082f49"
  rose-alert: "#e21d48"
  confirm-green: "#21c45d"
  amber-highlight: "#f59e0b"
  contrast-black: "#000000"
typography:
  display:
    fontFamily: "Cal Sans, DM Sans, sans-serif"
    fontWeight: 600
    letterSpacing: "normal"
  hero:
    fontFamily: "DM Sans, sans-serif"
    fontSize: "clamp(2.25rem, 5vw, 3.75rem)"
    fontWeight: 500
    lineHeight: 1.1
    letterSpacing: "-0.025em"
  body:
    fontFamily: "DM Sans, sans-serif"
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.625
  mono:
    fontFamily: "Source Code Pro, monospace"
rounded:
  button: "5px"
  input: "9px"
  card-sm: "10px"
  card: "16px"
  card-lg: "20px"
  pill: "9999px"
spacing:
  input: "3rem"
  input-sm: "2.5rem"
components:
  button-primary:
    backgroundColor: "{colors.ink}"
    textColor: "#ffffff"
    rounded: "{rounded.button}"
    padding: "0.5rem 1rem"
  button-primary-hover:
    backgroundColor: "#18181be6"
  button-outline:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.ink}"
    rounded: "{rounded.button}"
  button-destructive:
    backgroundColor: "{colors.rose-alert}"
    textColor: "#ffffff"
    rounded: "{rounded.button}"
  card-default:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.charcoal-text}"
    rounded: "{rounded.card}"
    padding: "1.5rem"
  input-default:
    backgroundColor: "{colors.paper}"
    rounded: "{rounded.input}"
    height: "2.5rem"
---

# Design System: Leviosa

## 1. Overview

**Creative North Star: "The Studio Drop"**

Leviosa runs two surfaces on one token system today, and they're headed in different directions. The staff/admin/client dashboard is a plain, traditional working tool — it stays that way; no branding ambition needed there. The public storefront is where the effort goes: booking should feel as frictionless as Doctolib, but the brand voice is shifting from "calm wellness portal" toward something with more edge — professional and trustworthy enough to hand over your time and payment, but confident and current rather than sedate, closer to a streetwear campaign than a spa brochure.

The tokens extracted below are the honest baseline as they exist in code right now: a near-grayscale neutral system (true zinc/charcoal, not a warm-tinted "AI cream") with a pale sky-blue accent and functional red/green/amber. The public hero already leans bolder than pure "calm" — full-bleed video, dark overlay, oversized type — but the rest of the storefront (cards, forms, badges) still reads soft and restrained. Closing that gap toward "The Studio Drop" energy is active brand work, not yet finished; treat this file as the floor to build up from, not the ceiling.

This system explicitly rejects a cold, clinical healthcare-portal feel and a generic corporate-SaaS-marketing look — per PRODUCT.md's anti-references.

**Key Characteristics:**
- True neutral grayscale (zinc-derived), not a warm/cream-tinted "AI default" base
- One functional accent (pale sky blue) used sparingly, not as a brand color
- Soft, tactile elevation: gentle shadows and rounded corners (9–20px) everywhere
- Admin/staff dashboard: traditional, utilitarian, unbranded by design
- Public storefront: currently restrained; directionally moving toward bolder, more energetic branding

## 2. Colors

Deliberately restrained: the palette is carried by near-black ink and white paper, with color reserved for function (errors, success, one soft accent) rather than decoration.

### Primary
- **Ink** (`#18181b`): the dominant dark tone — primary button backgrounds, body text on light surfaces, borders. Doubles as `--dark` and `--border` in code; this is the brand's real "color," not gray filler.

### Secondary
- **Sky Pale** (`#e1f3fe`) / **Sky Deep** (`#082f49`): the sole accent pair. Sky Pale as a soft background (badges, success-adjacent highlights), Sky Deep as its paired foreground text. Used sparingly — this is not the brand's signature hue, just a functional highlight.

### Tertiary
- **Amber Highlight** (`#f59e0b`): reserved for warning states and the one "tertiary" callout role (`--tertiary` in code). Not used decoratively.

### Neutral
- **Charcoal Text** (`#171717`): primary body text on paper.
- **Charcoal Text Soft** (`#525252`): secondary/muted text.
- **Paper** (`#ffffff`): base background.
- **Mist** (`#f4f4f5`): muted surface fill (badges, subtle backgrounds).
- **Cloud / Cloud Hover / Cloud Active** (`#ededed` / `#e0e0e0` / `#d4d4d4`): the surface interaction ramp — resting, hover, active/pressed states for non-card surfaces.
- **Hairline White** (`#fafafa`): section-level background separation, barely distinguishable from Paper by design.

### Functional
- **Rose Alert** (`#e21d48`): errors and destructive actions only.
- **Confirm Green** (`#21c45d`): success and confirmation states only.

Full dark-mode equivalents exist in `app.css` (`.dark` scope) — same roles, inverted lightness, same hue anchors. Any new component must define both.

### Named Rules
**The Function-Only Color Rule.** Outside of Ink and Paper, every color in this system maps to a state (error, success, warning, subtle highlight) — never to decoration. If a new UI element wants a color "because it looks nice," that's the signal the brand hasn't earned that color yet; push for boldness through type, layout, or motion first.

## 3. Typography

**Display Font:** Cal Sans (with DM Sans, sans-serif fallback)
**Body Font:** DM Sans (with sans-serif fallback)
**Label/Mono Font:** Source Code Pro (with monospace fallback)

**Character:** Cal Sans is a single self-hosted weight (600) — used deliberately, not as a full family, which keeps display moments rare and considered rather than everywhere. DM Sans carries everything else: a clean, geometric-humanist workhorse that reads professional without corporate stiffness — the base this system builds "confident and energetic" on top of.

### Hierarchy
- **Display** (Cal Sans, 600, size varies by context): section-level headings inside the staff dashboard and select public sections (`_how.svelte`). Reserved for moments that need a distinct brand mark, not routine headings.
- **Hero** (DM Sans, 500, `clamp(2.25rem, 5vw, 3.75rem)`, line-height 1.1, tracking -0.025em): the public homepage hero headline — the boldest type moment on the storefront today.
- **Body** (DM Sans, 400, 1rem, line-height 1.625): standard paragraph text; cap prose width at 65–75ch per general typography guidance.
- **Label** (DM Sans, 500, text-sm): buttons, form labels, nav items.

### Named Rules
**The One Display Rule.** Cal Sans is loaded in a single weight and appears only at section-heading scale — never body copy, never buttons. Its rarity is what makes it read as a brand mark instead of a font choice.

## 4. Elevation

Elevation is soft and functional, not decorative: shadows exist to separate interactive surfaces from the page, and to signal state change (hover, focus), not to add visual weight for its own sake.

### Shadow Vocabulary
- **Mini** (`0px 1px 0px 1px rgba(0,0,0,0.04)`): subtle elevation for inputs, inline chips, small cards.
- **Popover** (`0px 7px 12px 3px var(--dark-10)`): floating elements — dropdowns, modals, tooltips, and the "elevated" card variant.
- **Card** (`0px 2px 0px 1px rgba(0,0,0,0.04)`): standard resting card elevation.
- **Btn** (`0px 1px 0px 1px rgba(0,0,0,0.03)`): button press depth.
- **Kbd** (`0px 2px 0px 0px rgba(0,0,0,0.07)`): keyboard-shortcut chip styling.

### Named Rules
**The Escalation-On-Interaction Rule.** Interactive elements (cards, primary CTAs) start at `shadow-mini` or `shadow-card` at rest and escalate to `shadow-popover` plus a slight scale (`hover:scale-[1.02]`) on hover/interaction — never the reverse. Depth is earned by user attention, not present by default.

## 5. Components

### Buttons
- **Shape:** Gently rounded (5px default via `--radius-button`; `sm`/`lg` variants at 9px/10px).
- **Primary:** Ink background (`#18181b`), white text, `px-4 py-2` (default) up to `px-8 py-3.5` on hero/nav CTAs. Hero/nav CTAs additionally use `rounded-xl` (larger, ~12px) rather than the default button radius — a deliberate exception for the storefront's boldest CTAs.
- **Hover / Focus:** background dims to 90% opacity on hover; focus state is a 2px ring in the button's own color, offset 2px from the button edge — never a generic blue browser outline.
- **Outline / Secondary / Ghost:** Outline uses a bordered paper background; Secondary uses the Mist fill; Ghost is transparent until hover, filling with the Sky Pale accent.

### Cards
- **Corner Style:** 16px default (`rounded-card`), 10px small, 20px large.
- **Background:** Paper, with Ink-tinted border at 10% opacity (`--border-card`).
- **Shadow Strategy:** `shadow-card` at rest, `shadow-popover` for the "elevated" variant; see Elevation.
- **Interactive cards:** add `hover:scale-[1.02]` / `active:scale-[0.98]` — tactile, not just color feedback.
- **Internal Padding:** 24px default, 16px small, 32px large.

### Inputs / Fields
- **Style:** Floating-label pattern — label sits centered until focus/value, then animates up and shrinks. Rounded 9px, bordered, height ~3rem (`--spacing-input`) to fit the floating label comfortably.
- **Focus:** 2px ring in Ink/foreground color, 2px offset — matches button focus treatment for consistency.
- **Error:** border and helper text switch to Rose Alert; ring switches to match.

### Navigation
- **Style:** Fixed top bar, transparent over the hero video and transitioning to a blurred white bar (`bg-white/80 backdrop-blur-md`) on scroll or on non-hero pages. Desktop links use muted-foreground at rest, full foreground + semibold when active; mobile collapses into a full-height side drawer with 44px-minimum tap targets.
- **Active state:** bold weight + full-contrast color change, no underline or pill background — clarity over decoration.

### Trust Badge (signature component)
A pill-shaped, translucent badge (`bg-white/10`, `border-white/20`, `backdrop-blur-sm`, `rounded-full`) with a pulsing status dot (emerald `animate-ping`) — used in the hero to signal live availability ("Disponible maintenant"). This is the one place the system currently reaches for something more alive/energetic than the rest of the palette; a natural anchor point for extending "Studio Drop" energy elsewhere.

## 6. Do's and Don'ts

### Do:
- **Do** keep Ink + Paper as the load-bearing pair; let Sky/Rose/Green/Amber stay strictly functional.
- **Do** use Cal Sans sparingly, at heading scale only — its rarity is the brand signal.
- **Do** escalate shadow + scale together on hover/interactive states (`shadow-card` → `shadow-popover`, subtle scale) rather than shadow alone.
- **Do** treat the staff/admin/client dashboard as a plain, traditional tool — resist the urge to "brand" it; that effort belongs on the public storefront.
- **Do** keep the booking flow (service → person → slot → confirm) the least-decorated, fastest path in the whole product — per PRODUCT.md's positioning, decoration must never slow down the primary CTA.

### Don't:
- **Don't** introduce a warm cream/sand/beige base — this system is deliberately true-neutral (zinc-derived), not warmth-tinted.
- **Don't** let the public storefront settle for "calm and welcoming" as the finished brand — PRODUCT.md now calls for professional-but-energetic, streetwear-adjacent confidence; soft-and-tactile is the current floor, not the target ceiling.
- **Don't** carry new "branding effort" into the staff/admin/client dashboard — it should stay a traditional, utilitarian dashboard.
- **Don't** use a generic clinical/healthcare-portal look or a generic corporate-SaaS-marketing look, per PRODUCT.md's anti-references.
- **Don't** use `border-left`/`border-right` accent stripes, gradient text, or decorative glassmorphism outside the one deliberate trust-badge use.
