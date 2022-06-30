terraform {
    required_providers {
        sitehost = {
            source = "sitehostnz/sitehost"
            version = "~> 1.0"
        }
    }
}

provider "sitehost" {
    client_id = "****"
    api_key = "***"
    api_endpoint = "https://api.sitehost.nz/1.1/"
}

output "server_ips" {
	value = sitehost_server.web.ips
}

output "password" {
	value = nonsensitive(sitehost_server.web.password)
}

resource "sitehost_server" "web" {
    label = "webserver"
    location = "AKLCITY"
    product_code = "XENLIT"
    image = "ubuntu-jammy-pvh.amd64"
}
