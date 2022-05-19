resource "sitehost_server" "web" {
	label = "pw-trtest"
	location = "SHQLIN"
	product_code = "XENLIT"
	image = "ubuntu-focal.amd64"

 connection {
   host = sitehost_server.web.ips[0]
   user = "root"
   type = "ssh"
   #private_key = file(var.pvt_key)
   password = sitehost_server.web.password
   timeout = "2m"
 }

 provisioner "remote-exec" {
   inline = [
     "export PATH=$PATH:/usr/bin",
     # install nginx
     "sudo apt update",
     "sudo apt install -y nginx"
   ]
 }
}

