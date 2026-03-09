# Leviosa Terraform Configuration

Infrastructure as Code for Leviosa VPS deployment on AWS.

## What This Does

This Terraform configuration provisions:

- **EC2 Instance**: Ubuntu 24.04 LTS VPS in your preferred region
- **Security Group**: Firewall rules for HTTP (80), HTTPS (443), and SSH (22)
- **S3 Buckets**:
  - Media storage bucket (for user uploads)
  - Backup bucket (for database backups)
- **IAM Resources**:
  - IAM user with programmatic access for S3
  - IAM policy for least-privilege access
- **DNS Management**:
  - Route53 hosted zone (if using Cloudflare, skip this)
  - DNS records pointing to your VPS

## Prerequisites

### AWS Account

1. Create an AWS account if you don't have one
2. Configure your AWS CLI:

```bash
pip install awscli
aws configure
```

### Terraform Installation

```bash
# On Ubuntu/Debian
sudo apt install terraform

# Or download from https://www.terraform.io/downloads
```

### Domain

You need a domain name. You can either:
- Use Route53 (AWS DNS service)
- Use Cloudflare (recommended for better DDoS protection)
- Use any other DNS provider

## Quick Start

### 1. Configure Variables

Create a `terraform.tfvars` file:

```hcl
# Region
aws_region = "eu-central-1"

# Project
project_name = "leviosa"
environment = "production"

# Domain
domain_name = "leviosa.care"
create_route53_zone = false  # Set to true if using Route53

# Instance
instance_type = "t3.medium"
ssh_public_key = "ssh-rsa AAAA... your-key-here"

# Backup
backup_bucket_name = "leviosa-production-backups"
```

### 2. Initialize Terraform

```bash
terraform init
```

### 3. Review the Plan

```bash
terraform plan
```

### 4. Deploy Infrastructure

```bash
terraform apply
```

Type `yes` when prompted to confirm.

### 5. Get Outputs

```bash
terraform output server_ipv4_address
terraform output domain_name
terraform output iam_access_key_id
terraform output iam_access_key_secret
```

## Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `aws_region` | AWS region | `eu-central-1` | No |
| `project_name` | Project identifier | `leviosa` | No |
| `environment` | Environment name | `production` | No |
| `domain_name` | Your domain name | - | **Yes** |
| `instance_type` | EC2 instance type | `t3.medium` | No |
| `ssh_public_key` | Your SSH public key | - | **Yes** |
| `create_route53_zone` | Create Route53 zone | `false` | No |
| `backup_bucket_name` | S3 backup bucket name | - | **Yes** |

## Outputs

| Output | Description | Example |
|--------|-------------|---------|
| `server_ipv4_address` | VPS public IP | `46.225.25.238` |
| `domain_name` | Configured domain | `leviosa.care` |
| `iam_access_key_id` | AWS access key | `AKIAIOSFODNN7EXAMPLE` |
| `iam_access_key_secret` | AWS secret key | `wJalrXUtnFEMI/K7MDENG...` |
| `backup_bucket_name` | S3 backup bucket | `leviosa-backups` |

## Directory Structure

```
infra/terraform/
├── main.tf              # Main configuration
├── variables.tf         # Variable definitions
├── outputs.tf           # Output definitions
├── terraform.tfvars     # Your variables (not in git)
├── README.md            # This file
└── modules/
    ├── ec2/             # EC2 instance module
    ├── s3/              # S3 buckets module
    ├── iam/             # IAM resources module
    └── security/        # Security groups module
```

## Cost Estimates

Based on `eu-central-1` region:

| Resource | Cost (monthly) |
|----------|----------------|
| t3.medium EC2 | ~$25 |
| S3 Storage (100GB) | ~$2.30 |
| S3 Requests | ~$0.50 |
| Data Transfer (1TB) | ~$80 |
| Route53 Hosted Zone | ~$0.50 |
| **Total** | **~$108/month** |

*Prices are estimates and may vary by region.*

## Security Best Practices

### SSH Keys

Never commit your private SSH key. Use `ssh_public_key` variable with the public key only.

### AWS Credentials

- The IAM user created has least-privilege access (S3 only)
- Store credentials securely in Ansible Vault
- Rotate credentials periodically

### State File

- `terraform.tfstate` contains sensitive information
- Use Terraform Cloud or S3 backend for state management
- Never commit `.tfstate` files to git

### Remote State (Recommended)

Configure S3 backend for state:

```hcl
terraform {
  backend "s3" {
    bucket         = "leviosa-terraform-state"
    key            = "production/terraform.tfstate"
    region         = "eu-central-1"
    encrypt        = true
    dynamodb_table = "leviosa-terraform-locks"
  }
}
```

## Next Steps

After Terraform completes:

1. **Copy Outputs**: Save the IAM credentials for Ansible configuration
2. **Run Ansible**: Use `make setup` in the `../ansible` directory
3. **Configure DNS**: Point your domain to the server IP
4. **Test Access**: SSH to the server and verify services

## Troubleshooting

### "InvalidKeyPair.NotFound"

Your SSH key name doesn't exist in AWS. Either:
- Upload your key to AWS EC2 console
- Use `ssh_public_key` variable instead (recommended)

### "InsufficientPrivileges"

Your AWS IAM user doesn't have permissions to create resources. Ensure you have:
- EC2 Full Access
- S3 Full Access
- IAM Full Access
- Route53 Full Access (if using Route53)

### State Lock Timeout

If another process is holding the lock:
```bash
terraform force-unlock <LOCK_ID>
```

### Destroy Everything

```bash
terraform destroy
```

**Warning:** This will delete all resources including the S3 buckets (backups will be lost).

## Maintenance

### Update Infrastructure

```bash
terraform plan -out=tfplan
terraform apply tfplan
```

### Import Existing Resources

If you have existing AWS resources:

```bash
terraform import aws_instance.example i-1234567890abcdef0
```

### Taint Resources

Force recreation of a resource:

```bash
terraform taint aws_instance.example
terraform apply
```

## Multiple Environments

Create separate workspaces for staging and production:

```bash
terraform workspace new staging
terraform apply

terraform workspace new production
terraform apply
```

Each workspace has its own state file.
