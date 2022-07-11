
output "server_ips" {
  value = sitehost_server.web.ips
}

output "password" {
  value = nonsensitive(sitehost_server.web.password)
}
