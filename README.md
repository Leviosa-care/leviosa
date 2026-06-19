# Leviosa

Leviosa is booking and activity management software for service businesses — massage therapy practices, studios, and similar appointment-based businesses that need to manage rooms, staff schedules, and clients in one place.

## Features

- **Room scheduling & availability** — booking calendar, gap detection between bookings to surface unused room time, room utilization metrics (efficiency score, fragmentation, idle time), and suggested availability block durations based on the partner's product catalog
- **Product catalog & pricing** — categories and products (treatments, classes, events) with capacity tracking and Stripe-backed pricing/checkout
- **Payments** — Stripe integration for checkout and customer billing
- **Client messaging** — direct conversations between businesses and their clients
- **Notifications** — booking confirmations and reminders over email (Gmail SMTP) and SMS (Twilio)
- **Auth & accounts** — email/OTP sign-in, sessions, and Google/Apple OAuth
- **Settings** — per-business configuration (OTP policy, notification preferences, etc.)

## Architecture

```
├── frontend/           # SvelteKit 5 application
├── backend/            # Go modular monolith
├── config/             # Infrastructure configuration (Caddy, Loki, Vault, etc.)
├── infra/              # Terraform and Ansible deployment scripts
└── compose.yaml        # Production docker-compose configuration
```

- **Frontend**: [SvelteKit 5](https://kit.svelte.dev/) + TypeScript + Tailwind CSS v4.
- **Backend**: [Go](https://golang.org/) modular monolith — a single binary built from `backend/cmd/app`, with business domains (`authuser`, `booking`, `catalog`, `settings`, `notification`, `messaging`) organized as internal packages under `backend/internal/`, each following a hexagonal (ports & adapters) structure. See `backend/CLAUDE.md` for details.
- **Database**: PostgreSQL.
- **Cache**: Redis.
- **Message queue**: RabbitMQ.
- **File storage**: AWS S3.
- **Secrets**: HashiCorp Vault.
- **Payments**: Stripe.
- **Notifications**: Gmail SMTP, Twilio SMS.

## Local Development

```bash
make local-dev   # run frontend + backend in dev mode
make local-up    # start local docker-compose stack
make local-down  # stop local docker-compose stack
```

See `make help` at the repo root for the full list of commands (local dev, Docker image builds, VPS deploy, Terraform, Ansible). Frontend- and backend-specific commands are documented in `frontend/CLAUDE.md` and `backend/CLAUDE.md`.
