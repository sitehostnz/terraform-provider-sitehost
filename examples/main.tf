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

output "terra0" {
	value = sitehost_server.terra0
	sensitive = true
}

resource "sitehost_server" "terra0" {
	label = "trtest"
	location = "SHQLIN"
	product_code = "XENPRO"
	image = "ubuntu-xenial.amd64"
}
