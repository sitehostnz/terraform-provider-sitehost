<a href="https://terraform.io">
    <img src=".github/terraform_logo.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# Terraform Provider for SiteHost

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.18 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/sitehostnz/terraform-provider-sitehost`

```sh
$ git clone git@github.com:sitehostnz/terraform-provider-sitehost
```

Enter the provider directory and build the provider

```sh
$ make install
```

Example Usage
----------------------

```
# Create a new Server in the AKLCITY region
resource "sitehost_server" "web" {
    image  = "ubuntu-focal.amd64"
    label   = "web-1"
    location = "AKLCITY"
    product_code   = "XENLIT"
}
```