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

# terraform output loki_s3_access_key_id
# terraform output loki_s3_secret_access_key


# terraform output vault_user_access_key_id
# terraform output vault_user_secret_access_key

# terraform output vault_kms_key_id
# terraform output vault_kms_key_arn

# terraform output vault_kms_key_arn
