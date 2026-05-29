# ADR-0011: OTP Delivery via Direct In-Process Call

**Status**: Accepted  
**Date**: 2026-05

---

## Context

The authuser service needs to deliver OTP codes by email during registration and password reset. RabbitMQ infrastructure already exists in the codebase and OTP events are already published to `otp.notification.email.queue` — but no consumer was ever implemented, meaning OTP emails were never actually sent.

Two delivery paths were considered when wiring up the notification service:

1. **Async via RabbitMQ**: authuser publishes the OTP event, a notification consumer reads the queue and sends the email. Keeps authuser decoupled from email delivery.
2. **Synchronous in-process call**: authuser calls `notification.SendOTPEmail()` directly through the port interface.

## Decision

OTP email is delivered via a **direct in-process call** from authuser to the notification service. The existing RabbitMQ publish of OTP events is removed.

## Rationale

OTP delivery is **not fire-and-forget**. If the email send fails, the user is stuck at a registration or login step with no visible error — the async path hides this failure behind a queue. With a synchronous call, a delivery failure returns an error immediately and the caller can surface it to the user.

This is consistent with how booking notifications work (`BookingNotificationAdapter` is also a direct in-process call). The modular monolith architecture (ADR-0001) explicitly favours in-process calls for synchronous flows.

## Consequences

**Positive:**
- Email delivery failures are immediately visible and surfaced to the caller
- Simpler: no queue consumer to maintain, no dead-letter queue to monitor
- Consistent with booking notification pattern

**Negative:**
- A slow SMTP connection adds latency to the OTP request response
- authuser is now coupled to the notification service port at compile time (acceptable within the monolith boundary)

**Future path**: If SMTP latency becomes a problem, the notification port can be made async-by-default internally (e.g., goroutine with timeout), without changing the caller interface.
