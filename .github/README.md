# GitHub Actions Workflows

This directory contains the CI workflow for the Leviosa project.

## Workflow

### CI (`ci.yaml`)

**Trigger:** Pull requests to `main`
**Purpose:** Validate code before merge — does not deploy anything

#### Jobs:
- **frontend-checks**: Install deps, build the SvelteKit app
- **backend-checks**: Build the Go binary, run unit tests (`go test ./...`)
- **backend-integration-tests**: Integration tests via `make test-integration` (dependencies spun up per-test with testcontainers)
- **security-scan**: Go vulnerability scanning with `govulncheck`
- **summary**: Aggregates the above; fails the check if any job failed

**Failure Behavior:** PR cannot merge until all jobs pass

**Merge Strategy:** Rebase and merge only (preserves full commit history)

## Deployment

There is no GitHub Actions deployment workflow. Staging and production deploys are manual, driven by the root `Makefile`:

- `make image-build` / `make image-push` / `make image-release` — build and push Docker images (and the `-staging` variants)
- `make deploy` / `make deploy-staging` — quick SSH-based pull + restart on the VPS
- `make ansible-deploy` / `make ansible-deploy-staging` — full deploy via Ansible
- `make infra-*` — Terraform-managed infrastructure

Run `make help` at the repo root for the full list.
