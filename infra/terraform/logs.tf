resource "aws_s3_bucket" "loki_logs" {
  bucket        = "staging-leviosa-loki-logs"
  force_destroy = false

  tags = {
    Name        = "Loki Logs Bucket"
    Environment = "staging"
    Project     = "leviosa"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "loki_logs_lifecycle" {
  bucket = aws_s3_bucket.loki_logs.id

  rule {
    id     = "log-expiry"
    status = "Enabled"

    filter {
      prefix = "" # Empty prefix means apply to all objects
    }

    expiration {
      days = 90
    }

    transition {
      days          = 30
      storage_class = "GLACIER"
    }

    noncurrent_version_expiration {
      newer_noncurrent_versions = 1
      noncurrent_days           = 30
    }
  }
}


# Enable versioning
resource "aws_s3_bucket_versioning" "log_logs_versioning" {
  bucket = aws_s3_bucket.loki_logs.id

  versioning_configuration {
    status = "Enabled"
  }
}

# Enable server-side encryption (SSE-S3)
resource "aws_s3_bucket_server_side_encryption_configuration" "loki_logs_encryption" {
  bucket = aws_s3_bucket.loki_logs.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "loki_logs" {
  bucket = aws_s3_bucket.loki_logs.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_iam_user" "loki_s3_user" {
  name = "loki-s3-uploader"
  tags = {
    Environment = "staging"
    Project     = "leviosa"
  }
}

resource "aws_iam_user_policy" "loki_s3_user_policy" {
  name = "LokiS3AccessPolicy"
  user = aws_iam_user.loki_s3_user.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:ListBucket",
          "s3:DeleteObject"
        ]
        Resource = [
          "${aws_s3_bucket.loki_logs.arn}",
          "${aws_s3_bucket.loki_logs.arn}/*"
        ]
      },
    ]
  })
}

resource "aws_iam_access_key" "loki_s3_user_key" {
  user = aws_iam_user.loki_s3_user.name
}

