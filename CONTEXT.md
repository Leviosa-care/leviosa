# Leviosa — Domain Context

Leviosa is a **wellness booking platform** connecting clients with verified wellness professionals (therapists, coaches, practitioners). Clients browse services, book sessions in treatment rooms, and pay through the platform. Partners manage their availability and receive bookings. Administrators oversee the full operation.

---

## Ubiquitous Language

These are the canonical terms used in code, documentation, and conversation. Use them precisely.

### People & Roles

| Term | Definition |
|------|------------|
| **User** | Any registered account on the platform. Has a role that governs access. |
| **Client** | A user who books wellness sessions (role: `standard` or `premium`). |
| **Premium client** | A client with elevated access, able to book directly without extra steps. |
| **Partner** | A wellness professional who offers services through the platform (role: `partner`). Has a public profile with bio, certifications, and specializations. |
| **Staff** | Internal operational users managing rooms, allocations, and day-to-day scheduling (role: `staff`). |
| **Administrator** | Full platform access including user management, settings, and catalog (role: `administrator`). |

### Services & Pricing

| Term | Definition |
|------|------------|
| **Category** | A top-level grouping of services (e.g., Massage, Therapy, Coaching). |
| **Product** | A specific wellness service offered by a partner. Has a duration, buffer time, cancellation window, and images. Duration is always a multiple of 10 minutes, between 20 and 240 minutes. |
| **Buffer time** | Mandatory rest/preparation time after a session ends before the next booking can start. Defined on the product. |
| **Cancellation window** | Number of hours before a session within which cancellation is no longer permitted. Defined on the product. |
| **Price** | A Stripe-backed pricing record attached to a product. Can be one-time or recurring (monthly/yearly intervals). |

### Scheduling & Rooms

| Term | Definition |
|------|------------|
| **Building** | A physical location that contains treatment rooms. |
| **Room** | A treatment space within a building. Has a name, capacity, and allocation type. |
| **Allocation** | The assignment of a room to a partner. Two types: `dedicated` (fixed time period, one partner) or `shared` (ongoing availability pool). |
| **Availability** | A time window a partner offers for bookings. Can be a single occurrence or a recurring pattern (daily/weekly/monthly) with an optional end date. |
| **Booking** | A confirmed reservation by a client for a partner's service in a room. The core transactional entity. Tracks payment and progresses through a lifecycle. In client-facing UI the word "consultation" is used as a synonym, but the domain term and API surface always use "booking". |
| **Booking status** | The FSM governing a booking's lifecycle: `confirmed → completed | cancelled | no_show`. |
| **Earnings** | A read-only financial summary for a partner, derived by aggregating their completed bookings. Includes current-month revenue, last-month revenue, pending amounts, and a per-transaction history. Not a stored entity — computed on demand from the Booking table. |
| **Gap** | Unused time between bookings in a room schedule. The system detects gaps to surface upselling opportunities. |
| **Utilization** | A metric measuring how efficiently a room's time is used (utilization % minus fragmentation penalty). Computed via a materialized view. |
| **10-minute alignment** | All time slots are aligned to 10-minute boundaries for scheduling consistency. |

### Infrastructure Concepts

| Term | Definition |
|------|------------|
| **Session** | An authenticated user's server-side state, stored in Redis and referenced by a secure cookie. |
| **OTP** | One-time password used for email/phone verification steps. Cached in Redis with a TTL. |
| **Pepper** | A per-service encryption secret stored in HashiCorp Vault. Used by `encx` to encrypt sensitive fields at rest. |
| **encx** | The internal encryption library. Fields tagged `encx:"encrypt"` are transparently encrypted before persistence. |
| **Contract** | A typed message struct exchanged over RabbitMQ between services. Defined in `internal/common/messaging/contracts/`. |

---

## Domain Model

```
Category ──< Product ──< Price (Stripe)
                │
                └── Partner (User) ──< Availability
                                          │
Building ──< Room ──< Allocation ─────────┘
                          │
                          └──> Booking <── Client (User)
                                  │
                              Payment (Stripe)
```

### Key invariants

- A **Product** duration must be in `[20, 240]` minutes and a multiple of 10.
- A **Booking** can only be cancelled outside the product's cancellation window.
- A **Room** can have at most one `dedicated` allocation active at any time.
- **Availability** slots and booking start times must align to 10-minute boundaries.
- All PII fields (client/partner notes, building addresses, etc.) are encrypted via `encx` before hitting the database.

---

## Service Boundaries

Leviosa is a **modular monolith**: one deployed binary, five internal services with hard boundaries. Services communicate synchronously via function calls through port interfaces, and asynchronously via RabbitMQ for notifications and side effects.

| Service | Responsibility |
|---------|---------------|
| **authuser** | User registration, OTP verification, session management, partner profiles |
| **catalog** | Categories, products, images (S3), Stripe pricing |
| **booking** | Rooms, buildings, allocations, availabilities, bookings, metrics |
| **settings** | Platform-wide configuration key-value store |
| **notification** | Email and SMS delivery (consumed from RabbitMQ events) |

---

## Environments

| Environment | Purpose | Auth |
|-------------|---------|------|
| `development` | Local dev, mock user, mock data | `USE_MOCK_DATA=true` |
| `staging` | Integration testing, real services, password-protected | Real cookies, staging Vault |
| `production` | Public users | Real cookies, production Vault |

Cookie names and Vault paths differ per environment to prevent cross-contamination.
