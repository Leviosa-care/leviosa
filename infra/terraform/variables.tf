# TODO: where to use that one ?
variable "environment_name" {
  description = "Deployment environment (staging/production)"
  type        = string
  default     = "staging"
}
# Hetzner
variable "hcloud_token" {
  description = "The Hetzner Cloud API token"
  type        = string
  sensitive   = true
}

# Cloudflare
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

variable "production_s3_bucket" {
  description = "Production S3 bucket for CDN"
  type        = string
}

variable "staging_s3_bucket" {
  description = "Staging S3 bucket for CDN"
  type        = string
}

variable "contact_email" {
  description = "Contact email for DMARC reports"
  type        = string
}

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
