# ============================================
# Hetzner Cloud Server
# ============================================

# Reference existing SSH key from Hetzner console
data "hcloud_ssh_key" "default" {
  name = "terraform-leviosa"
}

# Main VPS server (shared by both staging and production)
resource "hcloud_server" "manager" {
  name        = var.project_name
  image       = "ubuntu-24.04"
  server_type = var.server_type
  location    = var.server_location
  ssh_keys    = [data.hcloud_ssh_key.default.id]
  backups     = var.enable_backups

  labels = {
    project     = var.project_name
    environment = var.environment
    managed_by  = "terraform"
  }

  user_data = file("${path.module}/cloud-init.yml.tftpl")

  lifecycle {
    prevent_destroy = false
  }
}

# ============================================
# S3 Assets Bucket
# ============================================

locals {
  bucket_name = var.bucket_suffix != "" ? "${var.environment}-${var.project_name}-assets-${var.bucket_suffix}" : "${var.environment}-${var.project_name}-assets"
}

resource "aws_s3_bucket" "bucket" {
  bucket = local.bucket_name

  tags = {
    Name        = "${var.project_name} Assets"
    Environment = var.environment
    ManagedBy   = "Terraform"
  }
}

resource "aws_s3_bucket_versioning" "assets_versioning" {
  bucket = aws_s3_bucket.bucket.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "assets_encryption" {
  bucket = aws_s3_bucket.bucket.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "assets_block" {
  bucket = aws_s3_bucket.bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_cors_configuration" "assets_cors" {
  bucket = aws_s3_bucket.bucket.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD", "PUT", "POST"]
    allowed_origins = var.allowed_origins
    expose_headers  = ["ETag"]
    max_age_seconds = 3600
  }
}

# ============================================
# IAM User for Application S3 Access
# ============================================

resource "aws_iam_user" "app_user" {
  name = "${var.environment}-${var.project_name}-app"

  tags = {
    Environment = var.environment
    Project     = var.project_name
  }
}

resource "aws_iam_access_key" "app_user_key" {
  user = aws_iam_user.app_user.name
}

resource "aws_iam_user_policy" "app_s3_policy" {
  name = "S3Access"
  user = aws_iam_user.app_user.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ]
      Resource = [
        aws_s3_bucket.bucket.arn,
        "${aws_s3_bucket.bucket.arn}/*"
      ]
    }]
  })
}
