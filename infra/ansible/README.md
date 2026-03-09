# Leviosa Ansible Configuration

Automated VPS provisioning and application deployment using Ansible.

## What This Does

### Security Hardening
- SSH key-only authentication (password auth disabled)
- Root login disabled (deploy user with restricted sudo only)
- Fail2ban for SSH brute-force protection (5 attempts = 1 hour ban)
- UFW firewall with only essential ports open
- Automatic security updates (unattended-upgrades)
- Kernel hardening (sysctl, secure shared memory)
- Modern cryptographic algorithms only (Curve25519, ChaCha20)
- **Restricted sudo** - deploy user can only run specific commands
- **Caddy reverse proxy** - application bound to localhost, Caddy handles TLS
- **Docker metrics** - bound to localhost only (127.0.0.1:9323)
- **Ansible Vault** - encrypted secrets management

### System Setup
- Docker & Docker Compose installation
- Non-root deploy user with **restricted sudo**
- Caddy reverse proxy with SSL/TLS and security headers
- Application directory structure
- Logrotate configuration

### SSL/TLS Configuration
- Caddy reverse proxy with automatic TLS (Let's Encrypt)
- Security headers (HSTS, X-Frame-Options, CSP)
- OCSP stapling enabled

### Application Deployment
- Docker Compose configuration
- Environment file management
- Health check scripts
- Service management scripts

### Database Backups
- Automated PostgreSQL backups to S3
- Tiered retention (daily/weekly/monthly)
- S3 lifecycle integration
- AWS CLI configuration

### Monitoring
- cAdvisor for container metrics (127.0.0.1:9323)
- node-exporter for host metrics (127.0.0.1:9100)
- Systemd service for automatic startup

## Prerequisites

### Local Machine

```bash
# Install Ansible
sudo apt update
sudo apt install ansible -y

# Or with pip
pip install ansible

# Install required collections
ansible-galaxy collection install community.docker community.general
```

### VPS Access

You need:
- Root SSH access to your VPS
- Your public SSH key (`~/.ssh/leviosa.pub`)

### Terraform Outputs

First, run Terraform to get your VPS details:

```bash
cd infra/terraform
terraform output server_ipv4_address
terraform output domain_name
```

## Quick Start

### 1. Initial VPS Setup

```bash
cd infra/ansible

# Set your SSH public key
export DEPLOY_SSH_PUBLIC_KEY="$(cat ~/.ssh/leviosa.pub)"

# Set your domain
export APP_DOMAIN="yourdomain.com"

# Get VPS IP from Terraform
export ANSIBLE_HOST=$(cd ../terraform && terraform output -raw server_ipv4_address)

# Run full setup
ansible-playbook -i inventory/hosts playbooks/site.yml \
  -e "ansible_host=$ANSIBLE_HOST" \
  -e "deploy_ssh_public_key=$DEPLOY_SSH_PUBLIC_KEY" \
  -e "app_domain=$APP_DOMAIN"
```

### 2. Configure Environment Variables

After setup completes, SSH into the server and configure your `.env`:

```bash
# SSH as deploy user
ssh deploy@$(cd ../terraform && terraform output -raw server_ipv4_address)

# Edit environment file
cd /opt/leviosa
nano .env
```

### 3. Deploy Application

```bash
# Using Ansible
ansible-playbook playbooks/deploy.yml \
  -e "ansible_host=$ANSIBLE_HOST"

# Or manually on the server
ssh deploy@<server-ip>
cd /opt/leviosa
./prod-start.sh
```

### 4. Setup Database Backups

```bash
# Get backup bucket from Terraform
BACKUP_BUCKET=$(cd ../terraform && terraform output -raw backup_bucket_name)

# Run backup setup
ansible-playbook playbooks/backup.yml \
  -e "ansible_host=$ANSIBLE_HOST" \
  -e "backup_s3_bucket=$BACKUP_BUCKET"
```

## Playbooks

### `site.yml` - Complete Setup

Run this on a fresh VPS to set up everything:

```bash
ansible-playbook playbooks/site.yml \
  -e "ansible_host=<server-ip>" \
  -e "deploy_ssh_public_key=$(cat ~/.ssh/leviosa.pub)" \
  -e "app_domain=yourdomain.com"
```

Tags:
- `system` - Base system configuration
- `ssh` - SSH hardening
- `firewall` - UFW firewall setup
- `docker` - Docker installation
- `user` - Deploy user setup
- `app` - Application deployment
- `monitoring` - Metrics collection setup

### `deploy.yml` - Application Deployment

Deploy or update the application only:

```bash
ansible-playbook playbooks/deploy.yml \
  -e "ansible_host=<server-ip>"
```

### `deploy-staging.yml` - Staging Deployment

Deploy to staging environment:

```bash
ansible-playbook playbooks/deploy-staging.yml \
  -e "ansible_host=<server-ip>"
```

### `backup.yml` - Database Backups

Configure automated database backups:

```bash
ansible-playbook playbooks/backup.yml \
  -e "ansible_host=<server-ip>" \
  -e "backup_s3_bucket=my-backups" \
  -e "backup_s3_region=eu-central-1"
```

## Variables

### Inventory Variables (`inventory/hosts`)

| Variable | Description | Default |
|----------|-------------|---------|
| `ansible_host` | VPS IP address | Required |
| `app_domain` | Your domain name | `leviosa.care` |
| `app_port` | Application port | `3000` |
| `ssh_port` | SSH port | `22` |

### Group Variables (`group_vars/`)

| Variable | Description | Default |
|----------|-------------|---------|
| `deploy_user` | Deploy user name | `leviosa` |
| `deploy_ssh_public_key` | SSH public key | Required |
| `auto_security_updates` | Enable auto updates | `true` |
| `fail2ban_enabled` | Enable fail2ban | `true` |
| `backup_enabled` | Enable backups | `true` |

## Directory Structure

```
infra/ansible/
├── ansible.cfg                 # Ansible configuration
├── Makefile                     # Quick commands
├── README.md                    # This file
├── SETUP_GUIDE.md               # Quick start guide
├── VAULT.md                     # Vault usage guide
├── inventory/
│   └── hosts                    # Inventory file
├── group_vars/
│   ├── secrets.yml              # Encrypted secrets
│   └── secrets.vault.example.yml # Vault template
├── playbooks/
│   ├── site.yml                 # Complete setup
│   ├── deploy.yml               # Application deployment
│   ├── deploy-staging.yml       # Staging deployment
│   ├── backup.yml               # Backup configuration
│   └── templates/               # Backup templates
└── roles/
    ├── system/                  # Base system hardening
    ├── ssh/                     # SSH configuration
    ├── firewall/                # UFW setup (via ufw role)
    ├── docker/                  # Docker installation
    ├── user/                    # Deploy user setup
    ├── app/                     # Application deployment
    ├── caddy/                   # Reverse proxy
    ├── fail2ban/                # Brute-force protection
    ├── monitoring/              # Metrics collection
    └── rclone/                  # Backup storage
```

## Security Features

### SSH Hardening
- Password authentication disabled (key-only)
- Root login disabled (use deploy user with sudo)
- SSH banner with security notice
- Modern cryptographic algorithms (Curve25519, ChaCha20-Poly1305)
- Diffie-Hellman key exchange
- Secure ciphers and MACs (encrypt-then-MAC)
- Disabled X11 forwarding, agent forwarding, TCP forwarding
- Client alive interval (5 min) to timeout idle connections
- Fail2ban: 5 failed attempts = 1 hour ban

### Restricted Sudo
- Deploy user can only run specific commands without password:
  - `docker`, `docker-compose`
  - `systemctl` (docker)
  - Application scripts in `/opt/leviosa/scripts/`
- Optional password-based sudo for other commands
- See `/opt/leviosa/scripts/test-sudo.sh` to verify

### Firewall (UFW)
- Default deny incoming
- Allow SSH (port 22), HTTP (80), HTTPS (443) only
- Rate limit SSH connections

### Caddy Reverse Proxy
- Application bound to localhost (127.0.0.1) only
- Caddy handles external connections with automatic TLS
- Security headers: HSTS, X-Frame-Options, CSP
- Rate limiting per IP address
- Large file upload support (configurable)

### Docker Security
- Metrics endpoint bound to localhost only (127.0.0.1:9323)
- Users in docker group have root-equivalent access
- Only add trusted users to docker group

### System Updates
- Automatic security updates (unattended-upgrades)
- Daily package list updates
- Configurable reboot behavior

## Maintenance

### Update Application

```bash
# Pull latest image and restart
ansible-playbook playbooks/deploy.yml \
  -e "ansible_host=<server-ip>"
```

### Quick Restart

```bash
# Restart only the app container
make restart

# Restart staging
make restart-staging
```

### Manual Backup

```bash
ssh deploy@<server-ip>
cd /opt/leviosa
./scripts/backup-db.sh daily
```

### View Logs

```bash
ssh deploy@<server-ip>
docker compose logs -f
```

### Check Service Status

```bash
ssh deploy@<server-ip>
docker compose ps
```

## Troubleshooting

### Connection Refused

Make sure you can reach the server:
```bash
ping <server-ip>
ssh root@<server-ip>
```

### Permission Denied

Check your SSH key is configured:
```bash
cat ~/.ssh/leviosa.pub
```

### Playbook Fails

Run with verbose output:
```bash
ansible-playbook playbooks/site.yml -vvv
```

### Service Not Starting

Check logs:
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
- [ ] Deploy user has sudo access
- [ ] Docker containers running as non-root
- [ ] Environment files have correct permissions (0600)
- [ ] Backups configured
- [ ] Monitoring endpoints accessible (localhost only)

## Next Steps

1. Set up monitoring dashboards (Grafana, Prometheus)
2. Configure alerting (optional)
3. Set up CI/CD pipeline
4. Configure log aggregation (Loki)
5. Set up staging environment
