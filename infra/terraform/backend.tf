# Terraform State Backend Configuration
# This file creates the S3 bucket needed for Terraform state with workspace support
#
# IMPORTANT: These resources should be created ONCE before using workspaces
# After creating these resources, you can comment out this file to prevent recreation
#
# Initial setup:
# 1. Comment out the entire `backend "s3"` block in main.tf
# 2. Run: terraform init
# 3. Run: terraform apply
# 4. Uncomment the backend block in main.tf
# 5. Run: terraform init -reconfigure
# 6. (Optional) Comment out this file to prevent accidental recreation

resource "aws_s3_bucket" "terraform_state" {
  bucket = "leviosa-terraform-state"

  tags = {
    Name        = "Leviosa Terraform State"
    ManagedBy   = "Terraform"
  }
}

resource "aws_s3_bucket_versioning" "terraform_state_versioning" {
  bucket = aws_s3_bucket.terraform_state.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "terraform_state_sse" {
  bucket = aws_s3_bucket.terraform_state.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "terraform_state_block" {
  bucket = aws_s3_bucket.terraform_state.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# DynamoDB table for state locking - REMOVED: single-user workflow
# resource "aws_dynamodb_table" "terraform_locks" {
#   name         = "leviosa-terraform-locks"
#   billing_mode = "PAY_PER_REQUEST"
#   hash_key     = "LockID"
#
#   attribute {
#     name = "LockID"
#     type = "S"
#   }
#
#   tags = {
#     Name        = "Leviosa Terraform Locks"
#     Environment = "staging"
#     ManagedBy   = "Terraform"
#   }
# }

