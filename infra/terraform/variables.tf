# ============================================
# Core Variables
# ============================================

variable "environment" {
  description = "Deployment environment (development, staging, production)"
  type        = string
  default     = "staging"

  validation {
    condition     = contains(["development", "staging", "production"], var.environment)
    error_message = "Environment must be one of: development, staging, production."
  }
}

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "leviosa"
}

variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "eu-central-1"
}

variable "allowed_origins" {
  description = "Allowed CORS origins for S3 bucket (e.g., your application domains)"
  type        = list(string)
  default = [
    "http://localhost:5173",
    "https://leviosa.care",
    "https://staging.leviosa.care",
  ]
}

# ============================================
# Hetzner Cloud Variables
# ============================================

variable "hcloud_token" {
  description = "The Hetzner Cloud API token"
  type        = string
  sensitive   = true
}

variable "server_type" {
  description = "Hetzner Cloud server type (cpx11, cpx21, cpx31, etc.)"
  type        = string
  default     = "cpx11"

  validation {
    condition     = can(regex("^cpx[0-9]{2}$", var.server_type))
    error_message = "Server type must be a valid CX plan (e.g., cpx11, cpx21, cpx31)."
  }
}

variable "server_location" {
  description = "Hetzner Cloud datacenter location"
  type        = string
  default     = "nbg1"

  validation {
    condition     = contains(["nbg1", "fsn1", "hel1", "hil", "ash", "sin"], var.server_location)
    error_message = "Location must be a valid Hetzner datacenter code."
  }
}

variable "enable_backups" {
  description = "Enable automatic backups for the server (additional cost)"
  type        = bool
  default     = true
}

# ============================================
# Cloudflare DNS Variables
# ============================================

variable "cloudflare_token" {
  description = "The Cloudflare API token"
  type        = string
  sensitive   = true
}

variable "zone_id" {
  description = "The Cloudflare zone ID"
  type        = string
}

variable "domain_name" {
  description = "The primary domain name"
  type        = string
}

variable "contact_email" {
  description = "Contact email for DMARC reports"
  type        = string
}

variable "production_s3_bucket" {
  description = "Production S3 bucket for CDN (legacy, for staging CDN record)"
  type        = string
}

variable "staging_s3_bucket" {
  description = "Staging S3 bucket for CDN"
  type        = string
}

# ============================================
# Email Configuration Variables
# ============================================

variable "mx_servers" {
  description = "Mail exchange servers"
  type = list(object({
    server   = string
    priority = number
  }))
  default = [
    {
      server   = "mx1.hostinger.com"
      priority = 10
    },
    {
      server   = "mx2.hostinger.com"
      priority = 20
    }
  ]
}

variable "email_spf_includes" {
  description = "Domains to include in SPF record (e.g., amazonses.com, _spf.hostinger.com)"
  type        = list(string)
  default     = ["amazonses.com"]
}

variable "email_dmarc_policy" {
  description = "DMARC policy for the domain (none, quarantine, reject)"
  type        = string
  default     = "quarantine"

  validation {
    condition     = contains(["none", "quarantine", "reject"], var.email_dmarc_policy)
    error_message = "DMARC policy must be one of: none, quarantine, reject."
  }
}

variable "email_dkim_records" {
  description = "Mailbox provider DKIM records (from your email provider's DNS settings page)"
  type        = map(object({ type = string, content = string }))
  default = {
    "hostingermail-a._domainkey" = {
      type    = "CNAME"
      content = "hostingermail-a.dkim.mail.hostinger.com"
    }
    "hostingermail-b._domainkey" = {
      type    = "CNAME"
      content = "hostingermail-b.dkim.mail.hostinger.com"
    }
    "hostingermail-c._domainkey" = {
      type    = "CNAME"
      content = "hostingermail-c.dkim.mail.hostinger.com"
    }
  }
}

# ============================================
# Backup Configuration Variables
# ============================================

variable "backup_retention_days" {
  description = "Number of days to retain daily database backups before deletion"
  type        = number
  default     = 90

  validation {
    condition     = var.backup_retention_days >= 7 && var.backup_retention_days <= 365
    error_message = "Backup retention must be between 7 and 365 days."
  }
}
