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

resource "cloudflare_dns_record" "cdn_leviosa_care" {
  zone_id = var.zone_id
  name    = "cdn.${var.domain_name}"
  type    = "CNAME"
  content = var.production_s3_bucket
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "cdn_staging_leviosa_care" {
  zone_id = var.zone_id
  name    = "cdn-staging.${var.domain_name}"
  type    = "CNAME"
  content = var.staging_s3_bucket
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "hostingermail_a_domainkey_leviosa_care" {
  zone_id = var.zone_id
  name    = "hostingermail-a._domainkey.${var.domain_name}"
  type    = "CNAME"
  content = "hostingermail-a.dkim.mail.hostinger.com"
  proxied = false
  ttl     = 18000
}

resource "cloudflare_dns_record" "hostingermail_b_domainkey_leviosa_care" {
  zone_id = var.zone_id
  name    = "hostingermail-b._domainkey.${var.domain_name}"
  type    = "CNAME"
  content = "hostingermail-b.dkim.mail.hostinger.com"
  proxied = false
  ttl     = 300
}

resource "cloudflare_dns_record" "hostingermail_c_domainkey_leviosa_care" {
  zone_id = var.zone_id
  name    = "hostingermail-c._domainkey.${var.domain_name}"
  type    = "CNAME"
  content = "hostingermail-c.dkim.mail.hostinger.com"
  proxied = false
  ttl     = 300
}

resource "cloudflare_dns_record" "leviosa_care_mx1" {
  zone_id  = var.zone_id
  name     = var.domain_name
  type     = "MX"
  content  = var.mx_servers[0].server
  priority = var.mx_servers[0].priority
  proxied  = false
  ttl      = 1
}

resource "cloudflare_dns_record" "leviosa_care_mx2" {
  zone_id  = var.zone_id
  name     = var.domain_name
  type     = "MX"
  content  = var.mx_servers[1].server
  priority = var.mx_servers[1].priority
  proxied  = false
  ttl      = 1
}

resource "cloudflare_dns_record" "dmarc_leviosa_care" {
  zone_id = var.zone_id
  name    = "_dmarc.${var.domain_name}"
  type    = "TXT"
  content = "\"v=DMARC1; p=none; rua=mailto:${var.contact_email}\""
  proxied = false
  ttl     = 1
}

resource "cloudflare_dns_record" "leviosa_care_txt" {
  zone_id = var.zone_id
  name    = "leviosa.care"
  type    = "TXT"
  content = "\"v=spf1 include:_spf.mail.hostinger.com ~all\""
  proxied = false
  ttl     = 3600
}
