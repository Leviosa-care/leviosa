# Leviosa Ansible Configuration

Automated VPS provisioning and application deployment using Ansible.

The VPS is shared by two app environments — production and staging — both deployed from the same Terraform-provisioned server (see `infra/terraform/README.md`). Staging shares production's Postgres, Redis, RabbitMQ, and Caddy instance to save resources; only the app container is separate.

## What This Does

### Security Hardening
- SSH key-only authentication (password auth disabled)
- Root login disabled (deploy user with restricted sudo only)
- Fail2ban for SSH/Caddy brute-force protection
- UFW firewall with only essential ports open
- Automatic security updates (unattended-upgrades)
- Kernel hardening (sysctl, secure shared memory)
- Modern cryptographic algorithms only (Curve25519, ChaCha20)
- Restricted sudo — deploy user can only run specific commands
- Caddy reverse proxy — application bound to localhost, Caddy handles TLS
- Docker metrics bound to localhost only (127.0.0.1:9323)
- HashiCorp Vault deployed alongside the app for secrets (transit encryption keys, KV secrets used by `encx`)

### System Setup
- Docker & Docker Compose installation
- Non-root deploy user with restricted sudo
- Caddy reverse proxy with SSL/TLS and security headers
- Application directory structure
- GPG configuration (for `rclone`-based backup encryption)

### SSL/TLS Configuration
- Caddy reverse proxy with automatic TLS (Let's Encrypt)
- Security headers (HSTS, X-Frame-Options, CSP)
- OCSP stapling enabled

### Application Deployment
- Docker Compose configuration
- Environment file management
- Vault initialization/unseal and admin-user seeding (staging deploy only — see `playbooks/deploy-staging.yml`)
- Health check scripts

### Database Backups
- `rclone`-based backups to S3 with GPG encryption
- AWS CLI/credentials configuration on the server

### Monitoring
- cAdvisor for container metrics (127.0.0.1, port configurable via `monitoring_cadvisor_port`)
- node-exporter for host metrics (127.0.0.1:9100)
- Disabled by default (`monitoring_enabled: false` in `group_vars/all.yml`)

## Prerequisites

### Local Machine

```bash
make install-deps
# or manually:
pip install ansible
ansible-galaxy collection install community.docker community.general
```

### VPS Access

You need:
- Root SSH access to your VPS
- Your public SSH key (`~/.ssh/leviosa.pub`, falls back to `~/.ssh/id_rsa.pub`)

### Terraform Outputs

The Makefile pulls the VPS IP and domain from Terraform automatically:

```bash
cd ../terraform
make apply-staging      # provisions the shared VPS, if not already up
```

### Secrets

Copy the example secrets files and fill in real values. These are gitignored — only the `.example` files are committed:

```bash
cd group_vars
cp leviosa_staging.example.yml leviosa_staging.yml
cp leviosa_production.example.yml leviosa_production.yml
nano leviosa_staging.yml
nano leviosa_production.yml
```

`group_vars/all.yml` holds shared, non-secret defaults (deploy user, Docker daemon options, fail2ban thresholds, etc.) and is committed to git.

## Quick Start

### 1. Initial VPS Setup (run once per environment)

```bash
make setup-staging
make setup-production
```

These targets pull the VPS IP/domain from Terraform, read your SSH public key, and run `playbooks/site.yml` with the right `group_vars/all.yml` + environment-specific secrets file. You'll be prompted to confirm before anything runs.

### 2. Deploy the Application

```bash
make deploy-staging
make deploy-production
```

These run `playbooks/deploy-staging.yml` / `playbooks/deploy.yml`, which deploy/update Caddy + the app container without touching system configuration. The staging deploy additionally creates the staging Postgres database/RabbitMQ vhost on the shared services, and handles Vault init/unseal and admin-user seeding.

### 3. Quick Restart (pull latest image, restart only)

```bash
make restart                    # production: all services
make restart-frontend           # production: frontend only
make restart-backend            # production: backend only
make restart-staging            # staging: all services
make restart-staging-frontend
make restart-staging-backend
```

### 4. Database Backups

```bash
make backup
```

Backups are not yet wired up to a Terraform-provisioned S3 bucket — `make backup` currently exits with an error reminding you to add one first. The `rclone` and `gpg` roles and `playbooks/backup.yml` exist and are ready to use once a backup bucket is configured.

## Playbooks

### `site.yml` — Complete Setup

Run on a fresh VPS to set up everything. Roles run in this order, each with its tags:

| Role | Tags | Purpose |
|------|------|---------|
| `system` | `system`, `setup` | Base system hardening (sysctl, shared memory, etc.) |
| `docker` | `docker`, `setup` | Docker & Docker Compose installation |
| `common` | `common`, `setup` | Shared baseline config |
| `ssh` | `ssh`, `security` | SSH hardening |
| `fail2ban` | `fail2ban`, `security` | Brute-force protection |
| `ufw` | `ufw`, `security` | Firewall setup |
| `caddy` | `caddy`, `web` | Reverse proxy |
| `gpg` | `gpg`, `security` | GPG config for backup encryption |
| `rclone` | `rclone`, `backup` | S3 backup tooling |
| `app` | `app`, `deploy` | Application deployment |
| `monitoring` | `monitoring`, `setup` | Metrics collection (only if `monitoring_enabled: true`) |

```bash
ansible-playbook -i inventory/hosts.yml playbooks/site.yml \
  -e "ansible_host=<server-ip>" \
  -e "deploy_ssh_public_key=$(cat ~/.ssh/leviosa.pub)" \
  -e "app_domain=yourdomain.com"
```

### `deploy.yml` / `deploy-staging.yml` — Application Deployment

Deploy or update the application only (Caddy + app container), targeting `leviosa_servers` or `leviosa_staging_servers` respectively. Run after the initial `site.yml` setup.

### `backup.yml` — Database Backups

```bash
ansible-playbook -i inventory/hosts.yml playbooks/backup.yml \
  -e "backup_s3_bucket=my-backups" \
  -e "backup_s3_region=eu-central-1"
```

## Variables

### Inventory (`inventory/hosts.yml`)

Two host groups share the same physical VPS IP:

| Group | Host | Purpose |
|-------|------|---------|
| `leviosa_servers` | `leviosa_vps` | Production app, owns Postgres/Redis/RabbitMQ/Caddy |
| `leviosa_staging_servers` | `leviosa_vps_staging` | Staging app, shares production's services; `caddy_enabled: false` |

Key per-group vars: `app_domain`, `app_port`, `app_docker_image`, `app_base_dir`, `ssh_port`, `caddy_enabled`.

### Shared Group Variables (`group_vars/all.yml`, committed)

| Variable | Description | Default |
|----------|-------------|---------|
| `deploy_user` | Deploy user name | `leviosa` |
| `deploy_sudo_restricted` | Restrict sudo to specific commands | `true` |
| `auto_security_updates` | Enable unattended-upgrades | `true` |
| `fail2ban_enabled` | Enable fail2ban | `true` |
| `backup_enabled` | Enable backup role behavior | `true` |
| `monitoring_enabled` | Enable cAdvisor/node-exporter | `false` |
| `staging_enabled` | Include staging site blocks in Caddyfile | `false` |

### Per-Environment Secrets (`group_vars/leviosa_staging.yml`, `leviosa_production.yml`, gitignored)

Database/Redis/session credentials, AWS keys, Stripe keys, SMTP config, Cloudflare API token for Caddy's DNS-01 challenge, and `seed_admins` (admin users upserted on every deploy). See `leviosa_production.example.yml` for the full list and structure.

## Directory Structure

```
infra/ansible/
├── ansible.cfg                     # Ansible configuration
├── Makefile                         # Workspace-aware setup/deploy/restart commands
├── README.md                        # This file
├── SETUP_GUIDE.md                   # Quick start guide
├── VAULT.md                         # HashiCorp Vault operations guide (init/unseal, troubleshooting)
├── VAULT_PRODUCTION_SETUP.md        # HashiCorp Vault production-mode reference
├── inventory/
│   └── hosts.yml                    # leviosa_servers + leviosa_staging_servers
├── group_vars/
│   ├── all.yml                      # Shared, non-secret defaults (committed)
│   ├── leviosa_staging.yml          # Staging secrets (gitignored)
│   ├── leviosa_staging.example.yml
│   ├── leviosa_production.yml       # Production secrets (gitignored)
│   └── leviosa_production.example.yml
├── playbooks/
│   ├── site.yml                     # Complete setup
│   ├── deploy.yml                   # Production app deployment
│   ├── deploy-staging.yml           # Staging app deployment (Vault init/unseal, seeding)
│   ├── backup.yml                   # Backup configuration
│   └── templates/                   # Backup/AWS-credentials templates
└── roles/
    ├── system/                      # Base system hardening
    ├── docker/                      # Docker installation
    ├── common/                      # Shared baseline config
    ├── ssh/                         # SSH configuration
    ├── fail2ban/                    # Brute-force protection
    ├── ufw/                         # Firewall setup
    ├── caddy/                       # Reverse proxy
    ├── gpg/                         # Backup encryption keys
    ├── rclone/                      # S3 backup storage
    ├── app/                         # Application deployment
    ├── monitoring/                  # Metrics collection
    └── user/                        # Deploy user setup (currently unused by any playbook)
```

## Security Features

### SSH Hardening
- Password authentication disabled (key-only)
- Root login disabled (use deploy user with sudo)
- Modern cryptographic algorithms (Curve25519, ChaCha20-Poly1305)
- Disabled X11 forwarding, agent forwarding, TCP forwarding
- Fail2ban: SSH brute-force protection

### Restricted Sudo
- Deploy user can only run specific commands without a password (`docker`, `docker-compose`, `systemctl` for docker, app scripts)
- See `deploy_sudo_restricted` / `deploy_sudo_password_fallback` in `group_vars/all.yml`

### Firewall (UFW)
- Default deny incoming
- Allow SSH, HTTP (80), HTTPS (443) only

### Caddy Reverse Proxy
- Application bound to localhost (127.0.0.1) only
- Caddy handles external connections with automatic TLS (Let's Encrypt, DNS-01 via Cloudflare)
- Security headers: HSTS, X-Frame-Options, CSP

### Docker Security
- Metrics endpoint bound to localhost only (127.0.0.1:9323)

### Vault
- HashiCorp Vault runs alongside the app, initialized/unsealed automatically on staging deploy
- Used for transit encryption keys and KV secrets (`encx` pepper material)

## Maintenance

```bash
make status              # production service status
make staging-status      # staging service status
make logs                # stream production logs
make staging-logs        # stream staging logs
make health-check         # run health check script on the server
make ssh-info             # show SSH connection details
make check-clean          # syntax-check all playbooks
make lint                 # ansible-lint (requires installation)
make vars-test            # display effective variables
```

## Troubleshooting

### Connection Refused
```bash
ping <server-ip>
ssh root@<server-ip>
```

### Permission Denied
```bash
cat ~/.ssh/leviosa.pub
```

### Playbook Fails
```bash
ansible-playbook playbooks/site.yml -vvv
```

### Service Not Starting
```bash
ssh deploy@<server-ip>
docker compose logs
```

## Security Checklist

After initial setup, verify:

- [ ] SSH password authentication disabled
- [ ] Root login disabled
- [ ] UFW firewall enabled
- [ ] Fail2ban running
- [ ] Deploy user has restricted sudo access
- [ ] Environment files have correct permissions (0600)
- [ ] Vault initialized and unsealed
- [ ] Monitoring endpoints accessible (localhost only, if enabled)
