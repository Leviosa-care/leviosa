# ADR-0013 — Partner-initiated messaging threads

**Status:** Accepted  
**Date:** 2026-06-02

## Context

The messaging service supports private threads between two users. Three initiation models were considered:

1. **Partner-initiated only** — only partners (and admins) can open a thread; clients can read and reply.
2. **Mutual initiation** — any authenticated user can open a thread with any other user.
3. **Client-initiated with booking check** — clients can open threads, but only with partners they have a booking with.

A booking relationship is a natural prerequisite for meaningful communication. Without it, clients could contact any partner on the platform, creating unsolicited messages and a moderation burden.

## Decision

Threads are partner-initiated. A partner may only open a thread with a client with whom they share at least one booking. Administrators may open threads unconditionally. Clients may read and reply to threads but may not create them.

This is enforced at the application layer (`ErrCannotInitiateThread` for Standard role) and at the route layer (`POST /threads` requires Partner role minimum).

## Consequences

- No backend changes are needed to support the client messages UI — clients already have `RequireStandard` access to read messages, send replies, mark threads as read, and receive SSE updates. Only `GET /threads` needs to be lowered from `RequirePartner` to `RequireStandard`.
- Clients who have never had a booking will see an empty messages page. This is correct — they have no threads yet.
- If a future requirement allows clients to initiate contact (e.g., pre-booking enquiries), the booking-relationship check would need to be inverted or replaced with a different gate.
