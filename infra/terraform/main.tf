# Leviosa Terraform Configuration
#
# Workspace Strategy:
# - default workspace = staging
# - production workspace = production

terraform {
  backend "s3" {
    bucket  = "leviosa-terraform-state"
    key     = "terraform.tfstate"
    region  = "eu-central-1"
    encrypt = true
    use_lockfile = true
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

  required_version = ">= 1.0"
}

provider "aws" {
  region = var.aws_region
}

provider "hcloud" {
  token = var.hcloud_token
}

provider "cloudflare" {
  api_token = var.cloudflare_token
}
