# ============================================
# Hetzner Cloud Server
# ============================================

# Reference existing SSH key from Hetzner console
# Update the name to match your SSH key in Hetzner
data "hcloud_ssh_key" "default" {
  name = "terraform-leviosa"
}

# Main VPS server
resource "hcloud_server" "manager" {
  name        = "${var.project_name}-${var.environment}"
  image       = "ubuntu-24.04"
  server_type = var.server_type
  location    = var.server_location
  ssh_keys    = [data.hcloud_ssh_key.default.id]

  # Enable automatic backups
  backups = var.enable_backups

  # Labels for organization
  labels = {
    project     = var.project_name
    environment = var.environment
    managed_by  = "terraform"
  }

  # User data for initial server setup
  user_data = file("${path.module}/cloud-init.yml.tftpl")

  # Don't destroy server if it has important data
  lifecycle {
    # Set to true in production to prevent accidental deletion
    prevent_destroy = false
  }
}

# ============================================
# S3 Assets Bucket
# ============================================

resource "aws_s3_bucket" "bucket" {
  bucket = "${var.environment}-${var.project_name}-assets"

  tags = {
    Name        = "${var.project_name} Assets"
    Environment = var.environment
    ManagedBy   = "Terraform"
    Purpose     = "CDN Assets"
  }
}

# Enable versioning for assets bucket
resource "aws_s3_bucket_versioning" "assets_versioning" {
  bucket = aws_s3_bucket.bucket.id

  versioning_configuration {
    status = "Enabled"
  }
}

# Server-side encryption for assets bucket
resource "aws_s3_bucket_server_side_encryption_configuration" "assets_encryption" {
  bucket = aws_s3_bucket.bucket.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# Block public access (access via CloudFront only)
resource "aws_s3_bucket_public_access_block" "assets_block" {
  bucket = aws_s3_bucket.bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# CORS configuration for assets bucket
resource "aws_s3_bucket_cors_configuration" "assets_cors" {
  bucket = aws_s3_bucket.bucket.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = var.allowed_origins
    expose_headers  = ["ETag"]
    max_age_seconds = 3600
  }
}

# Lifecycle configuration for assets
resource "aws_s3_bucket_lifecycle_configuration" "assets_lifecycle" {
  bucket = aws_s3_bucket.bucket.id

  rule {
    id     = "asset-cleanup"
    status = "Enabled"

    filter {
      prefix = ""
    }

    noncurrent_version_expiration {
      noncurrent_days = 30
    }

    abort_incomplete_multipart_upload {
      days_after_initiation = 1
    }
  }
}

# Bucket policy for CloudFront OAC access
resource "aws_s3_bucket_policy" "assets_policy" {
  bucket = aws_s3_bucket.bucket.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowCloudFrontOACAccess"
        Effect = "Allow"
        Principal = {
          Service = "cloudfront.amazonaws.com"
        }
        Action   = "s3:GetObject"
        Resource = "${aws_s3_bucket.bucket.arn}/*"
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = aws_cloudfront_distribution.cdn.arn
          }
        }
      }
    ]
  })

  depends_on = [aws_cloudfront_distribution.cdn]
}
