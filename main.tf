terraform {
	required_providers {
		sitehost = {
			version = "0.0.0"
			source = "sitehost.nz/terraform/sitehost"
		}
	}
}

provider "sitehost" {
	// required
	id = 970180
	// token = "xxx"

	// optional
	url = "https://mysth.safeserver.net.nz/1.0"
}

/*
data "sitehost_servers" "all" {}

output "all_servers" {
	value = data.sitehost_servers.all.servers
}
*/

output "terra0" {
	value = sitehost_server.terra0
}

resource "sitehost_server" "terra0" {
	// required
	label = "test"
	location = "AKLCITY"
	product_code = "XENLIT"
	image = "ubuntu-xenial.amd64"

	// optional
	name = "myserver"
	/*
	ips = [ "120.138.19.22", "2403:7000:8000:400::b4" ]
	ssh_keys = ["ssh-rsa AAAAB3NzaC1yc2EAAA... user@host"]
	contact_id = 3
	backup = false
	send_email = false
	*/
}

/*
resource "sitehost_snapshot" "foo" {
	// required
	server = var.sitehost_server.servername
	partition = "xvda2"

	// optional
	lifetime = 1
}
*/
