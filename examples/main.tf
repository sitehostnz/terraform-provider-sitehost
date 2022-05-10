terraform {
	required_providers {
		sitehost = {
			version = "0.1.1"
			source = "hashicorp.com/sh/sitehost"
		}
	}
}

provider "sitehost" {
	client_id = ""
	api_key = ""
	api_endpoint = "https://api.staging.sitehost.nz/1.1/"
}

output "server_ips" {
	value = sitehost_server.web.ips
}

output "password" {
	value = nonsensitive(sitehost_server.web.password)
}

resource "sitehost_server" "web" {
	label = "trtest"
	location = "SHQLIN"
	product_code = "XENPRO"
	image = "ubuntu-xenial.amd64"
}
