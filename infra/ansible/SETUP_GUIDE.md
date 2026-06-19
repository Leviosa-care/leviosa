# Leviosa Ansible Setup Guide

Quick guide to provision the VPS and deploy the app, after Terraform has stood up infrastructure (see `infra/terraform/README.md`). For full detail see `README.md`; for how HashiCorp Vault (the app's secrets backend) is handled, see `VAULT.md`.

---

## Prerequisites

Install Ansible and required collections locally:

```bash
make install-deps
# or manually:
pip install ansible
ansible-galaxy collection install community.docker community.general
```

Ensure you have your SSH key at `~/.ssh/leviosa.pub` (falls back to `~/.ssh/id_rsa.pub`).

---

## Quick Start (Complete Setup)

### 1. Provision the VPS with Terraform

```bash
cd ../terraform
make apply-staging       # provisions the shared VPS + DNS (staging/default workspace)
make apply-production     # provisions production-only AWS resources (S3/IAM)
```

### 2. Configure Secrets

Secrets are plain (gitignored) YAML files, **not** encrypted with `ansible-vault`:

```bash
cd ../ansible
cp group_vars/leviosa_staging.example.yml group_vars/leviosa_staging.yml
cp group_vars/leviosa_production.example.yml group_vars/leviosa_production.yml
nano group_vars/leviosa_staging.yml
nano group_vars/leviosa_production.yml
```

Fill in the required values (see the `.example` files for the full list):
```yaml
db_password: "your_secure_postgres_password"
redis_password: "your_secure_redis_password"
session_secret: "your_random_32_character_string_here_min_length"

aws_access_key_id: "<from terraform output>"
aws_secret_access_key: "<from terraform output>"
s3_bucket: "staging-leviosa-assets"

smtp_from_email: "noreply@leviosa.care"
smtp_from_name: "Leviosa"

stripe_secret_key: "sk_live_..."
stripe_webhook_secret: "whsec_..."
stripe_publishable_key: "pk_live_..."

caddy_cloudflare_api_token: "your_cloudflare_api_token_here"
```

To pull the AWS credential fields straight from the latest Terraform apply instead of copying by hand:
```bash
cd ../terraform
make update-staging-vault
make update-production-vault
```

`group_vars/all.yml` holds shared, non-secret defaults (deploy user, fail2ban thresholds, etc.) and is already committed — no setup needed there.

### 3. Run Initial Setup

```bash
make setup-staging
make setup-production
```

Each prompts for confirmation, pulls the VPS IP/domain from Terraform automatically, and runs `playbooks/site.yml` against the right host group with the right secrets file.

### 4. Deploy the Application

```bash
make deploy-staging
make deploy-production
```

This also handles HashiCorp Vault initialization/unsealing and admin-user seeding automatically — see `VAULT.md` if you need to inspect or recover it.

### 5. Configure Backups (Not Yet Available)

```bash
make backup
```

This currently exits with an error — there is no backup S3 bucket wired up via Terraform yet. The `rclone`/`gpg` roles and `playbooks/backup.yml` are ready to use once one is added.

---

## What `make setup-staging` / `make setup-production` Does

Runs `playbooks/site.yml`, which applies these roles in order:

| Role | Description |
|------|-------------|
| `system` | Base hardening, automatic updates, kernel parameters |
| `docker` | Installs Docker & Docker Compose |
| `common` | Shared baseline configuration |
| `ssh` | Key-only auth, root login disabled |
| `fail2ban` | Brute-force protection |
| `ufw` | Firewall with ports 22, 80, 443 only |
| `caddy` | Reverse proxy with automatic TLS, security headers |
| `gpg` | Backup encryption keys |
| `rclone` | S3 backup tooling |
| `app` | Application + Vault deployment |
| `monitoring` | cAdvisor/node-exporter (only if `monitoring_enabled: true`) |

---

## Manual Commands (Without Makefile)

### Initial Setup

```bash
ansible-playbook -i inventory/hosts.yml playbooks/site.yml \
  -e "ansible_host=<server-ip>" \
  -e "deploy_ssh_public_key=$(cat ~/.ssh/leviosa.pub)" \
  -e "app_domain=leviosa.care" \
  -e "@group_vars/all.yml" \
  -e "@group_vars/leviosa_staging.yml"
```

### Deploy Application

```bash
ansible-playbook -i inventory/hosts.yml playbooks/deploy.yml \
  -e "ansible_host=<server-ip>" \
  -e "@group_vars/all.yml" \
  -e "@group_vars/leviosa_production.yml"

ansible-playbook -i inventory/hosts.yml playbooks/deploy-staging.yml \
  -e "ansible_host=<server-ip>" \
  -e "@group_vars/all.yml" \
  -e "@group_vars/leviosa_production.yml" \
  -e "@group_vars/leviosa_staging.yml"
```

### Configure Backups

```bash
ansible-playbook -i inventory/hosts.yml playbooks/backup.yml \
  -e "ansible_host=<server-ip>" \
  -e "backup_s3_bucket=leviosa-backups" \
  -e "backup_s3_region=eu-central-1"
```

---

## Useful Makefile Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make setup-staging` / `make setup-production` | Complete VPS setup (first time per environment) |
| `make deploy-staging` / `make deploy-production` | Deploy/update the application |
| `make restart` / `make restart-frontend` / `make restart-backend` | Pull + restart production |
| `make restart-staging` / `make restart-staging-frontend` / `make restart-staging-backend` | Pull + restart staging |
| `make backup` | Configure database backups (currently a stub, see above) |
| `make ssh-info` | Display SSH connection details |
| `make status` / `make staging-status` | Check service status |
| `make logs` / `make staging-logs` | Stream logs |
| `make health-check` | Run health check on the server |
| `make check-clean` | Check playbooks for syntax errors |
| `make lint` | Run ansible-lint |
| `make vars-test` | Display effective variables |
| `make install-deps` | Install Ansible dependencies |
| `make terraform-outputs` | Show Terraform outputs |

---

## SSH Access

After setup, connect as the `leviosa` deploy user (see `deploy_user` in `group_vars/all.yml`):

```bash
make ssh-info     # view connection info
ssh -i ~/.ssh/leviosa leviosa@<server-ip>
```

Root login is **disabled** for ordinary operations; the initial setup connects as root once, then switches to the deploy user.

---

## Troubleshooting

### Run with Verbose Output

```bash
ansible-playbook playbooks/site.yml -vvv
```

### Check Service Status

```bash
ssh -i ~/.ssh/leviosa leviosa@<server-ip>
docker compose ps
```

### View Logs

```bash
ssh -i ~/.ssh/leviosa leviosa@<server-ip>
docker compose logs -f
```

### Vault Issues

See `VAULT.md` — covers checking status, manually unsealing, and where the unseal keys/root token live on the server.

---

## Security Checklist

After setup, verify:

- [ ] SSH password authentication disabled
- [ ] Root login disabled
- [ ] UFW firewall enabled
- [ ] Fail2ban running
- [ ] Deploy user has restricted sudo access
- [ ] Environment files have correct permissions (0600)
- [ ] HashiCorp Vault initialized and unsealed (`docker exec leviosa_vault vault status`)

---

## Terraform Outputs Reference

| Output | Example | Usage |
|--------|---------|-------|
| `server_ipv4_address` | `46.224.117.126` | Ansible `ansible_host` (pulled automatically by the Makefile) |
| `domain_name` | `leviosa.care` | Ansible `app_domain` |
| `media_bucket_name` | `staging-leviosa-assets` | Referenced by `s3_bucket` in the secrets file |
| `app_user_access_key_id` / `app_user_secret_access_key` | — | Synced into the secrets file by `make update-staging-vault` / `make update-production-vault` |

---

## File Structure

```
infra/ansible/
├── Makefile                          # Setup/deploy/restart commands
├── README.md                         # Detailed documentation
├── SETUP_GUIDE.md                    # This file
├── VAULT.md                          # HashiCorp Vault operations guide
├── VAULT_PRODUCTION_SETUP.md         # Vault production-mode reference
├── ansible.cfg                       # Ansible config
├── inventory/
│   └── hosts.yml                     # leviosa_servers + leviosa_staging_servers
├── group_vars/
│   ├── all.yml                       # Shared, non-secret defaults (committed)
│   ├── leviosa_staging.yml           # Staging secrets (gitignored)
│   ├── leviosa_staging.example.yml
│   ├── leviosa_production.yml        # Production secrets (gitignored)
│   └── leviosa_production.example.yml
├── playbooks/
│   ├── site.yml                      # Complete setup
│   ├── deploy.yml                    # Production application deployment
│   ├── deploy-staging.yml            # Staging application deployment
│   └── backup.yml                    # Backup configuration
└── roles/
    ├── system/                       # Base hardening
    ├── docker/                       # Docker installation
    ├── common/                       # Shared baseline config
    ├── ssh/                          # SSH configuration
    ├── fail2ban/                     # Brute-force protection
    ├── ufw/                          # Firewall setup
    ├── caddy/                        # Reverse proxy
    ├── gpg/                          # Backup encryption keys
    ├── rclone/                       # S3 backup storage
    ├── app/                          # Application + Vault deployment
    ├── monitoring/                   # Metrics collection
    └── user/                         # Deploy user setup (currently unused by any playbook)
```
