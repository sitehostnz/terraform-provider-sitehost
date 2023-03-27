terraform {
  required_providers {
    sitehost = {
      source  = "sitehostnz/sitehost"
      version = "~> 1.0"
    }
  }
}

provider "sitehost" {
  client_id    = var.sitehost_client_id
  api_key      = var.sitehost_api_key
  api_endpoint = "https://api.sitehost.nz/1.1/"
}

resource "sitehost_dns_zone" "zone" {
  name = var.domain_name
}

resource "sitehost_dns_record" "mx_main" {
  domain   = sitehost_dns_zone.zone.name
  name     = sitehost_dns_zone.zone.name
  type     = "MX"
  priority = 1
  record   = "ASPMX.L.GOOGLE.COM"
}

resource "sitehost_dns_record" "mx_alt1" {
  domain   = sitehost_dns_zone.zone.name
  name     = sitehost_dns_zone.zone.name
  type     = "MX"
  priority = 5
  record   = "ALT1.ASPMX.L.GOOGLE.COM"
}

resource "sitehost_dns_record" "mx_alt2" {
  domain   = sitehost_dns_zone.zone.name
  name     = sitehost_dns_zone.zone.name
  type     = "MX"
  priority = 5
  record   = "ALT2.ASPMX.L.GOOGLE.COM"
}

resource "sitehost_dns_record" "mx_alt3" {
  domain   = sitehost_dns_zone.zone.name
  name     = sitehost_dns_zone.zone.name
  type     = "MX"
  priority = 10
  record   = "ALT3.ASPMX.L.GOOGLE.COM"
}

resource "sitehost_dns_record" "mx_alt4" {
  domain   = sitehost_dns_zone.zone.name
  name     = sitehost_dns_zone.zone.name
  type     = "MX"
  priority = 10
  record   = "ALT4.ASPMX.L.GOOGLE.COM"
}

resource "sitehost_dns_record" "spf" {
  domain = sitehost_dns_zone.zone.name
  name   = sitehost_dns_zone.zone.name
  type   = "TXT"
  record = "v=spf1 include:_spf.google.com ~all"
}
