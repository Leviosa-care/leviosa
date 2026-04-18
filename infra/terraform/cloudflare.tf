# ============================================
# Application DNS Records
# ============================================

resource "cloudflare_dns_record" "leviosa_care_a" {
  zone_id = var.zone_id
  name    = var.domain_name
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
