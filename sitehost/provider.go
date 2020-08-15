// Package sitehost exposes a terraform provider.
package sitehost

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/terraform-provider-sitehost/internal/api"
)

// Provider represents a resource provider in Terraform, and properly
// implements all of the ResourceProvider API.
var Provider = schema.Provider{
	// Schema is the configuration schema for the provider
	// itself. You should define any API keys, etc. here. Schemas
	// are covered below.
	Schema: map[string]*schema.Schema{
		"id": {
			Type:         schema.TypeInt,
			Required:     true,
			DefaultFunc:  schema.EnvDefaultFunc("SITEHOST_ID", nil),
			Description:  "client identifier",
			ValidateFunc: unsigned,
		}, "token": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("SITEHOST_TOKEN", nil),
			Description: "api authentication key",
			Sensitive:   true,
		}, "url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "url prefix of the api server",
		},
	},
	// ConfigureContextFunc is the function used to configure a Provider.
	//
	// The interface{} value returned by this function is stored and passed into
	// the subsequent resources as the meta parameter. This return value is usually
	// used to pass along a configured API client, a configuration structure, etc.
	ConfigureContextFunc: providerConfigure,
	// DataSourcesMap is the collection of available data sources that this
	// provider implements, with a Resource instance defining the schema and Read
	// operation of each.
	//
	// Resource instances for data sources must have a Read function and must *not*
	// implement Create, Update or Delete.
	DataSourcesMap: map[string]*schema.Resource{
		"sitehost_servers": &dataSourceServers,
	},
	// ResourcesMap is the map of resources that this provider
	// supports. All keys are resource names and the values are the
	// schema.Resource structures implementing this resource.
	ResourcesMap: map[string]*schema.Resource{
		"sitehost_server": &resourceServer,
	},
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	c := api.Client{
		Token: d.Get("token").(string),
		ID:    uint(d.Get("id").(int)),
	}
	if v, ok := d.GetOk("url"); ok {
		u, err := url.Parse(v.(string))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		c.Base = *u
	}
	return c, nil
}
