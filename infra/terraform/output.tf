# ============================================
# Vault/Auto-unseal Outputs (existing)
# ============================================

output "vault_user_access_key_id" {
  value       = aws_iam_access_key.vault_user_key.id
  description = "Access key ID for Vault auto-unseal"
  sensitive   = true
}

output "vault_user_secret_access_key" {
  value       = aws_iam_access_key.vault_user_key.secret
  description = "Secret access key for Vault auto-unseal"
  sensitive   = true
}

output "vault_kms_key_id" {
  value       = aws_kms_key.vault_auto_unseal.key_id
  description = "KMS Key ID for Vault"
}

output "vault_kms_key_arn" {
  value       = aws_kms_key.vault_auto_unseal.arn
  description = "KMS Key ARN for Vault"
}

# ============================================
# Loki S3 Outputs (existing)
# ============================================

output "loki_s3_access_key_id" {
  value       = aws_iam_access_key.loki_s3_user_key.id
  description = "Access key ID for Loki S3 bucket"
  sensitive   = true
}

output "loki_s3_secret_access_key" {
  value       = aws_iam_access_key.loki_s3_user_key.secret
  description = "Access secret key for Loki S3 bucket"
  sensitive   = true
}

output "loki_s3_bucket_name" {
  value       = aws_s3_bucket.loki_logs.id
  description = "Loki S3 bucket name"
}

output "loki_s3_bucket_arn" {
  value       = aws_s3_bucket.loki_logs.arn
  description = "Loki S3 bucket ARN"
}

output "loki_s3_bucket_region" {
  value       = aws_s3_bucket.loki_logs.region
  description = "Loki S3 bucket region"
}

# ============================================
# S3 Assets Bucket Outputs
# ============================================

output "s3_bucket_name" {
  value       = aws_s3_bucket.bucket.id
  description = "S3 bucket name for assets"
}

output "s3_bucket_arn" {
  value       = aws_s3_bucket.bucket.arn
  description = "S3 bucket ARN for assets"
}

output "s3_bucket_region" {
  value       = aws_s3_bucket.bucket.region
  description = "S3 bucket region"
}

output "s3_public_url" {
  value       = aws_s3_bucket.bucket.bucket_regional_domain_name
  description = "S3 bucket public URL"
}

# ============================================
# CloudFront CDN Outputs
# ============================================

output "cloudfront_distribution_id" {
  value       = aws_cloudfront_distribution.cdn.id
  description = "CloudFront distribution ID for cache invalidation"
}

output "cloudfront_domain_name" {
  value       = aws_cloudfront_distribution.cdn.domain_name
  description = "CloudFront distribution domain name"
}

output "cdn_url" {
  value       = "https://cdn.${var.domain_name}"
  description = "CDN URL for application configuration"
}

# ============================================
# Database Backup Outputs
# ============================================

output "backup_bucket_name" {
  value       = aws_s3_bucket.backups.bucket
  description = "Name of the database backup S3 bucket"
}

output "backup_bucket_arn" {
  value       = aws_s3_bucket.backups.arn
  description = "ARN of the database backup S3 bucket"
}

output "backup_bucket_region" {
  value       = aws_s3_bucket.backups.region
  description = "AWS region of the backup bucket"
}

output "backup_env_variables" {
  value = {
    BACKUP_S3_BUCKET = aws_s3_bucket.backups.bucket
    BACKUP_S3_REGION = aws_s3_bucket.backups.region
  }
  description = "Environment variables for backup script"
}

# ============================================
# SES Email Outputs
# ============================================

output "ses_verification_token" {
  description = "SES domain verification token"
  value       = aws_ses_domain_identity.main.verification_token
}

output "ses_dkim_tokens" {
  description = "SES DKIM tokens"
  value       = aws_ses_domain_dkim.main.dkim_tokens
}

output "ses_mail_from_domain" {
  description = "SES MAIL FROM domain"
  value       = aws_ses_domain_mail_from.main.mail_from_domain
}

output "ses_configuration_set" {
  description = "SES configuration set name"
  value       = aws_ses_configuration_set.main.name
}

output "email_setup_status" {
  value       = <<-EOT
    ========================================
    Email Configuration
    (Existing Mailbox + Amazon SES API)
    ========================================

    Domain: ${var.domain_name}

    AUTOMATED (Terraform-managed):
    ----------------------------------------
    - SES domain identity & verification
    - SES DKIM authentication (3 CNAME records)
    - MAIL FROM domain (MX + SPF records)
    - SES verification TXT record
    - SES sending policy (attached to vault_user)
    - Mailbox provider DKIM (${length(var.email_dkim_records)} records configured)

    SENDING: Amazon SES API
    ----------------------------------------
    Uses existing vault_user IAM credentials (no SMTP credentials needed).
    Application .env:
      AWS_ACCESS_KEY_ID=<from terraform output vault_user_access_key_id>
      AWS_SECRET_ACCESS_KEY=<from terraform output vault_user_secret_access_key>
      SES_FROM_EMAIL=noreply@${var.domain_name}
      SES_FROM_NAME=Leviosa
      SES_REGION=${var.aws_region}

    RECEIVING: Existing Mailbox Provider
    ----------------------------------------
    Your existing mailbox configuration at:
    ${var.mx_servers[0].server} (priority ${var.mx_servers[0].priority})
    ${var.mx_servers[1].server} (priority ${var.mx_servers[1].priority})

    SPF Record: v=spf1 ${join(" ", [for d in var.email_spf_includes : "include:${d}"])} ~all
    DMARC Record: v=DMARC1; p=${var.email_dmarc_policy}; rua=mailto:${var.contact_email}

    REMAINING MANUAL STEPS:
      1. Request SES production access in AWS Console
      2. Backend code changes to use SES API instead of Gmail SMTP
      3. Verify email provider DKIM records in Cloudflare

    See: infra/terraform/ses.tf for details
    ========================================
    EOT
  description = "Email configuration status and instructions"
}

# ============================================
# Server Outputs
# ============================================

output "server_ipv4_address" {
  value       = hcloud_server.manager.ipv4_address
  description = "Public IPv4 address of the server"
}

output "server_ipv6_address" {
  value       = hcloud_server.manager.ipv6_address
  description = "Public IPv6 address of the server"
}

output "ssh_connection_string" {
  value       = "ssh root@${hcloud_server.manager.ipv4_address}"
  description = "SSH connection string for the server"
}

output "iam_user_name" {
  value       = aws_iam_user.vault_user.name
  description = "Name of the IAM user for application access"
}
