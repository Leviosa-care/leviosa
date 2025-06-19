terraform {
  backend "s3" {
    bucket  = "staging-leviosa-terraform-state"
    key     = "terraform.tfstate"
    region  = "eu-central-1"
    encrypt = true
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
      version = "~> 5"

    }
  }
}

provider "aws" {
  region = "eu-central-1"
}

provider "hcloud" {
  token = var.hcloud_token
}

provider "cloudflare" {
  api_token = var.cloudflare_token
}
