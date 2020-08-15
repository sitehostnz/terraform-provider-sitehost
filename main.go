package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: func() *schema.Provider { return &sitehost.Provider }})
}
