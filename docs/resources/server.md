---
page_title: "sitehost_server Resource - terraform-provider-sitehost"
subcategory: ""
description: Provides a SiteHost Server resource. This can be used to create, modify, and delete Servers.

---

# sitehost_server (Resource)

Provides a SiteHost Server resource. This can be used to create, modify, and delete Servers.

## Example Usage
```hcl
# Create a web server
resource "sitehost_server" "web" {
	label = "webserver"
	location = "AKLCITY"
	product_code = "XENLIT"
	image = "ubuntu-xenial.amd64"
	ssh_keys = []
}
```

## Schema

### Required

- `image` (String) The Server image slug.
- `label` (String) The Server label.
- `location` (String) The region to start in.
- `product_code` (String) The unique slug that identifies the type of Server.

### Optional

- `ips` (List of String) Assign specific IPs for this Server.
- `name` (String) The name (ID) of the Server.
- `ssh_keys` (List of String) The SSH public keys will added to Server.

### Read-Only

- `id` (String) The ID of this resource.
- `password` (String, Sensitive) The root/administrator password.


