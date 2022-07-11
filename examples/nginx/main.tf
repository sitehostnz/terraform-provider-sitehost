terraform {
    required_providers {
        sitehost = {
            source = "sitehostnz/sitehost"
            version = "~> 1.0"
        }
    }
}

provider "sitehost" {
    client_id = var.sitehost_client_id
    api_key = var.sitehost_api_key
    api_endpoint = "https://api.sitehost.nz/1.1/"
}

resource "sitehost_server" "web" {
    label = "webserver"
    location = var.location
    product_code = var.product_code
    image = var.image
    ssh_keys = [chomp(file(var.public_ssh_key))]

    connection {
        host = self.ips[0]
        user = "root"
        type = "ssh"
        private_key = chomp(file(var.private_ssh_key))
        timeout = "2m"
    }

    provisioner "remote-exec" {
        inline = [
            "export PATH=$PATH:/usr/bin",
            # install nginx
            "apt-get -q update",
            "apt-get -q -y install nginx",
            "echo it works ${self.label} > /var/www/html/index.html",
        ]
    }
}

