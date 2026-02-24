# Database Backup Resources for Leviosa
# This file defines S3 resources for database backups

# ============================================
# S3 Bucket for Database Backups
# ============================================

resource "aws_s3_bucket" "backups" {
  bucket = "${var.environment}-${var.project_name}-backups"

  tags = {
    Name        = "${var.project_name} Database Backups"
    Environment = var.environment
    ManagedBy   = "Terraform"
    Purpose     = "Database Backups"
  }
}

# Enable versioning for backup bucket (recover from accidental overwrites)
resource "aws_s3_bucket_versioning" "backups_versioning" {
  bucket = aws_s3_bucket.backups.id

  versioning_configuration {
    status = "Enabled"
  }
}

# Server-side encryption for backup bucket
resource "aws_s3_bucket_server_side_encryption_configuration" "backups_encryption" {
  bucket = aws_s3_bucket.backups.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# Block all public access - backups should never be public
resource "aws_s3_bucket_public_access_block" "backups_block" {
  bucket = aws_s3_bucket.backups.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Lifecycle configuration for backup retention
resource "aws_s3_bucket_lifecycle_configuration" "backups_lifecycle" {
  bucket = aws_s3_bucket.backups.id

  # Daily backups - standard retention
  rule {
    id     = "daily-backup-retention"
    status = "Enabled"

    filter {
      prefix = "daily/"
    }

    # Move to cheaper storage after 30 days (minimum for STANDARD_IA)
    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    # Move to Glacier after 90 days (must be > last transition)
    transition {
      days          = 90
      storage_class = "GLACIER"
    }

    # Delete daily backups after retention period (must be > last transition)
    expiration {
      days = var.backup_retention_days + 1
    }

    # Clean up old versions
    noncurrent_version_expiration {
      noncurrent_days = 7
    }
  }

  # Weekly backups - longer retention
  rule {
    id     = "weekly-backup-retention"
    status = "Enabled"

    filter {
      prefix = "weekly/"
    }

    # Move to STANDARD_IA after 30 days (minimum required by AWS)
    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    # Move to Glacier after 60 days
    transition {
      days          = 60
      storage_class = "GLACIER"
    }

    # Keep weekly backups longer
    expiration {
      days = var.backup_retention_days * 2
    }

    noncurrent_version_expiration {
      noncurrent_days = 14
    }
  }

  # Monthly backups - longest retention
  rule {
    id     = "monthly-backup-retention"
    status = "Enabled"

    filter {
      prefix = "monthly/"
    }

    # Move to STANDARD_IA after 30 days
    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    # Move to Glacier after 90 days
    transition {
      days          = 90
      storage_class = "GLACIER"
    }

    # Keep monthly backups for a year
    expiration {
      days = 365
    }

    noncurrent_version_expiration {
      noncurrent_days = 30
    }
  }

  # Abort incomplete multipart uploads
  rule {
    id     = "abort-incomplete-uploads"
    status = "Enabled"

    filter {}

    abort_incomplete_multipart_upload {
      days_after_initiation = 1
    }
  }
}

# ============================================
# IAM Policy for Backup Access
# ============================================

resource "aws_iam_policy" "backup_access" {
  name        = "${var.project_name}-backup-access"
  description = "Policy for ${var.project_name} to manage database backups in S3"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "S3BackupAccess"
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = [
          aws_s3_bucket.backups.arn,
          "${aws_s3_bucket.backups.arn}/*"
        ]
      }
    ]
  })
}

# Attach backup policy to the vault IAM user
resource "aws_iam_user_policy_attachment" "backup_access_attach" {
  user       = aws_iam_user.vault_user.name
  policy_arn = aws_iam_policy.backup_access.arn
}
