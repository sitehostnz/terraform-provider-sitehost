Terraform Provider
==================

- Website: https://docs.sitehost.nz

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by [SiteHost](https://sitehost.nz).

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.13.x
-	[Go](https://golang.org/doc/install) 1.15 (to build the provider plugin)

Usage
---------------------

```tf
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
  // token = "read from SITEHOST_TOKEN environment variable"

  // optional
  url = "https://mysth.safeserver.net.nz/1.0"
}

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
}
```

Building The Provider
---------------------

Clone repository:

```sh
$ git clone git@github.com:sitehostnz/terraform-provider-sitehost
```

Enter the provider directory and build the provider

```sh
$ cd terraform-provider-sitehost
$ go build
```

Using the provider
----------------------
## Fill in for each provider

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.15+ is *required*).

To compile the provider, run `go build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ go get github.com/sitehostnz/terraform-provider-sitehost
...
$ $GOPATH/bin/terraform-provider-sitehost
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ go test -short ./...
```

In order to run the full suite of tests.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ go test ./...
```
