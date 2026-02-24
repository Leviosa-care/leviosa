# ============================================
# Amazon SES Configuration
# ============================================
# SES resources for sending emails: domain verification, DKIM, IAM policy.
# DNS records (DKIM, verification, MAIL FROM) are created automatically in cloudflare.tf.
#
# The application sends email via the SES API (aws-sdk), using the same IAM
# credentials already provisioned for S3 (app_user). No SMTP credentials needed.
#
# REMAINING MANUAL STEP:
# 1. Request SES production access in AWS Console (to send to non-verified emails)
# 2. Backend code changes required to use SES API instead of Gmail SMTP
# ============================================

# ============================================
# SES Domain Identity
# ============================================

resource "aws_ses_domain_identity" "main" {
  domain = var.domain_name
}

# ============================================
# SES DKIM Authentication
# ============================================

resource "aws_ses_domain_dkim" "main" {
  domain = aws_ses_domain_identity.main.domain
}

# ============================================
# SES Configuration Set
# ============================================

resource "aws_ses_configuration_set" "main" {
  name = "${var.environment}-${var.project_name}-config"

  sending_enabled = true

  reputation_metrics_enabled = true
}

# ============================================
# SES Sending Policy (IAM)
# ============================================

resource "aws_iam_policy" "ses_send" {
  name        = "${var.project_name}-ses-send"
  description = "Policy for ${var.project_name} to send emails via SES"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "SESSendAccess"
        Effect = "Allow"
        Action = [
          "ses:SendEmail",
          "ses:SendRawEmail"
        ]
        Resource = "*"
      },
      {
        Sid    = "SESVerifyAccess"
        Effect = "Allow"
        Action = [
          "ses:VerifyEmailIdentity",
          "ses:GetIdentityVerificationAttributes",
          "ses:GetIdentityMailFromDomainAttributes"
        ]
        Resource = "*"
      }
    ]
  })
}

# Attach SES policy to existing vault user
resource "aws_iam_user_policy_attachment" "ses_send_attach" {
  user       = aws_iam_user.vault_user.name
  policy_arn = aws_iam_policy.ses_send.arn
}

# ============================================
# SES MAIL FROM Domain
# ============================================

resource "aws_ses_domain_mail_from" "main" {
  domain           = aws_ses_domain_identity.main.domain
  mail_from_domain = "mail.${aws_ses_domain_identity.main.domain}"
}
