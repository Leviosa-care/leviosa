# Leviosa Ansible Setup Guide

Quick guide to provision your VPS after Terraform deployment.

---

## Prerequisites

Install Ansible and required collections locally:

```bash
pip install ansible
ansible-galaxy collection install community.docker community.general
```

Ensure you have your SSH key at `~/.ssh/leviosa.pub`.

---

## Quick Start (Complete Setup)

### 1. Get Terraform Outputs

```bash
cd infra/terraform
terraform output server_ipv4_address
terraform output domain_name
```

### 2. Configure Ansible Vault

```bash
cd infra/ansible
cp group_vars/secrets.vault.example.yml group_vars/secrets.yml
nano group_vars/secrets.yml
```

Fill in the required values:
```yaml
# Generate strong passwords
vault_db_password: "your_secure_postgres_password"
vault_session_secret: "your_random_32_character_string"

# From Terraform outputs
vault_aws_access_key_id: "<from terraform>"
vault_aws_secret_access_key: "<from terraform>"
vault_aws_region: "eu-central-1"
vault_s3_bucket: "leviosa-media"

# Email (Gmail SMTP or AWS SES)
vault_smtp_host: "smtp.gmail.com"
vault_smtp_port: "587"
vault_smtp_secure: "false"
vault_smtp_user: "your@gmail.com"
vault_smtp_password: "your_app_password"
vault_smtp_from_email: "noreply@leviosa.care"
vault_smtp_from_name: "Leviosa"

# Optional: OAuth, Backup AWS credentials
```

Encrypt the vault:
```bash
ansible-vault encrypt group_vars/secrets.yml
```

### 3. Run Initial Setup

```bash
make setup
```

This prompts for confirmation, then runs the complete site playbook.

### 4. Deploy Application

```bash
make deploy
```

### 5. Configure Backups (Optional)

```bash
make backup
```

---

## What `make setup` Does

| Task | Description |
|------|-------------|
| **System** | Base hardening, automatic updates, kernel parameters |
| **SSH** | Key-only auth, root login disabled, fail2ban |
| **Firewall** | UFW with ports 22, 80, 443 only |
| **User** | Creates `leviosa` user with restricted sudo |
| **Docker** | Installs Docker & Docker Compose |
| **Caddy** | Reverse proxy with automatic TLS, security headers |

---

## Manual Commands (Without Makefile)

### Initial Setup

```bash
ansible-playbook -i inventory/hosts playbooks/site.yml \
  --ask-vault-pass \
  -e "ansible_host=46.225.25.238" \
  -e "deploy_ssh_public_key=$(cat ~/.ssh/leviosa.pub)" \
  -e "app_domain=leviosa.care"
```

### Deploy Application

```bash
ansible-playbook -i inventory/hosts playbooks/deploy.yml \
  --ask-vault-pass \
  -e "ansible_host=46.225.25.238"
```

### Configure Backups

```bash
ansible-playbook -i inventory/hosts playbooks/backup.yml \
  --ask-vault-pass \
  -e "ansible_host=46.225.25.238" \
  -e "backup_s3_bucket=leviosa-backups" \
  -e "backup_s3_region=eu-central-1"
```

---

## Useful Makefile Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make setup` | Complete VPS setup (first time) |
| `make deploy` | Deploy/update production |
| `make deploy-staging` | Deploy/update staging |
| `make restart` | Quick restart (pull + restart app) |
| `make restart-staging` | Quick restart staging |
| `make backup` | Configure database backups |
| `make ssh-info` | Display SSH connection details |
| `make status` | Check production service status |
| `make staging-status` | Check staging service status |
| `make logs` | Stream production logs |
| `make staging-logs` | Stream staging logs |
| `make health-check` | Run health check |
| `make check-clean` | Check playbooks for syntax errors |
| `make install-deps` | Install Ansible dependencies |

---

## SSH Access

After setup, connect as the `leviosa` user:

```bash
# View connection info
make ssh-info

# Connect
ssh leviosa@46.225.25.238
```

Root login is **disabled** for security. Use `leviosa` with sudo.

---

## Troubleshooting

### View Vault Contents

```bash
ansible-vault view group_vars/secrets.yml
```

### Edit Vault

```bash
ansible-vault edit group_vars/secrets.yml
```

### Run with Verbose Output

```bash
ansible-playbook playbooks/site.yml -vvv
```

### Check Service Status

```bash
ssh leviosa@46.225.25.238
docker compose ps
```

### View Logs

```bash
ssh leviosa@46.225.25.238
docker compose logs -f
```

---

## Security Checklist

After setup, verify:

- [ ] SSH password authentication disabled
- [ ] Root login disabled
- [ ] UFW firewall enabled
- [ ] Fail2ban running
- [ ] Deploy user has sudo access
- [ ] Environment files have correct permissions (0600)

---

## Terraform Outputs Reference

| Output | Example | Usage |
|--------|---------|-------|
| `server_ipv4_address` | `46.225.25.238` | Ansible `ansible_host` |
| `domain_name` | `leviosa.care` | Ansible `app_domain` |
| `backup_bucket_name` | `leviosa-backups` | Backup playbook |

---

## File Structure

```
infra/ansible/
├── Makefile                      # Quick commands
├── README.md                     # Detailed documentation
├── SETUP_GUIDE.md                # This file
├── VAULT.md                      # Vault usage guide
├── ansible.cfg                   # Ansible config
├── inventory/
│   └── hosts                     # Server inventory
├── group_vars/
│   ├── secrets.yml               # Your encrypted secrets
│   └── secrets.vault.example.yml # Vault template
├── playbooks/
│   ├── site.yml                  # Complete setup
│   ├── deploy.yml                # Application deployment
│   ├── deploy-staging.yml        # Staging deployment
│   └── backup.yml                # Backup configuration
└── roles/
    ├── system/                   # Base hardening
    ├── ssh/                      # SSH configuration
    ├── ufw/                      # Firewall setup
    ├── docker/                   # Docker installation
    ├── user/                     # Deploy user
    ├── app/                      # Application deployment
    ├── caddy/                    # Reverse proxy
    ├── fail2ban/                 # Brute-force protection
    └── monitoring/               # Metrics collection
```
