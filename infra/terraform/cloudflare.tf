# ============================================
# Application DNS Records
# ============================================

resource "cloudflare_dns_record" "leviosa_care_a" {
  zone_id = var.zone_id
  name    = "leviosa.care"
  type    = "A"
  content = hcloud_server.manager.ipv4_address
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "leviosa_care_aaaa" {
  zone_id = var.zone_id
  name    = var.domain_name
  type    = "AAAA"
  content = hcloud_server.manager.ipv6_address
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "www_leviosa_care" {
  zone_id = var.zone_id
  name    = "www.${var.domain_name}"
  type    = "CNAME"
  content = var.domain_name
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "staging_leviosa_care" {
  zone_id = var.zone_id
  name    = "staging.${var.domain_name}"
  type    = "CNAME"
  content = var.domain_name
  proxied = true
  ttl     = 1
}

# ============================================
# CDN DNS Record (CloudFront)
# ============================================

# Changed from S3 CNAME (proxied) to CloudFront domain CNAME (not proxied)
resource "cloudflare_dns_record" "cdn_leviosa_care" {
  zone_id = var.zone_id
  name    = "cdn.${var.domain_name}"
  type    = "CNAME"
  content = aws_cloudfront_distribution.cdn.domain_name
  proxied = false # Cannot proxy CloudFront
  ttl     = 3600
}

resource "cloudflare_dns_record" "cdn_staging_leviosa_care" {
  zone_id = var.zone_id
  name    = "cdn-staging.${var.domain_name}"
  type    = "CNAME"
  content = var.staging_s3_bucket
  proxied = true
  ttl     = 1
}

# ============================================
# Email DNS Records
# ============================================
# Existing mailbox provider (e.g., Hostinger) for receiving + Amazon SES API for sending
#
# SENDING: Amazon SES via AWS SDK (uses existing IAM credentials)
# RECEIVING: Existing mailbox provider (configured in MX records)
#
# MX records route inbound mail to existing provider.
# SPF authorizes all providers in var.email_spf_includes.
# SES DKIM is automated; existing provider DKIM added via var.email_dkim_records.
# ============================================

# MX records - route inbound email to existing mailbox provider
resource "cloudflare_dns_record" "leviosa_care_mx1" {
  zone_id  = var.zone_id
  name     = var.domain_name
  type     = "MX"
  content  = var.mx_servers[0].server
  priority = var.mx_servers[0].priority
  proxied  = false
  ttl      = 3600
}

resource "cloudflare_dns_record" "leviosa_care_mx2" {
  zone_id  = var.zone_id
  name     = var.domain_name
  type     = "MX"
  content  = var.mx_servers[1].server
  priority = var.mx_servers[1].priority
  proxied  = false
  ttl      = 3600
}

# SPF record - variable-driven, authorizes sending providers
resource "cloudflare_dns_record" "leviosa_care_txt" {
  zone_id = var.zone_id
  name    = var.domain_name
  type    = "TXT"
  content = "\"v=spf1 ${join(" ", [for d in var.email_spf_includes : "include:${d}"])} ~all\""
  proxied = false
  ttl     = 3600
}

# ============================================
# Amazon SES DNS Records (automated from ses.tf)
# ============================================

# SES domain verification TXT record
resource "cloudflare_dns_record" "ses_verification" {
  zone_id = var.zone_id
  name    = "_amazonses.${var.domain_name}"
  type    = "TXT"
  content = "\"${aws_ses_domain_identity.main.verification_token}\""
  proxied = false
  ttl     = 3600
}

# SES DKIM CNAME records (AWS always generates exactly 3 tokens)
resource "cloudflare_dns_record" "ses_dkim" {
  count = 3

  zone_id = var.zone_id
  name    = "${aws_ses_domain_dkim.main.dkim_tokens[count.index]}._domainkey"
  type    = "CNAME"
  content = "${aws_ses_domain_dkim.main.dkim_tokens[count.index]}.dkim.amazonses.com"
  proxied = false
  ttl     = 3600
}

# MAIL FROM MX record
resource "cloudflare_dns_record" "ses_mail_from_mx" {
  zone_id  = var.zone_id
  name     = "mail.${var.domain_name}"
  type     = "MX"
  content  = "feedback-smtp.${var.aws_region}.amazonses.com"
  priority = 10
  proxied  = false
  ttl      = 3600
}

# MAIL FROM SPF record
resource "cloudflare_dns_record" "ses_mail_from_spf" {
  zone_id = var.zone_id
  name    = "mail.${var.domain_name}"
  type    = "TXT"
  content = "\"v=spf1 include:amazonses.com ~all\""
  proxied = false
  ttl     = 3600
}

# ============================================
# Existing Mailbox Provider DKIM Records
# ============================================

resource "cloudflare_dns_record" "mailbox_dkim" {
  for_each = var.email_dkim_records

  zone_id = var.zone_id
  name    = each.key
  type    = each.value.type
  content = each.value.type == "TXT" ? "\"${each.value.content}\"" : each.value.content
  proxied = false
  ttl     = 3600
}

# ============================================
# DMARC Record
# ============================================

resource "cloudflare_dns_record" "dmarc_leviosa_care" {
  zone_id = var.zone_id
  name    = "_dmarc.${var.domain_name}"
  type    = "TXT"
  content = "\"v=DMARC1; p=${var.email_dmarc_policy}; rua=mailto:${var.contact_email}\""
  proxied = false
  ttl     = 3600
}
