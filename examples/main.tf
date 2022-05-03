terraform {
	required_providers {
		sitehost = {
			version = "0.1.0"
			source = "hashicorp.com/sh/sitehost"
		}
	}
}

provider "sitehost" {
	client_id = "XXX"
	apikey = "XXX"
}

output "terra0" {
	value = sitehost_server.terra0
	sensitive = true
}

resource "sitehost_server" "terra0" {
	label = "TerraForm"
	location = "SHQLIN"
	product_code = "XENPRO"
	image = "ubuntu-xenial.amd64"
}
