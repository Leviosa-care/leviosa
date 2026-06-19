# Leviosa Terraform Configuration

Infrastructure as Code for the Leviosa VPS and supporting AWS resources.

## What This Does

This Terraform configuration provisions:

- **VPS**: A single Hetzner Cloud server (`hcloud_server.manager`), shared by both the staging and production application environments (see [Workspace Strategy](#workspace-strategy))
- **DNS**: Cloudflare DNS records for the apex domain, `www`, `staging`, `admin`, `staff`, `admin-staging`, and `staff-staging` subdomains, all proxied through Cloudflare
- **S3 Assets Bucket**: One per environment (staging/production), versioned, encrypted, with public access blocked and CORS configured for the frontend origins
- **IAM**: A per-environment IAM user with least-privilege S3 access for the application, plus a separate `app-leviosa` IAM user (production workspace only) with access to both environments' asset and backup buckets
- **Terraform State Backend**: An S3 bucket (`leviosa-terraform-state`, defined in `backend.tf`) for remote state, created once and then left alone

## Workspace Strategy

This configuration uses Terraform workspaces instead of separate directories per environment:

- **`default` workspace = staging** — owns the shared VPS, Cloudflare DNS records, and staging S3/IAM resources
- **`production` workspace** — owns *only* the production S3/IAM resources (`app_leviosa_*` resources are also gated to this workspace via `count = terraform.workspace == "production" ? 1 : 0`)

The VPS and DNS records are **not** duplicated in the production workspace — they live once in staging/default and are shared by both app environments (see `infra/ansible` for how staging and production share the same VPS).

## Prerequisites

### Hetzner Cloud
1. Create a Hetzner Cloud account and project
2. Upload an SSH public key to the Hetzner console named `terraform-leviosa` (referenced via `data "hcloud_ssh_key" "default"` in `leviosa.tf` — it must already exist, Terraform does not create it)
3. Generate an API token for the project

### Cloudflare
1. Add your domain to Cloudflare and note the Zone ID
2. Generate an API token with DNS edit permissions for that zone

### AWS (for S3 + IAM only — no EC2 is used)
```bash
pip install awscli
aws configure --profile terraform-leviosa
```
The Makefile exports `AWS_PROFILE=terraform-leviosa` for all commands.

### Terraform
```bash
# On Ubuntu/Debian
sudo apt install terraform

# Or download from https://www.terraform.io/downloads
```

## Quick Start

### 1. Configure Variables

Copy the example and fill in values for each environment:

```bash
cp terraform.tfvars.example terraform.tfvars.staging
cp terraform.tfvars.example terraform.tfvars.production
```

```hcl
aws_region   = "eu-central-1"
project_name = "leviosa"
environment  = "staging"  # or "production"

hcloud_token    = "your_hetzner_api_token_here"
server_type     = "cpx11"
server_location = "nbg1"
enable_backups  = true

cloudflare_token = "your_cloudflare_api_token_here"
zone_id          = "your_zone_id_here"
domain_name      = "leviosa.care"

allowed_origins = [
  "http://localhost:5173",
  "https://leviosa.care",
  "https://staging.leviosa.care",
]
```

`terraform.tfvars`, `terraform.tfvars.staging`, and `terraform.tfvars.production` are gitignored (only `*.example` files are committed).

### 2. Initialize Terraform

```bash
make init-staging      # initializes and selects the default (staging) workspace
make init-production   # initializes and selects/creates the production workspace
```

### 3. Plan and Apply

```bash
make plan-staging
make apply-staging

make plan-production
make apply-production
```

Staging applies the VPS, Cloudflare DNS, and staging S3/IAM. Production applies only the production S3/IAM resources (see `STAGING_TARGETS` / `PRODUCTION_TARGETS` in the Makefile).

### 4. Get Outputs

```bash
make output            # all outputs
make credentials       # AWS S3 credentials
make server-info        # VPS connection details
make dns-info            # DNS records
```

## Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `aws_region` | AWS region for S3/IAM resources | `eu-central-1` | No |
| `project_name` | Project identifier | `leviosa` | No |
| `bucket_suffix` | Optional suffix to avoid S3 bucket name collisions | `""` | No |
| `environment` | `staging` or `production` | `staging` | No |
| `allowed_origins` | CORS origins for the S3 assets bucket | localhost + leviosa.care | No |
| `hcloud_token` | Hetzner Cloud API token | - | **Yes** |
| `server_type` | Hetzner Cloud server type | `cpx11` | No |
| `server_location` | Hetzner Cloud datacenter (`nbg1`, `fsn1`, `hel1`, `hil`, `ash`, `sin`) | `nbg1` | No |
| `enable_backups` | Enable Hetzner server backups | `true` | No |
| `cloudflare_token` | Cloudflare API token | - | **Yes** |
| `zone_id` | Cloudflare zone ID | - | **Yes** |
| `domain_name` | Primary domain name | - | **Yes** |

## Outputs

| Output | Description |
|--------|-------------|
| `domain_name` | Configured domain |
| `server_ipv4_address` / `server_ipv6_address` | VPS public IP addresses |
| `ssh_connection_string` | `ssh root@<ip>` |
| `media_bucket_name` / `media_bucket_region` / `media_bucket_arn` | S3 assets bucket details |
| `app_user_access_key_id` / `app_user_secret_access_key` | Per-environment S3 app credentials (sensitive) |
| `app_leviosa_access_key_id` / `app_leviosa_access_key_secret` | Cross-environment app credentials, production workspace only (sensitive) |
| `credentials` | Formatted summary of AWS credentials for Ansible vault |

## Directory Structure

```
infra/terraform/
├── main.tf              # Providers, S3 state backend config
├── variables.tf         # Variable definitions
├── leviosa.tf            # VPS, S3 assets bucket, per-environment IAM
├── app-leviosa.tf        # Cross-environment app-leviosa IAM user (production workspace only)
├── cloudflare.tf          # DNS records
├── backend.tf             # One-time bootstrap of the Terraform state S3 bucket
├── output.tf              # Output definitions
├── cloud-init.yml.tftpl   # Cloud-init template applied to the VPS
├── terraform.tfvars.example  # Template — copy to .staging/.production
├── Makefile               # Workspace-aware plan/apply/destroy commands
└── README.md              # This file
```

## Cost Estimates

Based on `eu-central-1` (AWS) / `nbg1` (Hetzner):

| Resource | Cost (monthly) |
|----------|----------------|
| Hetzner `cpx11` VPS | ~€4.59 |
| Hetzner server backups | ~€0.92 |
| S3 Storage (per environment, light usage) | ~$1–3 |
| Cloudflare DNS | Free |
| **Total (both environments)** | **~€10–15/month** |

*Prices are estimates and may vary by region/usage.*

## Security Best Practices

### SSH Keys
The VPS references an SSH key that must already exist in the Hetzner console (`terraform-leviosa`). Never commit private keys.

### AWS Credentials
- IAM users created have least-privilege access (S3 only)
- Credentials feed into Ansible Vault via `make update-staging-vault` / `make update-production-vault`
- Rotate credentials periodically

### State File
- Remote state lives in the `leviosa-terraform-state` S3 bucket (configured in `main.tf`/`backend.tf`)
- Never commit `.tfstate` files to git (already gitignored)
- State locking uses S3's native `use_lockfile` (no DynamoDB table — single-user workflow)

## Next Steps

After Terraform completes:

1. **Update Ansible Vault**: `make update-staging-vault` / `make update-production-vault`
2. **Run Ansible**: `make setup-staging` / `make setup-production` from `../ansible` (see `infra/ansible/README.md`)
3. **Verify DNS**: `make dns-verify`

## Troubleshooting

### "ssh_key not found" / hcloud_ssh_key data lookup fails
The SSH key named `terraform-leviosa` doesn't exist in your Hetzner project yet — upload it via the Hetzner console first.

### AWS "InsufficientPrivileges" / "AccessDenied"
Ensure the `terraform-leviosa` AWS profile has S3 Full Access and IAM Full Access.

### State Lock Timeout
```bash
terraform force-unlock <LOCK_ID>
```

### Destroy
```bash
make destroy             # destroys resources in the current workspace
make destroy-staging
make destroy-production
```

**Warning:** Destroying the staging/default workspace removes the shared VPS used by both environments.

## Maintenance

```bash
make fmt        # format files
make validate   # fmt -check + terraform validate
make refresh    # refresh state
make show       # show current state
make graph      # generate dependency graph
```
