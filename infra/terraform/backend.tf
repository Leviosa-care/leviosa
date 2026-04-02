resource "aws_s3_bucket" "terraform_state" {
  bucket        = "staging-leviosa-terraform-state"
  force_destroy = false
  tags = {
    Environment = "staging"
    Name        = "Terraform State Bucket"
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

