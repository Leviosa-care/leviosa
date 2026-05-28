# ADR 0009 — Inline guest account creation on the booking confirmation page

**Status:** Accepted  
**Date:** 2026-05-28

After a guest booking is confirmed, the confirmation page offers an inline, two-phase account creation form pre-filled with the guest's name and email (or phone — see ADR 0008). The guest enters only a password, submits, receives an OTP inline (no page navigation), and on verification their account is created and their booking is silently claimed. The resulting account is partial — gender, birthdate, and address are not collected at this step and are prompted later via an in-app nudge.

## Why not redirect to the standard `/auth` registration flow

The multi-step registration flow (email OTP → general info → address → password) exists for new users who arrive cold. A guest who just booked has already given their name and contact — routing them through that flow would ask for information we already have and force them to navigate away from the confirmation context. The inline form honours the implied promise: "you're almost there, one more step."

## Why OTP is still required

Skipping OTP for the inline form was considered. A guest booking that uses an unverified email is the guest's problem (no confirmation, no invoice). An account created with an unverified email is our problem — broken password resets, booking claims linked to the wrong person, support burden. The trust asymmetry justifies the gate. The OTP is handled inline so the user never leaves the page.

## Consequences

- The `authuser` service must expose an account-creation endpoint that accepts partial profile data (name, email, phone, password) without requiring gender, birthdate, or address.
- The booking service must expose a claim endpoint triggered immediately after account creation in this flow.
- A profile-completion prompt must be shown in the `(client)/` layout for accounts created via this path.
- Phone-only guests (no email on the booking) are shown the same inline form with an additional email field — email remains required for account identity.
