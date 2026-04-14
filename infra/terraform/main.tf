# Leviosa Terraform Configuration
# This file configures the Terraform backend and providers
#
# Workspace Strategy:
# - default workspace = staging
# - production workspace = production
#
# Each workspace has separate state and AWS resources

terraform {
  # S3 backend for state storage with workspace support
  # State files: terraform.tfstate (staging/default), production-terraform.tfstate (production)
  backend "s3" {
    bucket         = "leviosa-terraform-state"
    key            = "terraform.tfstate"
    region         = "eu-central-1"
    encrypt        = true
    # dynamodb_table = "leviosa-terraform-locks"  # Removed: single-user workflow
    use_lockfile   = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = "1.51.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.0"
    }
  }

  # Required Terraform version
  required_version = ">= 1.0"
}

# ============================================
# Providers
# ============================================

provider "aws" {
  region = var.aws_region
}

provider "hcloud" {
  token = var.hcloud_token
}

provider "cloudflare" {
  api_token = var.cloudflare_token
}

# AWS Provider for us-east-1 (required for ACM certificates in CloudFront)
provider "aws" {
  alias  = "us_east_1"
  region = "us-east-1"
}
