# vault user
resource "aws_iam_user" "vault_user" {
  name = "vault-unseal"
}

resource "aws_iam_access_key" "vault_user_key" {
  user = aws_iam_user.vault_user.name
}

# KMS key
resource "aws_kms_key" "vault_auto_unseal" {
  description             = "KMS key for Vault auto-unseal"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_kms_alias" "vault_auto_unseal_alias" {
  name          = "alias/vault-auto-unseal"
  target_key_id = aws_kms_key.vault_auto_unseal.id
}

resource "aws_iam_policy" "vault_access_policy" {
  name = "vault-access-policy"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      # KMS access
      {
        Effect = "Allow"
        Action = [
          "kms:Decrypt",
          "kms:Encrypt",
          "kms:GenerateDataKey",
          "kms:DescribeKey"
        ]
        Resource = aws_kms_key.vault_auto_unseal.arn
      },
      # S3 access to Vault storage bucket
      {
        Effect = "Allow",
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ],
        Resource = [
          aws_s3_bucket.vault_storage.arn,
          "${aws_s3_bucket.vault_storage.arn}/*"
        ]
      }
    ]
  })
}

resource "aws_iam_user_policy_attachment" "vault_user_attachment" {
  user       = aws_iam_user.vault_user.name
  policy_arn = aws_iam_policy.vault_access_policy.arn
}

# S3 vault storage
resource "aws_s3_bucket" "vault_storage" {
  bucket        = "${var.environment}-${var.project_name}-vault-storage"
  force_destroy = false
  tags = {
    Environment = var.environment
    Name        = "Vault Storage Bucket"
  }
}

# Enable versioning
resource "aws_s3_bucket_versioning" "vault_storage_versioning" {
  bucket = aws_s3_bucket.vault_storage.id

  versioning_configuration {
    status = "Enabled"
  }
}

# Enable server-side encryption (SSE-S3)
resource "aws_s3_bucket_server_side_encryption_configuration" "vault_storage_encryption" {
  bucket = aws_s3_bucket.vault_storage.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}
