# Leviosa

Leviosa is booking and activity management software for service businesses (e.g. massage therapy practices) — room scheduling, availability, capacity tracking, payments, and client messaging.

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

## Deployment

Deployment to the VPS is driven by the root `Makefile` (Docker image build/push, Ansible playbooks, Terraform-managed infrastructure) — run `make help` to see available targets. GitHub Actions (`.github/workflows/ci.yaml`) handles pull request validation (build, unit tests, integration tests, vulnerability scanning) but does not perform deployments.
