# ============================================
# Domain Output
# ============================================

output "domain_name" {
  value       = var.domain_name
  description = "Primary domain name"
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

# ============================================
# S3 Assets Bucket Outputs
# ============================================

output "media_bucket_name" {
  value       = aws_s3_bucket.bucket.id
  description = "S3 bucket name for assets"
}

output "media_bucket_region" {
  value       = aws_s3_bucket.bucket.region
  description = "S3 bucket region"
}

output "media_bucket_arn" {
  value       = aws_s3_bucket.bucket.arn
  description = "S3 bucket ARN"
}

# ============================================
# App IAM User Outputs
# ============================================

output "app_user_access_key_id" {
  value       = aws_iam_access_key.app_user_key.id
  description = "Access key ID for the app IAM user"
  sensitive   = true
}

output "app_user_secret_access_key" {
  value       = aws_iam_access_key.app_user_key.secret
  description = "Secret access key for the app IAM user"
  sensitive   = true
}

output "credentials" {
  value = <<-EOT
    ==========================================
    AWS Credentials for Ansible vault
    ==========================================

    AWS_ACCESS_KEY_ID     = ${aws_iam_access_key.app_user_key.id}
    AWS_SECRET_ACCESS_KEY = (run: terraform output -raw app_user_secret_access_key)
    S3_BUCKET             = ${aws_s3_bucket.bucket.id}
    S3_REGION             = ${aws_s3_bucket.bucket.region}

    ==========================================
  EOT
  description = "Credentials summary for Ansible vault update"
  sensitive   = true
}
