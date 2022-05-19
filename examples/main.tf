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
	product_code = "XENLIT"
	image = "ubuntu-xenial.amd64"
	ssh_keys = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCxZhDviV2VIBYLWAdrkSMoTS4RqNZya1OJ6izeDrYmpbTZGqr2CkUf6y+TdqZkngzOdmFO03/DFn3k4O8+XPOv9NoVXZoDWPcxRl4jBEIOR+bB+ZDSWLGCEZGluPCOS1w6TuI4SYkFFLDnBI/w9y/hRogNTZgSaXJolNjkxTO4fQyzUbNa7cShpFQ0sWVe/uKyJXDrPmY3DqBX2G21k0y+/zwOG/n6a6SQapIEY4dMYQSt0KS7tRvkBERwD/Oa/HCYG+EmU4twRGOKHOh/EmGqGyKxY54vY+7wAXL9jChq5fikbMMeGBBUvGySNIT69uiVzc5WlMJCtxk3e0bfSlsTWLTsgMkMorlkFM0EKO+C0Eq+wjcMYZjkyFIx4r7Qz0xzg12OR1UA2Q5ODcWOfNtZ0KGD+wKc3s6flp6HImUFptI9yC9yZY7rLbowAz3ZZapxYl7JTd1uoZNe747uP34/tQgk1Tk4yfIH/wXLxh6vTn4fGpFeN24mfASt7w7qwj1BVTX0UGxSXUfBHMokr/7fCy1PnPmQBVLE4jEWYW2o+R+yGAoizLtWtwmsAKjGzRbCSkH5Fh9ZKA5yv9WsiA8T5sYDLp1a4siF9L5HHEtPveZA8dyoofzm3g/kJw+VDytE4Y8k4N7RaZh/BDRzSIlFFIZ6b4SKis2U5hH5kxEITw== gonzalo@sitehost.co.nz"]
}
