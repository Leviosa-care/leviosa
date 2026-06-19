# =============================================================================
# app-leviosa IAM User
# Least-privilege credential for leviosa's running application — S3-only
# access to leviosa's own assets and backups buckets in both environments.
# Distinct from terraform-leviosa (the Terraform admin user) and from the
# per-environment app_user resource in leviosa.tf.
#
# IAM users are account-global, but this config is applied across two
# Terraform workspaces (default=staging, production=production) that share
# these files. Gating on terraform.workspace ensures the user is created
# exactly once instead of colliding when applied in the second workspace.
# =============================================================================

locals {
  app_leviosa_bucket_names = [
    "production-${var.project_name}-assets-lc",
    "staging-${var.project_name}-assets",
    "staging-${var.project_name}-backups",
  ]
  app_leviosa_bucket_arns = [for b in local.app_leviosa_bucket_names : "arn:aws:s3:::${b}"]
  app_leviosa_object_arns = [for b in local.app_leviosa_bucket_names : "arn:aws:s3:::${b}/*"]
}

resource "aws_iam_policy" "app_leviosa" {
  count       = terraform.workspace == "production" ? 1 : 0
  name        = "app-leviosa-policy"
  description = "S3-only access to leviosa's own buckets (production + staging assets and backups)"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid      = "ListLeviosaBuckets"
        Effect   = "Allow"
        Action   = "s3:ListBucket"
        Resource = local.app_leviosa_bucket_arns
      },
      {
        Sid    = "ReadWriteDeleteLeviosaObjects"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = local.app_leviosa_object_arns
      }
    ]
  })
}

resource "aws_iam_user" "app_leviosa" {
  count = terraform.workspace == "production" ? 1 : 0
  name  = "app-leviosa"
  path  = "/applications/"

  tags = {
    Name      = "app-leviosa"
    Project   = var.project_name
    ManagedBy = "Terraform"
  }
}

resource "aws_iam_user_policy_attachment" "app_leviosa" {
  count      = terraform.workspace == "production" ? 1 : 0
  user       = aws_iam_user.app_leviosa[0].name
  policy_arn = aws_iam_policy.app_leviosa[0].arn
}

resource "aws_iam_access_key" "app_leviosa" {
  count = terraform.workspace == "production" ? 1 : 0
  user  = aws_iam_user.app_leviosa[0].name
}

output "app_leviosa_access_key_id" {
  value       = try(aws_iam_access_key.app_leviosa[0].id, null)
  description = "Access key ID for app-leviosa (only set when applied in the production workspace)"
  sensitive   = true
}

output "app_leviosa_access_key_secret" {
  value       = try(aws_iam_access_key.app_leviosa[0].secret, null)
  description = "Secret access key for app-leviosa (only set when applied in the production workspace)"
  sensitive   = true
}
