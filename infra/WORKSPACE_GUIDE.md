# Leviosa Workspace Guide

This guide explains how to use Terraform workspaces to manage staging and production environments for Leviosa.

## Architecture Overview

```
┌─────────────────────────────────────────────┐
│           Single Hetzner VPS                │
│  ┌──────────────────┐  ┌──────────────────┐ │
│  │   Staging        │  │   Production     │ │
│  │  - Port 3001     │  │  - Port 3000     │ │
│  │  - DB: leviosa_  │  │  - DB: leviosa_  │ │
│  │    staging       │  │    prod          │ │
│  │  - Redis DB 1    │  │  - Redis DB 0    │ │
│  │  - Dir: /opt/    │  │  - Dir: /opt/    │ │
│  │    leviosa-      │  │    leviosa       │ │
│  │    staging       │  │                  │ │
│  └──────────────────┘  └──────────────────┘ │
│         │                      │            │
│         └──────────┬───────────┘            │
│                    │                        │
│         Shared Services (Production)        │
│  - PostgreSQL (databases 0 and 1)          │
│  - Redis (DB indexes 0 and 1)              │
│  - Caddy (reverse proxy for both)          │
└─────────────────────────────────────────────┘
         │                         │
         └───────────┬─────────────┘
                     │
        Terraform (workspaces) → Outputs → Ansible
                     │
        ┌────────────┴────────────┐
        │  AWS Resources (per env) │
        │  - S3 buckets           │
        │  - CloudFront           │
        │  - IAM credentials      │
        └─────────────────────────┘
```

## Terraform Workspaces

### Workspace Strategy

- **default workspace** = staging environment
- **production workspace** = production environment

Each workspace has:
- Separate state file in S3
- Separate AWS resources (S3 buckets, IAM users, CloudFront distributions)
- Same VPS resources (Hetzner server, DNS records)

### Common Workflows

#### Initial Setup

```bash
cd /home/henga/Documents/projects/livio/leviosa/infra

# 1. Initialize Terraform backend (one time)
cd terraform
# Comment out backend "s3" block in main.tf
terraform init
terraform apply  # Creates state bucket
# Uncomment backend "s3" block
terraform init -reconfigure

# 2. Create production workspace
terraform workspace new production

# 3. Apply staging infrastructure (default workspace)
terraform workspace select default
terraform apply -var-file=terraform.tfvars.staging

# 4. Apply production infrastructure
terraform workspace select production
terraform apply -var-file=terraform.tfvars.production
```

#### Daily Operations

```bash
cd /home/henga/Documents/projects/livio/leviosa/infra

# Using the infra Makefile (recommended)

# Staging
make plan-staging
make apply-staging
make update-staging-vault
make deploy-staging

# Production
make plan-production
make apply-production
make update-production-vault
make deploy-production
```

#### Using Terraform Directly

```bash
cd terraform

# Staging (default workspace)
terraform workspace select default
terraform plan -var-file=terraform.tfvars.staging
terraform apply -var-file=terraform.tfvars.staging
../scripts/update-staging-vault.sh

# Production
terraform workspace select production
terraform plan -var-file=terraform.tfvars.production
terraform apply -var-file=terraform.tfvars.production
../scripts/update-production-vault.sh
```

### Vault Update Scripts

After applying Terraform changes, update Ansible vaults:

```bash
# Staging
./scripts/update-staging-vault.sh

# Production
./scripts/update-production-vault.sh
```

These scripts:
1. Switch to the appropriate workspace
2. Extract Terraform outputs (IAM credentials, S3 buckets, etc.)
3. Update Ansible vault files with the credentials

## Ansible Deployment

### Deployment Architecture

- **Production**: Full stack (app + postgres + redis + caddy)
- **Staging**: App only (shares postgres/redis/caddy from production)

### Environment Variables

| Variable | Production | Staging |
|----------|-----------|---------|
| `app_name` | leviosa | leviosa_staging |
| `app_env` | production | staging |
| `app_port` | 3000 | 3001 |
| `app_base_dir` | /opt/leviosa | /opt/leviosa-staging |
| `db_name` | leviosa | leviosa_staging |
| `redis_db` | 0 | 1 |
| `s3_bucket` | production-leviosa-assets | staging-leviosa-assets |

### Deployment Commands

```bash
cd ansible

# Production
make setup-production      # First time setup
make deploy-production     # Deploy application
make restart-production    # Quick restart

# Staging
make setup-staging         # First time setup
make deploy-staging        # Deploy application
make restart-staging       # Quick restart
```

## AWS Resources per Environment

### Staging Workspace (default)

- S3: `staging-leviosa-assets`
- S3: `staging-leviosa-backups`
- S3: `staging-leviosa-vault-storage`
- S3: `staging-leviosa-loki-logs`
- CloudFront: `cdn.leviosa.care`
- IAM: `vault-unseal`

### Production Workspace

- S3: `production-leviosa-assets`
- S3: `production-leviosa-backups`
- S3: `production-leviosa-vault-storage`
- S3: `production-leviosa-loki-logs`
- CloudFront: (separate distribution)
- IAM: (separate user)

## URL Structure

| Environment | URL | Notes |
|-------------|-----|-------|
| Production | https://leviosa.care | Main production site |
| Staging | https://staging.leviosa.care | Password-protected testing |

## Troubleshooting

### Terraform Issues

**Problem**: Can't find workspace
```bash
terraform workspace list
terraform workspace new production
```

**Problem**: State is locked
```bash
# Force unlock (use with caution!)
terraform force-unlock <LOCK_ID>
```

### Ansible Issues

**Problem**: Credentials out of date
```bash
# Re-run vault update
./scripts/update-staging-vault.sh
./scripts/update-production-vault.sh
```

**Problem**: Container won't start
```bash
# Check logs
make staging-logs    # or production-logs
make status          # Check service status
```

### DNS Issues

After Terraform apply, verify DNS records in Cloudflare:
1. Go to Cloudflare Dashboard
2. Select leviosa.care zone
3. Check DNS records for:
   - A record for staging.leviosa.care → VPS IP
   - CNAME for cdn.leviosa.care → CloudFront domain

## Best Practices

1. **Always use workspaces** when applying Terraform changes
2. **Update vaults after every Terraform apply**
3. **Test in staging first** before deploying to production
4. **Keep secrets separate** - never commit vault files
5. **Backup database** before major production changes

## Migration from Single Environment

If you're migrating from the old single-environment setup:

1. Backup current state: `terraform pull`
2. Create production workspace: `terraform workspace new production`
3. Move current resources to staging:
   - Current state becomes staging (default workspace)
   - Apply production workspace with `terraform.tfvars.production`
4. Update Ansible inventory for both environments
5. Deploy staging to `/opt/leviosa-staging`

## Further Reading

- [Terraform Workspaces](https://www.terraform.io/docs/language/settings/backends/s3.html#workspaces)
- [Ansible Inventory](https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html)
- [Docker Compose](https://docs.docker.com/compose/)
