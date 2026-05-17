# ADR-0004: Session-Based Auth with Redis and Environment-Scoped Cookies

**Status**: Accepted  
**Date**: 2025-07

---

## Context

The platform needs authenticated sessions for web users across three environments (dev, staging, production) with distinct subdomains (admin, staff, public). Cookie names and domains must not bleed between environments, and the development loop must not require a real auth server.

## Decision

**Session storage**: Server-side sessions in Redis. The cookie holds only a session ID; all session state lives server-side. Redis provides sub-millisecond lookup.

**Cookie strategy**: Cookie name and domain are differentiated per environment:
- Development: mock user injected via `USE_MOCK_DATA=true`, no real cookie validation
- Staging: distinct cookie name from production to prevent session hijacking across environments
- Production: `Secure; HttpOnly; SameSite=Lax` with scoped domain

**Frontend validation**: `hooks.server.ts` validates the session cookie on every server request by calling `/users/me`. On failure, it clears the cookie and redirects to login. This is the single point of auth enforcement for all routes.

**Subdomain routing**: Admin and staff users are detected by hostname (`admin.*`, `staff.*`). Role-based redirects enforce that only users with the matching role can access those subdomains.

**OTP**: One-time passwords are cached in Redis with a short TTL. Redis is authoritative for OTP validity; the database is not involved.

## Consequences

**Positive:**
- Sessions can be invalidated server-side instantly (Redis delete) without waiting for token expiry
- Mock mode allows full frontend development without a running auth service
- Environment isolation prevents session bleed between staging and production

**Negative:**
- Redis is a hard runtime dependency — session lookups fail if Redis is unavailable
- Every authenticated request hits Redis (acceptable latency, mitigated by Redis being in-process network)
- Cookie domain configuration must be kept consistent between the `Set-Cookie` and deletion paths (past source of bugs — see git history)
