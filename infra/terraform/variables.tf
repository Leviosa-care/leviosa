# ============================================
# Core Variables
# ============================================

variable "environment" {
  description = "Deployment environment (staging, production)"
  type        = string
  default     = "staging"

  validation {
    condition     = contains(["staging", "production"], var.environment)
    error_message = "Environment must be one of: staging, production."
  }
}

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "leviosa"
}

variable "bucket_suffix" {
  description = "Optional unique suffix appended to the S3 bucket name to avoid global name collisions"
  type        = string
  default     = ""
}

variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "eu-central-1"
}

variable "allowed_origins" {
  description = "Allowed CORS origins for S3 bucket"
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
